// Package repo implement the function that contact the db required by the service package
package repo

import (
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/component/oauth2/service"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/database/dbconn"
)

// Make sure the interface is implemented correctly
var _ service.RepoInterface = (*repo)(nil)

// Implement interface
type repo struct {
	repo service.RepoInterface
}

// New returns a new repo instance.
func New() *repo {
	return &repo{}
}

// FindByAccessToken find access in Database by access token
func (r *repo) FindByAccessToken(accessToken string) (at service.Access, err error) {
	if err := dbconn.DB.Where("access_token = ?", accessToken).First(&at).Error; err != nil {
		return service.Access{}, err
	}
	return
}
