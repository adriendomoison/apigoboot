// Package repo implement the function that contact the db required by the service package
package repo

import (
	"github.com/RangelReale/osin"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/component/oauth2/service"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/database/dbconn"
	"github.com/go-errors/errors"
	"github.com/jinzhu/copier"
)

// SaveAccess writes AccessData.
// If RefreshToken is not blank, it must save in a way that can be loaded using LoadRefresh.
func (s *Storage) SaveAccess(data *osin.AccessData) (err error) {
	var access service.Access
	prev := ""
	authorizeData := &osin.AuthorizeData{}

	if data.AccessData != nil {
		prev = data.AccessData.AccessToken
	}

	if data.AuthorizeData != nil {
		authorizeData = data.AuthorizeData
	}

	userId, ok := data.UserData.(uint)
	if !ok {
		return errors.New("cannot assert user_id is uint")
	}

	tx := dbconn.DB.Begin()

	if data.RefreshToken != "" {
		if err := s.saveRefresh(tx, data.RefreshToken, data.AccessToken); err != nil {
			return err
		}
	}

	if data.Client == nil {
		return errors.New("data.Client must not be nil")
	}

	copier.Copy(&access, data)
	access.Client = data.Client.GetId()
	access.UserId = userId
	access.Previous = prev
	access.Authorize = authorizeData.Code

	if err := tx.Create(&access).Error; err != nil {
		if rbe := tx.Rollback(); rbe != nil {
			return errors.New(rbe)
		}
		return errors.New(err)
	}

	if err = tx.Commit().Error; err != nil {
		return errors.New(err)
	}

	return nil
}

// LoadAccess retrieves access data by token. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (s *Storage) LoadAccess(code string) (*osin.AccessData, error) {
	var userId uint
	var result osin.AccessData
	var access service.Access

	if err := dbconn.DB.Where("code = ?", code).Find(&access).Error; err != nil {
		return nil, err
	}

	copier.Copy(&result, &access)
	result.UserData = userId

	client, err := s.GetClient(access.Client)
	if err != nil {
		return nil, err
	}
	result.Client = client
	result.AuthorizeData, _ = s.LoadAuthorize(access.Authorize)
	prevAccess, _ := s.LoadAccess(access.Previous)
	result.AccessData = prevAccess
	return &result, nil
}

// RemoveAccess revokes or deletes an AccessData.
func (s *Storage) RemoveAccess(code string) (err error) {
	return dbconn.DB.Where("access_token = ?", code).Delete(&service.Access{}).Error
}
