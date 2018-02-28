package dbmodel

import (
	"github.com/jinzhu/gorm"
	profiledbmodel "github.com/adriendomoison/gobootapi/profile/repo/dbmodel"
)

type Interface interface {
	Create(user Entity) bool
	CreateWithProfile(user Entity, profile profiledbmodel.Entity) error
	FindByID(id uint) (user Entity, err error)
	FindByEmail(email string) (user Entity, err error)
	Update(user Entity) error
	Delete(user Entity) error
}

type Entity struct {
	gorm.Model
	Email    string `gorm:"NOT NULL;UNIQUE"`
	Username string
	Password string `gorm:"NOT NULL"`
}

func (Entity) TableName() string {
	return "user"
}
