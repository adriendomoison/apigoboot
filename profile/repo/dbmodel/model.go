package dbmodel

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Interface interface {
	Create(profile Entity) bool
	FindByID(id uint) (profile Entity, err error)
	FindByPublicId(publicId string) (profile Entity, err error)
	Update(profile Entity) error
	Delete(profile Entity) error
}

// Entity is the model of profile for the database
type Entity struct {
	gorm.Model
	PublicId          string `gorm:"UNIQUE;NOT NULL"`
	FirstName         string
	LastName          string
	ProfilePictureUrl string
	Birthday          *time.Time
	OrderAmount       uint
	UserID            uint
}

func (Entity) TableName() string {
	return "profile"
}
