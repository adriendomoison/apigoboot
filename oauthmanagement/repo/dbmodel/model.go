package dbmodel

import oauthdbmodel "github.com/adriendomoison/gobootapi/oauth/repo/dbmodel"

type Interface interface {
	FindByAccessToken(token string) (user oauthdbmodel.Access, err error)
}

