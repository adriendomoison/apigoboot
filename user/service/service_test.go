package service

import (
	"os"
	"testing"
	"github.com/adriendomoison/gobootapi/apicore/config"
	"github.com/adriendomoison/gobootapi/database/dbconn"
	"github.com/adriendomoison/gobootapi/user/repo"
	"github.com/adriendomoison/gobootapi/user/repo/dbmodel"
	"github.com/adriendomoison/gobootapi/user/rest/jsonmodel"
	"github.com/adriendomoison/gobootapi/apicore/helpers/servicehelper"
)

var email = "john.doe@example.dev"
var password = "mySecretPassword"
var s *service

func TestMain(m *testing.M) {
	config.SetToTestingEnv()
	dbconn.Connect()
	defer dbconn.DB.Close()
	s = New(repo.New())

	code := m.Run()

	dbconn.DB.DropTable(&dbmodel.Entity{})

	os.Exit(code)
}

func TestAddUser(t *testing.T) {
	reqDTO := jsonmodel.RequestDTO{
		Email:    email,
		Password: password,
	}

	_, Err := s.Add(reqDTO)

	if Err != nil {
		t.Error(Err.Detail)
	}
}

func TestAddUserWithSameEmail(t *testing.T) {
	reqDTO := jsonmodel.RequestDTO{
		Email:    email,
		Password: password,
	}

	_, Err := s.Add(reqDTO)

	if Err == nil {
		t.Error("User already exist and no error was raised")
	}

	if Err != nil && Err.Code != servicehelper.AlreadyExist {
		t.Errorf("Code was supposde to be Already Exist but got %s", Err.Code)
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
		t.Errorf("Code should be: Not Found, got %s", Err.Code)
	}
}

func TestErrorMessageOnEditUserThatDoesNotExist(t *testing.T) {

	newEmail := "john.john@example.dev"

	reqDTOEmail := jsonmodel.RequestDTOPutEmail{
		Email:       "toto@example.dev",
		Password:    password,
		NewEmail:    newEmail,
	}
	_, Err := s.EditEmail(reqDTOEmail)

	if Err == nil {
		t.Error("No error raised while used doesn't exist")
	}

	if Err != nil && Err.Code != servicehelper.NotFound {
		t.Errorf("Code should be: Not Found, got %s", Err.Code)
	}

	reqDTOPassword := jsonmodel.RequestDTOPutPassword{
		Email:       "toto@example.dev",
		Password:    password,
		NewPassword: password + "New",
	}

	_, Err = s.EditPassword(reqDTOPassword)

	if Err == nil {
		t.Error("No error raised while used doesn't exist")
	}

	if Err != nil && Err.Code != servicehelper.NotFound {
		t.Errorf("Code should be: Not Found, got %s", Err.Code)
	}

}

func TestEditUserEmail(t *testing.T) {

	newEmail := "john.john@example.dev"

	reqDTO := jsonmodel.RequestDTOPutEmail{
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

	reqDTO := jsonmodel.RequestDTOPutEmail{
		Email:    email,
		NewEmail: newEmail,
	}

	_, Err := s.EditEmail(reqDTO)

	if Err == nil {
		t.Error("There should be an error")
	}
}

func TestEditUserPassword(t *testing.T) {
	reqDTO := jsonmodel.RequestDTOPutPassword{
		Email:       email,
		Password:    password,
		NewPassword: "myNewSecretPassword",
	}

	_, Err := s.EditPassword(reqDTO)

	reqDTO = jsonmodel.RequestDTOPutPassword{
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
	reqDTO := jsonmodel.RequestDTOPutPassword{
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
	if s.Remove("toto@example.dev") == nil  {
		t.Error("Delete success on user that did not exist")
	}
}
