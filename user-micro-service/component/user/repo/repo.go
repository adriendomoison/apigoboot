package repo

import (
	"github.com/adriendomoison/apigoboot/user-micro-service/database/dbconn"
	"github.com/adriendomoison/apigoboot/user-micro-service/component/user/service"
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

// Create create a user in Database
func (repo *repo) Create(user service.Entity) bool {
	if dbconn.DB.NewRecord(user) {
		dbconn.DB.Create(&user)
	}
	return !dbconn.DB.NewRecord(user)
}

// FindByID find user in Database by ID
func (repo *repo) FindByID(id uint) (user service.Entity, err error) {
	if err = dbconn.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return service.Entity{}, err
	}
	return user, nil
}

// FindByEmail find user in Database by email
func (repo *repo) FindByEmail(email string) (user service.Entity, err error) {
	if err = dbconn.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return service.Entity{}, err
	}
	return user, nil
}

// Update edit user in Database
func (repo *repo) Update(user service.Entity) error {
	return dbconn.DB.Save(&user).Error
}

// Delete remove user from Database
func (repo *repo) Delete(user service.Entity) error {
	return dbconn.DB.Delete(&user).Error
}
