// Package repo implement the function that contact the db required by the service package
package repo

import (
	"github.com/RangelReale/osin"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/component/oauth2/service"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/database/dbconn"
	"github.com/go-errors/errors"
	"github.com/jinzhu/copier"
)

// GetClient loads the client by id
func (s *Storage) GetClient(id string) (osin.Client, error) {
	var client service.Client
	if err := dbconn.DB.Where("id = ?", id).Find(&client).Error; err != nil {
		return nil, err
	}

	var c osin.DefaultClient
	copier.Copy(&c, &client)
	c.UserData = client.UserId
	return &c, nil
}

// UpdateClient updates the client (identified by it's id) and replaces the values with the values of client.
func (s *Storage) UpdateClient(c osin.Client) error {
	userId, ok := c.GetUserData().(uint)
	if !ok {
		return errors.New("cannot assert user_id is uint")
	}

	var client service.Client
	client.Id = c.GetId()
	client.Secret = c.GetSecret()
	client.RedirectUri = c.GetRedirectUri()
	client.UserId = userId

	return dbconn.DB.Save(&client).Error
}

// CreateClient stores the client in the database and returns an error, if something went wrong.
func (s *Storage) CreateClient(c osin.Client) error {
	userId, ok := c.GetUserData().(uint)
	if !ok {
		return errors.New("cannot assert user_id is uint")
	}

	var client service.Client
	client.Id = c.GetId()
	client.Secret = c.GetSecret()
	client.RedirectUri = c.GetRedirectUri()
	client.UserId = userId

	if !dbconn.DB.NewRecord(client) {
		return errors.New("client already exist")
	} else if err := dbconn.DB.Create(&client).Error; err != nil {
		return err
	} else if dbconn.DB.NewRecord(client) {
		return errors.New("client was not created")
	}
	return nil
}

// RemoveClient removes a client (identified by id) from the database. Returns an error if something went wrong.
func (s *Storage) RemoveClient(id string) (err error) {
	return dbconn.DB.Where("id = ?", id).Delete(&service.Client{}).Error
}
