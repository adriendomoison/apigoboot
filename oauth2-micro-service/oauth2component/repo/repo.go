package repo

import (
	"github.com/adriendomoison/gobootapi/oauth2-micro-service/database/dbconn"
	"github.com/adriendomoison/gobootapi/oauth2-micro-service/oauth2component/service"
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
