package repo

import (
	"github.com/adriendomoison/gobootapi/database/dbconn"
	"github.com/adriendomoison/gobootapi/user/repo/dbmodel"
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

// Create create a user in Database
func (repo *repo) Create(user dbmodel.Entity) bool {
	if dbconn.DB.NewRecord(user) {
		dbconn.DB.Create(&user)
	}
	return !dbconn.DB.NewRecord(user)
}

// FindByID Find user in Database by ID
func (repo *repo) FindByID(id uint) (user dbmodel.Entity, err error) {
	if err = dbconn.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return dbmodel.Entity{}, err
	}
	return user, nil
}

// FindByEmail Find user in Database by email
func (repo *repo) FindByEmail(email string) (user dbmodel.Entity, err error) {
	if err = dbconn.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return dbmodel.Entity{}, err
	}
	return user, nil
}

// Update edit user in Database
func (repo *repo) Update(user dbmodel.Entity) error {
	return dbconn.DB.Save(&user).Error
}

// Delete remove user from Database
func (repo *repo) Delete(user dbmodel.Entity) error {
	return dbconn.DB.Delete(&user).Error
}
