package repo

import (
	"time"
	"github.com/jinzhu/copier"
	"github.com/go-errors/errors"
	"github.com/RangelReale/osin"
	"github.com/adriendomoison/gobootapi/database/dbconn"
	"github.com/adriendomoison/gobootapi/oauth/repo/dbmodel"
)

// SaveAuthorize saves authorize data.
func (s *Storage) SaveAuthorize(data *osin.AuthorizeData) (err error) {
	userId, ok := data.UserData.(uint)
	if !ok {
		return errors.New("cannot assert user_id is uint")
	}

	var authorize dbmodel.Authorize
	copier.Copy(&authorize, &data)
	authorize.UserId = userId
	authorize.Client = data.Client.GetId()

	if !dbconn.DB.NewRecord(authorize) {
		return errors.New("authorize already exist")
	} else if err := dbconn.DB.Create(&authorize).Error; err != nil {
		return err
	} else if dbconn.DB.NewRecord(authorize) {
		return errors.New("authorize was not created")
	}
	return nil
}

// LoadAuthorize looks up AuthorizeData by a code.
// Client information MUST be loaded together.
// Optionally can return error if expired.
func (s *Storage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {

	// Load Authorize from DB
	var authorize dbmodel.Authorize
	if err := dbconn.DB.Where("code = ?", code).Find(&authorize).Error; err != nil {
		return nil, err
	}

	// Copy Authorize in OSIN AuthorizeData
	var data osin.AuthorizeData
	copier.Copy(&data, &authorize)

	c, err := s.GetClient(authorize.Client)
	if err != nil {
		return nil, err
	}

	if data.ExpireAt().Before(time.Now()) {
		return nil, errors.Errorf("Token expired at %s.", data.ExpireAt().String())
	}

	data.Client = c
	return &data, nil
}

// RemoveAuthorize revokes or deletes the authorization code.
func (s *Storage) RemoveAuthorize(code string) (err error) {
	return dbconn.DB.Where("code = ?", code).Delete(&dbmodel.Authorize{}).Error
}
