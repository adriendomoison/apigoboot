package service

import (
	"github.com/adriendomoison/gobootapi/oauthmanagement/repo/dbmodel"
	"github.com/adriendomoison/gobootapi/oauthmanagement/service/model"
)

// Make sure the interface is implemented correctly
var _ model.Interface = (*service)(nil)

// service implement interface
type service struct {
	repo dbmodel.Interface
}

// New returns a new service instance.
func New(repo dbmodel.Interface) *service {
	return &service{repo}
}

// GetResourceOwnerId ask database to retrieve the id of the user owning the access entity with this access_token
func (s *service) GetResourceOwnerId(token string) uint {
	accessToken, err := s.repo.FindByAccessToken(token)
	if err != nil {
		return 0
	}
	return accessToken.UserId
}