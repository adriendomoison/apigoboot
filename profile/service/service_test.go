package service

import (
	"os"
	"testing"
	"github.com/adriendomoison/gobootapi/apicore/config"
	"github.com/adriendomoison/gobootapi/database/dbconn"
	"github.com/adriendomoison/gobootapi/profile/rest/jsonmodel"
	"github.com/adriendomoison/gobootapi/profile/repo/dbmodel"
	"github.com/adriendomoison/gobootapi/profile/repo"
	"github.com/adriendomoison/gobootapi/apicore/helpers/servicehelper"
)

var publicId string
var firstName = "John"
var lastName = "Doe"
var birthday = "1995-11-23"
var s *service

func TestMain(m *testing.M) {
	config.SetToTestingEnv()
	dbconn.Connect()
	defer dbconn.DB.Close()
	s = New(repo.New())

	dbconn.DB.AutoMigrate(&dbmodel.Entity{})

	code := m.Run()

	dbconn.DB.DropTable(&dbmodel.Entity{})

	os.Exit(code)
}

func TestCreateProfile(t *testing.T) {

	reqDTO := jsonmodel.RequestDTO{
		FirstName:         firstName,
		LastName:          lastName,
		Birthday:          birthday,
	}

	restDTO, err := s.Add(reqDTO)

	if err != nil {
		t.Fail()
	}

	publicId = restDTO.PublicId
}

func TestGetProfile(t *testing.T) {

	resDTO, err := s.Retrieve(publicId)

	if err != nil {
		t.Errorf("could not retrieve profile: %s", err.Detail)
	}

	if resDTO.PublicId == "" {
		t.Error("profile is missing public id")
	}

	if resDTO.FirstName != firstName {
		t.Errorf("expected first name to be %s but got %s", firstName, resDTO.FirstName)
	}

	if resDTO.LastName != lastName {
		t.Errorf("expected last name to be %s but got %s", lastName, resDTO.LastName)
	}

	if resDTO.Birthday != birthday {
		t.Errorf("expected birthday to be %s but got %s", birthday, resDTO.Birthday)
	}

}

func TestEditProfile(t *testing.T) {
	profileDTO := jsonmodel.RequestDTO{
		PublicId:          publicId,
		FirstName:         "Alfred",
		LastName:          "Smith",
		Birthday:          "2001-05-30",
		ProfilePictureUrl: "http://afghanjustice.org/uploads/new/img2-1425146355.png",
	}

	resDTO, err := s.Edit(profileDTO)

	if err != nil {
		t.Errorf("could not retrieve profile: %s", err.Detail)
	}

	if resDTO.PublicId == "" {
		t.Error("profile is missing public id")
	}

	if resDTO.FirstName != "Alfred" {
		t.Errorf("expected first name to be %s but got %s", "Alfred", resDTO.FirstName)
	}

	if resDTO.LastName != "Smith" {
		t.Errorf("expected last name to be %s but got %s", "Smith", resDTO.LastName)
	}

	if resDTO.Birthday != "2001-05-30" {
		t.Errorf("expected birthday to be %s but got %s", "2001-05-30", resDTO.Birthday)
	}

	if resDTO.ProfilePictureUrl != "http://afghanjustice.org/uploads/new/img2-1425146355.png" {
		t.Errorf("expected birthday to be %s but got %s", "http://afghanjustice.org/uploads/new/img2-1425146355.png", resDTO.ProfilePictureUrl)
	}
}

func TestDeleteProfile(t *testing.T) {
	res := s.Remove(publicId)
	if res != nil {
		t.Error("profile was not deleted")
	}

	resDTO, err := s.Retrieve(publicId)
	if err.Code != servicehelper.NotFound {
		t.Error("profile is not 'not found'")
	}

	if resDTO.PublicId == publicId {
		t.Error("profile was found -> still exist")
	}
}
