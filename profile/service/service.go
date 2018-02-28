/*
	service do the link in between the web and the database for the user.
	It is responsible to handle data and transform them to be readable either from the database to the web or from the web to the database
*/
package service

import (
	"time"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/adriendomoison/gobootapi/tool"
	"github.com/adriendomoison/gobootapi/profile/repo/dbmodel"
	"github.com/adriendomoison/gobootapi/profile/service/model"
	"github.com/adriendomoison/gobootapi/profile/rest/jsonmodel"
	"github.com/adriendomoison/gobootapi/apicore/helpers/servicehelper"
	userrepo "github.com/adriendomoison/gobootapi/user/repo"
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

// GetResourceOwnerId ask database to retrieve a user ID from its profile public id
func (s *service) GetResourceOwnerId(publicId string) (userId uint) {
	entity, err := s.repo.FindByPublicId(publicId)
	if err != nil {
		return 0
	}
	return entity.UserID
}

// createDTOFromEntity copy all data from an entity to a Response DTO
func createDTOFromEntity(entity dbmodel.Entity) (resDTO jsonmodel.ResponseDTO, error *servicehelper.Error) {
	copier.Copy(&resDTO, &entity)
	location, _ := time.LoadLocation("UTC")
	resDTO.Birthday = entity.Birthday.In(location).Format("2006-01-02")
	if userEntity, err := userrepo.New().FindByID(entity.UserID); err != nil {
		return jsonmodel.ResponseDTO{}, &servicehelper.Error{
			Detail: errors.New("could not find user associated to this profile"),
			Code:   servicehelper.UnexpectedError,
		}
	} else {
		resDTO.Email = userEntity.Email
	}
	return resDTO, error
}

// createEntityFromDTO copy all data from a Request DTO to an entity and initialize entity with some initialization values
func createEntityFromDTO(reqDTO jsonmodel.RequestDTO, init bool) (entity dbmodel.Entity, error *servicehelper.Error) {
	copier.Copy(&entity, &reqDTO)
	birthday, err := time.Parse("2006-01-02", reqDTO.Birthday)
	if err != nil {
		return dbmodel.Entity{}, &servicehelper.Error{
			Detail: errors.New("birthday is malformed"),
			Code:   servicehelper.UnexpectedError,
		}
	}
	entity.Birthday = &birthday
	if init {
		if userEntity, err := userrepo.New().FindByEmail(reqDTO.Email); err != nil {
			return dbmodel.Entity{}, &servicehelper.Error{
				Detail: errors.New("could not retrieve user to associate with new profile"),
				Code:   servicehelper.UnexpectedError,
			}
		} else {
			entity.UserID = userEntity.ID
		}
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
		return jsonmodel.ResponseDTO{}, err
	} else if s.repo.Create(entity) {
		return createDTOFromEntity(entity)
	}
	return jsonmodel.ResponseDTO{}, &servicehelper.Error{Detail: errors.New("could not be created"), Code: servicehelper.AlreadyExist}
}

// Retrieve ask database to retrieve a profile from its public_id
func (s *service) Retrieve(publicId string) (resDTO jsonmodel.ResponseDTO, error *servicehelper.Error) {
	entity, err := s.repo.FindByPublicId(publicId)
	if err != nil {
		return jsonmodel.ResponseDTO{}, &servicehelper.Error{Detail: errors.New("no result found"), Code: servicehelper.NotFound}
	}
	return createDTOFromEntity(entity)
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

	if err := s.repo.Update(entity); err != nil {
		return jsonmodel.ResponseDTO{}, &servicehelper.Error{
			Detail:  errors.New("could not update profile"),
			Message: "We could not update the profile, please contact us or try again later",
			Code:    servicehelper.UnexpectedError}
	}
	return createDTOFromEntity(entity)
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
