/*
	Package repo implement the function that contact the db required by the service package
*/
package repo

import (
	"github.com/adriendomoison/apigoboot/profile-micro-service/component/profile/service"
	"github.com/adriendomoison/apigoboot/profile-micro-service/database/dbconn"
)

// Make sure the interface is implemented correctly
var _ service.RepoInterface = (*repo)(nil)

// Implement interface
type repo struct {
	repo service.RepoInterface
}

// New return a new repo instance
func New() *repo {
	dbconn.DB.AutoMigrate(&service.Entity{})
	return &repo{}
}

// Create create profile in Database
func (crud *repo) Create(profile service.Entity) bool {
	if dbconn.DB.NewRecord(profile) {
		dbconn.DB.Create(&profile)
	}
	return !dbconn.DB.NewRecord(profile)
}

// FindByID find profile in Database by ID
func (crud *repo) FindByID(id uint) (profile service.Entity, err error) {
	if err = dbconn.DB.Where("id = ?", id).First(&profile).Error; err != nil {
		return service.Entity{}, err
	}
	return profile, nil
}

// FindByPublicId find profile in Database by public_id
func (crud *repo) FindByPublicId(publicId string) (profile service.Entity, err error) {
	if err = dbconn.DB.Where("public_id = ?", publicId).First(&profile).Error; err != nil {
		return service.Entity{}, err
	}
	return profile, nil
}

// Update edit profile in Database
func (crud *repo) Update(profile service.Entity) error {
	return dbconn.DB.Save(&profile).Error
}

// Delete remove profile from Database
func (crud *repo) Delete(profile service.Entity) error {
	return dbconn.DB.Delete(&profile).Error
}
