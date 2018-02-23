/*
	service do the link in between the web and the database for the user profile.
	It is responsible to handle data and transform them to be readable either from the database to the web or from the web to the database
*/
package service

import (
	"errors"
	"github.com/jinzhu/copier"
	"github.com/elithrar/simple-scrypt"
	"github.com/adriendomoison/go-boot-api/apicore/helpers/servicehelper"
	"github.com/adriendomoison/go-boot-api/user/repo/dbmodel"
	"github.com/adriendomoison/go-boot-api/user/service/model"
	"github.com/adriendomoison/go-boot-api/user/rest/jsonmodel"
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
	return resDTO
}

// createEntityFromDTO copy all data from a Request DTO to an entity and initialize entity with some initialization values
func createEntityFromDTO(reqDTO jsonmodel.RequestDTO, init bool) (entity dbmodel.Entity, error *servicehelper.Error) {
	copier.Copy(&entity, &reqDTO)
	if init {
		hashedPassword, err := scrypt.GenerateFromPassword([]byte(reqDTO.Password), scrypt.DefaultParams)
		if err != nil {
			return dbmodel.Entity{}, &servicehelper.Error{
				Detail:  errors.New("could not hash the password"),
				Message: "We could not create your password, please try something different",
				Code:    servicehelper.UnexpectedError,
			}
		}
		entity.Password = string(hashedPassword[:])
	}
	return
}

// Add set up and create a user
func (s *service) Add(reqDTO jsonmodel.RequestDTO) (resDTO jsonmodel.ResponseDTO, error *servicehelper.Error) {
	entity, Err := createEntityFromDTO(reqDTO, true)
	if Err != nil {
		return jsonmodel.ResponseDTO{}, Err
	} else if s.repo.Create(entity) {
		return createDTOFromEntity(entity), error
	}
	return jsonmodel.ResponseDTO{}, &servicehelper.Error{Detail: errors.New("could not be created"), Code: servicehelper.AlreadyExist}
}

// Retrieve ask database to retrieve a user from its email
func (s *service) Retrieve(email string) (resDTO jsonmodel.ResponseDTO, error *servicehelper.Error) {
	entity, err := s.repo.FindByEmail(email)
	if err != nil {
		return jsonmodel.ResponseDTO{}, &servicehelper.Error{Detail: errors.New("no result found"), Code: servicehelper.NotFound}
	}
	return createDTOFromEntity(entity), error
}

// EditEmail check user's password and change its email
func (s *service) EditEmail(reqDTO jsonmodel.RequestDTOPutEmail) (resDTO jsonmodel.ResponseDTO, error *servicehelper.Error) {
	entity, err := s.repo.FindByEmail(reqDTO.Email)
	if err != nil {
		return jsonmodel.ResponseDTO{}, &servicehelper.Error{Detail: errors.New("no result found"), Code: servicehelper.NotFound}
	}

	if _, err := s.repo.FindByEmail(reqDTO.NewEmail); err == nil {
		return jsonmodel.ResponseDTO{}, &servicehelper.Error{
			Detail:  errors.New("this email is already used for another account"),
			Message: "Email already used, please chose another email",
			Param:   "new_email",
			Code:    servicehelper.AlreadyExist}
	}

	if err := scrypt.CompareHashAndPassword([]byte(entity.Password), []byte(reqDTO.Password)); err != nil {
		return jsonmodel.ResponseDTO{}, &servicehelper.Error{
			Detail:  errors.New("invalid password"),
			Message: "We could not authenticate you, please type you password again",
			Param:   "password",
			Code:    servicehelper.BadRequest}
	}
	entity.Email = reqDTO.NewEmail

	if err := s.repo.Update(entity); err != nil {
		return jsonmodel.ResponseDTO{}, &servicehelper.Error{
			Detail:  errors.New("could not update user"),
			Message: "We could not update the email of this account, please contact us or try again later",
			Code:    servicehelper.UnexpectedError}
	}

	return createDTOFromEntity(entity), error
}

// EditPassword check user password and overwrite it with a new one
func (s *service) EditPassword(reqDTO jsonmodel.RequestDTOPutPassword) (resDTO jsonmodel.ResponseDTO, error *servicehelper.Error) {
	entity, err := s.repo.FindByEmail(reqDTO.Email)
	if err != nil {
		return jsonmodel.ResponseDTO{}, &servicehelper.Error{Detail: errors.New("no result found"), Code: servicehelper.NotFound}
	}

	if scrypt.CompareHashAndPassword([]byte(entity.Password), []byte(reqDTO.Password)) != nil {
		return jsonmodel.ResponseDTO{}, &servicehelper.Error{
			Detail:  errors.New("invalid password"),
			Message: "We could not authenticate you, please type you password again",
			Param:   "password",
			Code:    servicehelper.BadRequest}
	}

	hashedPassword, err := scrypt.GenerateFromPassword([]byte(reqDTO.NewPassword), scrypt.DefaultParams)
	if err != nil {
		return jsonmodel.ResponseDTO{}, &servicehelper.Error{
			Detail:  errors.New("could not hash the password"),
			Message: "We could not create your password, please try something different",
			Param:   "password",
			Code:    servicehelper.UnexpectedError,
		}
	}
	entity.Password = string(hashedPassword[:])

	if err := s.repo.Update(entity); err != nil {
		return jsonmodel.ResponseDTO{}, &servicehelper.Error{
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
