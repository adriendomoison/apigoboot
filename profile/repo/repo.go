package repo

import (
	"github.com/adriendomoison/go-boot-api/database/dbconn"
	"github.com/adriendomoison/go-boot-api/profile/repo/dbmodel"
)

// Make sure the interface is implemented correctly
var _ dbmodel.Interface = (*repo)(nil)

// Implement interface
type repo struct {
	repo dbmodel.Interface
}

// New return a new repo instance
func New() *repo {
	dbconn.DB.AutoMigrate(&dbmodel.Entity{})
	return &repo{}
}

// Create create profile in Database
func (crud *repo) Create(profile dbmodel.Entity) bool {
	if dbconn.DB.NewRecord(profile) {
		dbconn.DB.Create(&profile)
	}
	return !dbconn.DB.NewRecord(profile)
}

// FindByID Find profile in Database by ID
func (crud *repo) FindByID(id uint) (profile dbmodel.Entity, err error) {
	if err = dbconn.DB.Where("id = ?", id).First(&profile).Error; err != nil {
		return dbmodel.Entity{}, err
	}
	return profile, nil
}

// FindByPublicId Find profile in Database by public_id
func (crud *repo) FindByPublicId(publicId string) (profile dbmodel.Entity, err error) {
	if err = dbconn.DB.Where("public_id = ?", publicId).First(&profile).Error; err != nil {
		return dbmodel.Entity{}, err
	}
	return profile, nil
}

// Update edit profile in Database
func (crud *repo) Update(profile dbmodel.Entity) error {
	return dbconn.DB.Save(&profile).Error
}

// Delete remove profile from Database
func (crud *repo) Delete(profile dbmodel.Entity) error {
	return dbconn.DB.Delete(&profile).Error
}
