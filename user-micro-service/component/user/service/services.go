// Package service implement the services required by the rest package
package service

import (
	"errors"
	"github.com/adriendomoison/apigoboot/api-tool/errorhandling/servicehelper"
	"github.com/adriendomoison/apigoboot/user-micro-service/component/user/rest"
	"github.com/elithrar/simple-scrypt"
	"github.com/jinzhu/copier"
	"time"
)

// RepoInterface is the model for the repo package of user
type RepoInterface interface {
	Create(user Entity) bool
	FindByID(id uint) (user Entity, err error)
	FindByEmail(email string) (user Entity, err error)
	Update(user Entity) error
	Delete(user Entity) error
}

// Entity is the model of a user in the database
type Entity struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string `gorm:"NOT NULL;UNIQUE"`
	Username  string
	Password  string `gorm:"NOT NULL"`
}

// TableName allow to gives a specific name to the user table
func (Entity) TableName() string {
	return "user"
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

// createDTOFromEntity copy all data from an entity to a Response DTO
func createDTOFromEntity(entity Entity) (resDTO rest.ResponseDTO) {
	copier.Copy(&resDTO, &entity)
	return resDTO
}

// createEntityFromDTO copy all data from a Request DTO to an entity and initialize entity with some initialization values
func createEntityFromDTO(reqDTO rest.RequestDTO, init bool) (entity Entity, error *servicehelper.Error) {
	copier.Copy(&entity, &reqDTO)
	if init {
		hashedPassword, err := scrypt.GenerateFromPassword([]byte(reqDTO.Password), scrypt.DefaultParams)
		if err != nil {
			return Entity{}, &servicehelper.Error{
				Detail:  errors.New("could not hash the password"),
				Message: "We could not create your password, please try something different",
				Code:    servicehelper.UnexpectedError,
			}
		}
		entity.Password = string(hashedPassword[:])
	}
	return
}

// GetResourceOwnerId ask database to retrieve a user ID from its email
func (s *service) GetResourceOwnerId(email string) (userId uint) {
	entity, err := s.repo.FindByEmail(email)
	if err != nil {
		return 0
	}
	return entity.ID
}

// Add set up and create a user
func (s *service) Add(reqDTO rest.RequestDTO) (resDTO rest.ResponseDTO, error *servicehelper.Error) {
	entity, Err := createEntityFromDTO(reqDTO, true)
	if Err != nil {
		return rest.ResponseDTO{}, Err
	} else if s.repo.Create(entity) {
		return createDTOFromEntity(entity), error
	}
	return rest.ResponseDTO{}, &servicehelper.Error{
		Detail:  errors.New("user could not be created"),
		Message: "We could not created this account because an account already use this email address",
		Param:   "email",
		Code:    servicehelper.AlreadyExist,
	}
}

// Retrieve ask database to retrieve a user from its email
func (s *service) Retrieve(email string) (resDTO rest.ResponseDTO, error *servicehelper.Error) {
	entity, err := s.repo.FindByEmail(email)
	if err != nil {
		return rest.ResponseDTO{}, &servicehelper.Error{
			Detail:  errors.New("user could not be found"),
			Message: "We could not find any user with the provided email address",
			Param:   "email",
			Code:    servicehelper.NotFound,
		}
	}
	return createDTOFromEntity(entity), error
}

// EditEmail check user's password and change its email
func (s *service) EditEmail(reqDTO rest.RequestDTOPutEmail) (resDTO rest.ResponseDTO, error *servicehelper.Error) {
	entity, err := s.repo.FindByEmail(reqDTO.Email)
	if err != nil {
		return rest.ResponseDTO{}, &servicehelper.Error{
			Detail:  errors.New("user could not be found"),
			Message: "We could not find any user with the provided email address",
			Param:   "email",
			Code:    servicehelper.NotFound,
		}
	}

	if _, err := s.repo.FindByEmail(reqDTO.NewEmail); err == nil {
		return rest.ResponseDTO{}, &servicehelper.Error{
			Detail:  errors.New("this email is already used for another account"),
			Message: "Email already used, please chose another email",
			Param:   "new_email",
			Code:    servicehelper.AlreadyExist,
		}
	}

	if err := scrypt.CompareHashAndPassword([]byte(entity.Password), []byte(reqDTO.Password)); err != nil {
		return rest.ResponseDTO{}, &servicehelper.Error{
			Detail:  errors.New("invalid password"),
			Message: "We could not authenticate you, please type you password again",
			Param:   "password",
			Code:    servicehelper.BadRequest,
		}
	}
	entity.Email = reqDTO.NewEmail

	if err := s.repo.Update(entity); err != nil {
		return rest.ResponseDTO{}, &servicehelper.Error{
			Detail:  errors.New("could not update user"),
			Message: "We could not update the email of this account, please contact us or try again later",
			Code:    servicehelper.UnexpectedError,
		}
	}

	return createDTOFromEntity(entity), error
}

// EditPassword check user password and overwrite it with a new one
func (s *service) EditPassword(reqDTO rest.RequestDTOPutPassword) (resDTO rest.ResponseDTO, error *servicehelper.Error) {
	entity, err := s.repo.FindByEmail(reqDTO.Email)
	if err != nil {
		return rest.ResponseDTO{}, &servicehelper.Error{
			Detail:  errors.New("user could not be found"),
			Message: "We could not find any user with the provided email address",
			Param:   "email",
			Code:    servicehelper.NotFound,
		}
	}

	if scrypt.CompareHashAndPassword([]byte(entity.Password), []byte(reqDTO.Password)) != nil {
		return rest.ResponseDTO{}, &servicehelper.Error{
			Detail:  errors.New("invalid password"),
			Message: "We could not authenticate you, please type you password again",
			Param:   "password",
			Code:    servicehelper.BadRequest,
		}
	}

	hashedPassword, err := scrypt.GenerateFromPassword([]byte(reqDTO.NewPassword), scrypt.DefaultParams)
	if err != nil {
		return rest.ResponseDTO{}, &servicehelper.Error{
			Detail:  errors.New("could not hash the password"),
			Message: "We could not create your password, please try something different",
			Param:   "password",
			Code:    servicehelper.UnexpectedError,
		}
	}
	entity.Password = string(hashedPassword[:])

	if err := s.repo.Update(entity); err != nil {
		return rest.ResponseDTO{}, &servicehelper.Error{
			Detail:  errors.New("could not update user"),
			Message: "We could not update the password of this account, please contact us or try again later",
			Code:    servicehelper.UnexpectedError}
	}

	return createDTOFromEntity(entity), error
}

// Remove find a user in the database and delete it
func (s *service) Remove(email string) (error *servicehelper.Error) {
	if entity, err := s.repo.FindByEmail(email); err != nil {
		return &servicehelper.Error{
			Detail:  errors.New("could not find user"),
			Message: "We could not find any user, please check the provided email",
			Param:   "email",
			Code:    servicehelper.BadRequest,
		}
	} else if err := s.repo.Delete(entity); err != nil {
		return &servicehelper.Error{
			Detail:  errors.New("failed to delete user"),
			Message: "We could not delete the user, please try again later",
			Code:    servicehelper.UnexpectedError,
		}
	}
	return
}

// IsThatTheUserId check if userIdToCheck is the same than the resource
func (s *service) IsThatTheUserId(email string, userIdToCheck uint) (same bool, error *servicehelper.Error) {
	entity, err := s.repo.FindByEmail(email)
	if err != nil {
		return false, &servicehelper.Error{
			Detail:  errors.New("could not find user"),
			Message: "We could not find any user, please check the provided email",
			Param:   "email",
			Code:    servicehelper.BadRequest,
		}
	}
	return entity.ID == userIdToCheck, nil
}
