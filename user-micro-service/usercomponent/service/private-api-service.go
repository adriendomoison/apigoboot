package service

import (
	"errors"
	"github.com/adriendomoison/gobootapi/errorhandling/servicehelper"
	"github.com/adriendomoison/gobootapi/user-micro-service/usercomponent/rest"
)

// RetrieveUserInfoByEmail ask database to retrieve a user from its email
func (s *service) RetrieveUserInfoByEmail(email string) (resDTO rest.ResponseDTOInternalUserInfo, error *servicehelper.Error) {
	entity, err := s.repo.FindByEmail(email)
	if err != nil {
		return rest.ResponseDTOInternalUserInfo{}, &servicehelper.Error{Detail: errors.New("no result found"), Code: servicehelper.NotFound}
	}
	return rest.ResponseDTOInternalUserInfo{
		UserId: entity.ID,
		Email: entity.Email,
	}, nil
}

// RetrieveUserInfoByEmail ask database to retrieve a user from its user id
func (s *service) RetrieveUserInfoByUserId(userId uint) (resDTO rest.ResponseDTOInternalUserInfo, error *servicehelper.Error) {
	entity, err := s.repo.FindByID(userId)
	if err != nil {
		return rest.ResponseDTOInternalUserInfo{}, &servicehelper.Error{Detail: errors.New("no result found"), Code: servicehelper.NotFound}
	}
	return rest.ResponseDTOInternalUserInfo{
		UserId: entity.ID,
		Email: entity.Email,
	}, nil
}
