// Package service implement the services required by the rest package
package service

import (
	"errors"
	"github.com/adriendomoison/apigoboot/api-tool/errorhandling/servicehelper"
	"github.com/adriendomoison/apigoboot/user-micro-service/component/user/rest"
)

// RetrieveUserInfoByEmail ask database to retrieve a user from its email
func (s *service) RetrieveUserInfoByEmail(email string) (resDTO rest.ResponseDTOUserInfo, error *servicehelper.Error) {
	entity, err := s.repo.FindByEmail(email)
	if err != nil {
		return rest.ResponseDTOUserInfo{}, &servicehelper.Error{Detail: errors.New("no result found"), Code: servicehelper.NotFound}
	}
	return rest.ResponseDTOUserInfo{
		UserId: entity.ID,
		Email:  entity.Email,
	}, nil
}

// RetrieveUserInfoByEmail ask database to retrieve a user from its user id
func (s *service) RetrieveUserInfoByUserId(userId uint) (resDTO rest.ResponseDTOUserInfo, error *servicehelper.Error) {
	entity, err := s.repo.FindByID(userId)
	if err != nil {
		return rest.ResponseDTOUserInfo{}, &servicehelper.Error{Detail: errors.New("no result found"), Code: servicehelper.NotFound}
	}
	return rest.ResponseDTOUserInfo{
		UserId: entity.ID,
		Email:  entity.Email,
	}, nil
}
