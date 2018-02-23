/*
	service do the link in between the web and the database for the user.
	It is responsible to handle data and transform them to be readable either from the database to the web or from the web to the database
*/
package service

import (
	"time"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/adriendomoison/go-boot-api/tool"
	"github.com/adriendomoison/go-boot-api/apicore/helpers/servicehelper"
	"github.com/adriendomoison/go-boot-api/profile/service/model"
	"github.com/adriendomoison/go-boot-api/profile/rest/jsonmodel"
	"github.com/adriendomoison/go-boot-api/profile/repo/dbmodel"
)

// Make sure the interface is implemented correctly
var _ model.Interface = (*service)(nil)

// service implement interface
type service struct {
	repo dbmodel.Interface
}

// New return a new service instance
func New(repo dbmodel.Interface) *service {
	return &service{repo}
}

// createDTOFromEntity copy all data from an entity to a Response DTO
func createDTOFromEntity(entity dbmodel.Entity) (resDTO jsonmodel.ResponseDTO) {
	copier.Copy(&resDTO, &entity)
	location, _ := time.LoadLocation("UTC")
	resDTO.Birthday = entity.Birthday.In(location).Format("2006-01-02")
	return resDTO
}

// createEntityFromDTO copy all data from a Request DTO to an entity and initialize entity with some initialization values
func createEntityFromDTO(reqDTO jsonmodel.RequestDTO, init bool) (entity dbmodel.Entity, err error) {
	copier.Copy(&entity, &reqDTO)
	birthday, err := time.Parse("2006-01-02", reqDTO.Birthday)
	if err != nil {
		return dbmodel.Entity{}, errors.New("birthday is malformed")
	}
	entity.Birthday = &birthday
	if init {
		entity.PublicId = tool.GenerateRandomString(8)
		entity.OrderAmount = 0
		entity.ProfilePictureUrl = "https://x1.xingassets.com/assets/frontend_minified/img/users/nobody_m.original.jpg"
	}
	return
}

// Add set up and create a profile
func (s *service) Add(reqDTO jsonmodel.RequestDTO) (resDTO jsonmodel.ResponseDTO, error *servicehelper.Error) {
	entity, err := createEntityFromDTO(reqDTO, true)
	if err != nil {
		return jsonmodel.ResponseDTO{}, &servicehelper.Error{Detail: err, Code: servicehelper.BadRequest}
	} else if s.repo.Create(entity) {
		return createDTOFromEntity(entity), error
	}
	return jsonmodel.ResponseDTO{}, &servicehelper.Error{Detail: errors.New("could not be created"), Code: servicehelper.AlreadyExist}
}

// Retrieve ask database to retrieve a profile from its public_id
func (s *service) Retrieve(publicId string) (resDTO jsonmodel.ResponseDTO, error *servicehelper.Error) {
	entity, err := s.repo.FindByPublicId(publicId)
	if err != nil {
		return jsonmodel.ResponseDTO{}, &servicehelper.Error{Detail: errors.New("no result found"), Code: servicehelper.NotFound}
	}
	return createDTOFromEntity(entity), error
}

// Edit edit user profile and ask database to save changes
func (s *service) Edit(reqDTO jsonmodel.RequestDTO) (resDTO jsonmodel.ResponseDTO, error *servicehelper.Error) {
	entity, err := s.repo.FindByPublicId(reqDTO.PublicId)
	if err != nil {
		return jsonmodel.ResponseDTO{}, &servicehelper.Error{Detail: errors.New("no result found"), Code: servicehelper.NotFound}
	}

	entity.FirstName = reqDTO.FirstName
	entity.LastName = reqDTO.LastName
	entity.ProfilePictureUrl = reqDTO.ProfilePictureUrl
	birthday, err := time.Parse("2006-01-02", reqDTO.Birthday)
	if err != nil {
		return jsonmodel.ResponseDTO{}, &servicehelper.Error{Detail: errors.New("birthday is malformed"), Code: servicehelper.BadRequest}
	}
	entity.Birthday = &birthday

	s.repo.Update(entity)

	return createDTOFromEntity(entity), error
}

// Remove find a profile in the database and delete it
func (s *service) Remove(publicId string) (error *servicehelper.Error) {
	if entity, err := s.repo.FindByPublicId(publicId); err != nil {
		return &servicehelper.Error{
			Detail:  errors.New("could not find profile"),
			Message: "We could not find any profile, please check the provided public_id",
			Param:   "public_id",
			Code:    servicehelper.BadRequest,
		}
	} else if err := s.repo.Delete(entity); err != nil {
		return &servicehelper.Error{
			Detail:  errors.New("failed to delete profile"),
			Message: "We could not delete the profile, please try again later",
			Code:    servicehelper.UnexpectedError,
		}
	}
	return
}
