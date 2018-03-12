package service

import (
	"github.com/adriendomoison/apigoboot/errorhandling/servicehelper"
	"github.com/adriendomoison/apigoboot/user-micro-service/component/user/rest"
	"github.com/adriendomoison/apigoboot/user-micro-service/config"
	"github.com/adriendomoison/apigoboot/user-micro-service/database/dbconn"
	"os"
	"testing"
)

var email = "john.doe@example.dev"
var password = "mySecretPassword"
var s *service

// Make sure the interface is implemented correctly
var _ RepoInterface = (*repo)(nil)

// Implement interface
type repo struct {
	repo RepoInterface
}

// New return a new repo instance
func NewRepoMock() *repo {
	dbconn.DB.AutoMigrate(&Entity{})
	return &repo{}
}

// Create create a user in Database
func (repo *repo) Create(user Entity) bool {
	if dbconn.DB.NewRecord(user) {
		dbconn.DB.Create(&user)
	}
	return !dbconn.DB.NewRecord(user)
}

// FindByID find user in Database by ID
func (repo *repo) FindByID(id uint) (user Entity, err error) {
	if err = dbconn.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return Entity{}, err
	}
	return user, nil
}

// FindByEmail find user in Database by email
func (repo *repo) FindByEmail(email string) (user Entity, err error) {
	if err = dbconn.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return Entity{}, err
	}
	return user, nil
}

// Update edit user in Database
func (repo *repo) Update(user Entity) error {
	return dbconn.DB.Save(&user).Error
}

// Delete remove user from Database
func (repo *repo) Delete(user Entity) error {
	return dbconn.DB.Delete(&user).Error
}

func TestMain(m *testing.M) {
	config.SetToTestingEnv()
	dbconn.Connect()
	defer dbconn.DB.Close()
	s = New(NewRepoMock())

	code := m.Run()

	dbconn.DB.DropTable(&Entity{})

	os.Exit(code)
}

func TestAddUser(t *testing.T) {
	reqDTO := rest.RequestDTO{
		Email:    email,
		Password: password,
	}

	_, Err := s.Add(reqDTO)

	if Err != nil {
		t.Error(Err.Detail)
	}
}

func TestAddUserWithSameEmail(t *testing.T) {
	reqDTO := rest.RequestDTO{
		Email:    email,
		Password: password,
	}

	_, Err := s.Add(reqDTO)

	if Err == nil {
		t.Error("User already exist and no error was raised")
	}

	if Err != nil && Err.Code != servicehelper.AlreadyExist {
		t.Errorf("Code was supposde to be Already Exist but got %v", Err.Code)
	}
}

func TestRetrieveUser(t *testing.T) {
	resDTO, Err := s.Retrieve(email)

	if Err != nil {
		t.Error(Err.Detail)
	}

	if resDTO.Email != email {
		t.Error("Wrong user was retrieved")
	}
}

func TestErrorMessageOnRetrieveUserThatDoesNotExist(t *testing.T) {
	_, Err := s.Retrieve("toto@example.dev")

	if Err == nil {
		t.Error("No error raised while used doesn't exist")
	}

	if Err != nil && Err.Code != servicehelper.NotFound {
		t.Errorf("Code should be: Not Found, got %v", Err.Code)
	}
}

func TestErrorMessageOnEditUserThatDoesNotExist(t *testing.T) {

	newEmail := "john.john@example.dev"

	reqDTOEmail := rest.RequestDTOPutEmail{
		Email:    "toto@example.dev",
		Password: password,
		NewEmail: newEmail,
	}
	_, Err := s.EditEmail(reqDTOEmail)

	if Err == nil {
		t.Error("No error raised while used doesn't exist")
	}

	if Err != nil && Err.Code != servicehelper.NotFound {
		t.Errorf("Code should be: Not Found, got %v", Err.Code)
	}

	reqDTOPassword := rest.RequestDTOPutPassword{
		Email:       "toto@example.dev",
		Password:    password,
		NewPassword: password + "New",
	}

	_, Err = s.EditPassword(reqDTOPassword)

	if Err == nil {
		t.Error("No error raised while used doesn't exist")
	}

	if Err != nil && Err.Code != servicehelper.NotFound {
		t.Errorf("Code should be: Not Found, got %v", Err.Code)
	}

}

func TestEditUserEmail(t *testing.T) {

	newEmail := "john.john@example.dev"

	reqDTO := rest.RequestDTOPutEmail{
		Email:    email,
		Password: password,
		NewEmail: newEmail,
	}

	email = newEmail

	resDTO, Err := s.EditEmail(reqDTO)

	if Err != nil {
		t.Error(Err.Detail)
	}

	if resDTO.Email != email {
		t.Error("Email was not updated")
	}
}

func TestEditUserEmailWithoutProvidingPassword(t *testing.T) {

	newEmail := "john.fake@example.dev"

	reqDTO := rest.RequestDTOPutEmail{
		Email:    email,
		NewEmail: newEmail,
	}

	_, Err := s.EditEmail(reqDTO)

	if Err == nil {
		t.Error("There should be an error")
	}
}

func TestEditUserPassword(t *testing.T) {
	reqDTO := rest.RequestDTOPutPassword{
		Email:       email,
		Password:    password,
		NewPassword: "myNewSecretPassword",
	}

	_, Err := s.EditPassword(reqDTO)

	reqDTO = rest.RequestDTOPutPassword{
		Email:       email,
		Password:    "myNewSecretPassword",
		NewPassword: "mySecretPassword",
	}

	_, Err = s.EditPassword(reqDTO)

	if Err != nil {
		t.Error("User password was not edited")
	}
}

func TestEditUserPasswordWithoutProvidingOldPassword(t *testing.T) {
	reqDTO := rest.RequestDTOPutPassword{
		Email:       email,
		NewPassword: "myNewSecretPassword",
	}

	_, Err := s.EditPassword(reqDTO)

	if Err == nil {
		t.Error("User password should not have been edited")
	}
}

func TestRemoveUser(t *testing.T) {
	if s.Remove(email) != nil {
		t.Error("Could not deleted user")
	}
	_, Err := s.Retrieve(email)
	if Err == nil {
		t.Error("Could find user after is was supposed to be deleted")
	}
}

func TestRemoveUserThatDoesNotExist(t *testing.T) {
	if s.Remove("toto@example.dev") == nil {
		t.Error("Delete success on user that did not exist")
	}
}
