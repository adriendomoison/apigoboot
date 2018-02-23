package repo

import (
	"github.com/jinzhu/gorm"
	"github.com/RangelReale/osin"
	"github.com/go-errors/errors"
	"github.com/adriendomoison/go-boot-api/database/dbconn"
	"github.com/adriendomoison/go-boot-api/oauth/repo/model"
)

// LoadRefresh retrieves refresh AccessData. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (s *Storage) LoadRefresh(code string) (*osin.AccessData, error) {
	var refresh model.Refresh
	if err := dbconn.DB.Where("token = ?", code).Find(&refresh).Error; err != nil {
		return nil, err
	}
	return s.LoadAccess(refresh.Access)
}

// RemoveRefresh revokes or deletes refresh AccessData.
func (s *Storage) RemoveRefresh(code string) error {
	return dbconn.DB.Where("token = ?", code).Delete(&model.Refresh{}).Error
}

func (s *Storage) saveRefresh(tx *gorm.DB, refresh string, access string) (err error) {
	if err := tx.Create(&model.Refresh{Access: access, Token: refresh}).Error; err != nil {
		if rbe := tx.Rollback(); rbe != nil {
			return errors.New(rbe)
		}
		return errors.New(err)
	}
	return nil
}
