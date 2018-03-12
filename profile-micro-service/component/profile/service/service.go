// Package service implement the services required by the rest package
package service

import (
	"errors"
	"github.com/adriendomoison/apigoboot/api-tool/gentool"
	"github.com/adriendomoison/apigoboot/errorhandling/servicehelper"
	"github.com/adriendomoison/apigoboot/profile-micro-service/component/profile/rest"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"time"
)

// RepoInterface is the model for the repo package of profile
type RepoInterface interface {
	Create(profile Entity) bool
	FindByID(id uint) (profile Entity, err error)
	FindByPublicId(publicId string) (profile Entity, err error)
	Update(profile Entity) error
	Delete(profile Entity) error
}

// Entity is the model of a profile in the database
type Entity struct {
	gorm.Model
	PublicId          string `gorm:"UNIQUE;NOT NULL"`
	FirstName         string
	LastName          string
	ProfilePictureUrl string
	Birthday          *time.Time
	OrderAmount       uint
	UserID            uint
}

// TableName allow to gives a specific name to the profile table
func (Entity) TableName() string {
	return "profile"
}

// Make sure the interface is implemented correctly
var _ rest.ServiceInterface = (*service)(nil)

// service implement interface
type service struct {
	repo RepoInterface
}

// New return a new service instance
func New(repo RepoInterface) *service {
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
func createDTOFromEntity(entity Entity) (resDTO rest.ResponseDTO, error *servicehelper.Error) {
	copier.Copy(&resDTO, &entity)
	location, _ := time.LoadLocation("UTC")
	resDTO.Birthday = entity.Birthday.In(location).Format("2006-01-02")
	userInfo, err := askUserServiceForUserEmail(entity.UserID)
	if err != nil {
		return rest.ResponseDTO{}, &servicehelper.Error{
			Detail: errors.New("can't retrieve user associated to this profile"),
			Code:   servicehelper.UnexpectedError,
		}
	}
	resDTO.Email = userInfo.Email
	return resDTO, error
}

// createEntityFromDTO copy all data from a Request DTO to an entity and initialize entity with some initialization values
func createEntityFromDTO(reqDTO rest.RequestDTOCreation, init bool) (entity Entity, error *servicehelper.Error) {
	copier.Copy(&entity, &reqDTO)
	birthday, err := time.Parse("2006-01-02", reqDTO.Birthday)
	if err != nil {
		return Entity{}, &servicehelper.Error{
			Detail: errors.New("birthday is malformed"),
			Code:   servicehelper.UnexpectedError,
		}
	}
	entity.Birthday = &birthday
	if init {
		userInfo, err := askUserServiceForUserId(reqDTO.Email)
		if err != nil {
			return Entity{}, &servicehelper.Error{
				Detail: errors.New("can't retrieve user associated to this profile"),
				Code:   servicehelper.UnexpectedError,
			}
		}
		entity.UserID = userInfo.UserId
		entity.PublicId = gentool.GenerateRandomString(16)
		entity.OrderAmount = 0
		entity.ProfilePictureUrl = "https://x1.xingassets.com/assets/frontend_minified/img/users/nobody_m.original.jpg"
	}
	return
}

// Add set up and create a profile
func (s *service) Add(reqDTO rest.RequestDTOCreation) (resDTO rest.ResponseDTO, error *servicehelper.Error) {
	entity, err := createEntityFromDTO(reqDTO, true)
	if err != nil {
		return rest.ResponseDTO{}, err
	} else if s.repo.Create(entity) {
		return createDTOFromEntity(entity)
	}
	return rest.ResponseDTO{}, &servicehelper.Error{Detail: errors.New("could not be created"), Code: servicehelper.AlreadyExist}
}

// Retrieve ask database to retrieve a profile from its public_id
func (s *service) Retrieve(publicId string) (resDTO rest.ResponseDTO, error *servicehelper.Error) {
	entity, err := s.repo.FindByPublicId(publicId)
	if err != nil {
		return rest.ResponseDTO{}, &servicehelper.Error{Detail: errors.New("no result found"), Code: servicehelper.NotFound}
	}
	return createDTOFromEntity(entity)
}

// Edit edit user profile and ask database to save changes
func (s *service) Edit(reqDTO rest.RequestDTO) (resDTO rest.ResponseDTO, error *servicehelper.Error) {
	entity, err := s.repo.FindByPublicId(reqDTO.PublicId)
	if err != nil {
		return rest.ResponseDTO{}, &servicehelper.Error{Detail: errors.New("no result found"), Code: servicehelper.NotFound}
	}

	entity.FirstName = reqDTO.FirstName
	entity.LastName = reqDTO.LastName
	entity.ProfilePictureUrl = reqDTO.ProfilePictureUrl
	birthday, err := time.Parse("2006-01-02", reqDTO.Birthday)
	if err != nil {
		return rest.ResponseDTO{}, &servicehelper.Error{Detail: errors.New("birthday is malformed"), Code: servicehelper.BadRequest}
	}
	entity.Birthday = &birthday

	if err := s.repo.Update(entity); err != nil {
		return rest.ResponseDTO{}, &servicehelper.Error{
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

// IsThatTheUserId check if userIdToCheck is the same than the resource
func (s *service) IsThatTheUserId(publicId string, userIdToCheck uint) (same bool, error *servicehelper.Error) {
	entity, err := s.repo.FindByPublicId(publicId)
	if err != nil {
		return false, &servicehelper.Error{
			Detail:  errors.New("could not find profile"),
			Message: "We could not find any profile, please check the provided profile public id",
			Param:   "public_id",
			Code:    servicehelper.BadRequest,
		}
	}
	return entity.ID == userIdToCheck, nil
}
