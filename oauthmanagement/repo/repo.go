package repo

import (
	"github.com/adriendomoison/gobootapi/database/dbconn"
	"github.com/adriendomoison/gobootapi/oauthmanagement/repo/dbmodel"
	oauthdbmodel "github.com/adriendomoison/gobootapi/oauth/repo/dbmodel"
)

// Make sure the interface is implemented correctly
var _ dbmodel.Interface = (*repo)(nil)

// Implement interface
type repo struct {
	repo dbmodel.Interface
}

// New returns a new repo instance.
func New() *repo {
	return &repo{}
}

// FindByAccessToken find access in Database by access token
func (r *repo) FindByAccessToken(accessToken string) (at oauthdbmodel.Access, err error) {
	if err := dbconn.DB.Where("access_token = ?", accessToken).First(&at).Error; err != nil {
		return oauthdbmodel.Access{}, err
	}
	return
}
