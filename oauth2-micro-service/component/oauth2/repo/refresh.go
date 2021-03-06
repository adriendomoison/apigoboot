// Package repo implement the function that contact the db required by the service package
package repo

import (
	"github.com/RangelReale/osin"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/component/oauth2/service"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/database/dbconn"
	"github.com/go-errors/errors"
	"github.com/jinzhu/gorm"
)

// LoadRefresh retrieves refresh AccessData. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (s *Storage) LoadRefresh(code string) (*osin.AccessData, error) {
	var refresh service.Refresh
	if err := dbconn.DB.Where("token = ?", code).Find(&refresh).Error; err != nil {
		return nil, err
	}
	return s.LoadAccess(refresh.Access)
}

// RemoveRefresh revokes or deletes refresh AccessData.
func (s *Storage) RemoveRefresh(code string) error {
	return dbconn.DB.Where("token = ?", code).Delete(&service.Refresh{}).Error
}

func (s *Storage) saveRefresh(tx *gorm.DB, refresh string, access string) (err error) {
	if err := tx.Create(&service.Refresh{Access: access, Token: refresh}).Error; err != nil {
		if rbe := tx.Rollback(); rbe != nil {
			return errors.New(rbe)
		}
		return errors.New(err)
	}
	return nil
}
