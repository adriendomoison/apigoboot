package main_test

import (
	"github.com/adriendomoison/apigoboot/api-tool/apitool"
	"github.com/adriendomoison/apigoboot/profile-micro-service/component/profile"
	"github.com/adriendomoison/apigoboot/profile-micro-service/component/profile/repo"
	"github.com/adriendomoison/apigoboot/profile-micro-service/component/profile/rest"
	"github.com/adriendomoison/apigoboot/profile-micro-service/component/profile/service"
	"github.com/adriendomoison/apigoboot/profile-micro-service/config"
	"github.com/adriendomoison/apigoboot/profile-micro-service/database/dbconn"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"testing"
)

var publicBaseUrl = config.GAppUrl + "/api/v1"
var privateBaseUrl = config.GAppUrl + "/api/private-v1"
var profilePublicId = ""

func getAccessTokenOwnerUserIdMock(c *gin.Context) {
	accessToken := c.Param("accessToken")
	if accessToken == "XXX" {
		c.JSON(http.StatusOK, gin.H{"user_id": 1})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{})
	}
}

func getUserById(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "1" {
		c.JSON(http.StatusOK, rest.ResponseDTOUserInfo{
			Email:  "test00@example.dev",
			UserId: 1,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{})
	}
}

func getUserByEmail(c *gin.Context) {
	email := c.Param("email")
	if email == "test00@example.dev" {
		c.JSON(http.StatusOK, rest.ResponseDTOUserInfo{
			Email:  "test00@example.dev",
			UserId: 1,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{})
	}
}

func TestMain(m *testing.M) {
	// Init Env
	config.SetToTestingEnv()

	// Init DB
	dbconn.Connect()
	defer dbconn.DB.Close()

	// Init router
	router := gin.Default()
	router.Use(cors.New(apitool.DefaultCORSConfig()))

	// Append routes to server
	profileComponent := profile.New(rest.New(service.New(repo.New())))
	profileComponent.AttachPublicAPI(router.Group("/api/v1"))
	profileComponent.AttachPrivateAPI(router.Group("/api/private-v1"))

	// Add mocked other micro-services called by this service
	router.GET("/api/private-v1/access-token/:accessToken/get-owner", getAccessTokenOwnerUserIdMock)
	router.GET("/api/private-v1/user/id/:userId", getUserById)
	router.GET("/api/private-v1/user/email/:email", getUserByEmail)

	// Start service
	go router.Run(":" + config.GPort)

	// Wait and check if the http server is running
	apitool.WaitForServerToStart(publicBaseUrl + "/")

	// Start tests
	code := m.Run()

	// Drop test tables
	dbconn.DB.DropTable(&service.Entity{})

	// Stop tests
	os.Exit(code)
}

func TestPost(t *testing.T) {

	// init test variable
	firstName := "John"
	lastName := "Doe"
	birthday := "1980-10-20"
	email := "test00@example.dev"

	// build JSON request body
	requestBody := rest.RequestDTOCreation{
		FirstName: firstName,
		LastName:  lastName,
		Birthday:  birthday,
		Email:     email,
	}

	// call api
	var profileDTO rest.ResponseDTO
	resp, _ := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method:        "POST",
		URL:           privateBaseUrl + "/profiles",
		ContentType:   "application/json",
		Authorization: "Bearer XXX",
	}, requestBody, &profileDTO)
	defer resp.Body.Close()

	// test response
	if resp.StatusCode != 201 {
		t.Errorf("Expected %s to be %s, got %s", "status", "201", resp.Status)
	} else if profileDTO.Email != email {
		t.Errorf("Expected %s to be %s, got %s", "email", email, profileDTO.Email)
	} else if profileDTO.FirstName != firstName {
		t.Errorf("Expected %s to be %s, got %s", "first name", firstName, profileDTO.FirstName)
	} else if profileDTO.Birthday != birthday {
		t.Errorf("Expected %s to be %s, got %s", "birthday", birthday, profileDTO.Birthday)
	} else if profileDTO.ProfilePictureUrl == "" {
		t.Errorf("Expected %s to be %s, got %s", "profile pircture", "filled with default value", "nothing")
	} else if profileDTO.PublicId == "" {
		t.Errorf("Expected %s to be %s, got %s", "profile public id", "generated", "nothing")
	}
	profilePublicId = profileDTO.PublicId
}

func TestGet(t *testing.T) {

	// init test variable
	publicId := profilePublicId
	birthday := "1980-10-20"
	email := "test00@example.dev"

	// print test variable for easy debug
	t.Log("testing with following parameters:")
	t.Log(publicId)

	// call api
	var profileDTO rest.ResponseDTO
	resp, _ := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method:        "GET",
		URL:           publicBaseUrl + "/profiles/" + publicId,
		Authorization: "Bearer XXX",
	}, nil, &profileDTO)
	defer resp.Body.Close()

	// test response
	if resp.StatusCode != 200 {
		t.Errorf("Expected %s to be %s, got %s", "status", "200", resp.Status)
	} else if profileDTO.Email != email {
		t.Errorf("Expected %s to be %s, got %s", "email", email, profileDTO.Email)
	} else if profileDTO.Birthday != birthday {
		t.Errorf("Expected %s to be %s, got %s", "first name", birthday, profileDTO.Birthday)
	}
}

func TestPut(t *testing.T) {

	// init test variable
	publicId := profilePublicId
	email := "test00@example.dev"
	firstName := "Johnny"
	lastName := "Blop"
	birthday := "1981-11-21"
	profilePictureUrl := "http://new.picture.dev/pic"

	// build JSON request body
	requestBody := rest.RequestDTO{
		PublicId:          publicId,
		Birthday:          birthday,
		ProfilePictureUrl: profilePictureUrl,
		FirstName:         firstName,
		LastName:          lastName,
	}

	// call api
	var profileDTO rest.ResponseDTO
	resp, _ := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method:        "PUT",
		URL:           publicBaseUrl + "/profiles/" + publicId,
		ContentType:   "application/json",
		Authorization: "Bearer XXX",
	}, requestBody, &profileDTO)
	defer resp.Body.Close()

	// test response
	if resp.StatusCode != 200 {
		t.Errorf("Expected %s to be %s, got %s", "status", "200", resp.Status)
	} else if profileDTO.Email != email {
		t.Errorf("Expected %s to be %s, got %s", "email", email, profileDTO.Email)
	} else if profileDTO.FirstName != firstName {
		t.Errorf("Expected %s to be %s, got %s", "first name", firstName, profileDTO.FirstName)
	} else if profileDTO.Birthday != birthday {
		t.Errorf("Expected %s to be %s, got %s", "birthday", birthday, profileDTO.Birthday)
	}
}

func TestPutMissingParam(t *testing.T) {

	// init test variable
	publicId := profilePublicId

	// build JSON request body
	requestBody := rest.RequestDTO{
		PublicId: publicId,
	}

	// call api
	var profileDTO rest.ResponseDTO
	resp, apiError := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method:        "PUT",
		URL:           publicBaseUrl + "/profiles/" + publicId,
		ContentType:   "application/json",
		Authorization: "Bearer XXX",
	}, requestBody, &profileDTO)
	defer resp.Body.Close()

	// test response
	if resp.StatusCode != 400 {
		t.Errorf("Expected %s to be %s, got %s", "status", "400", resp.Status)
	} else if len(apiError.Errors) == 0 {
		t.Errorf("Expected %s to be %s, got %s", "apiError.Errors", "containing errors", "no errors")
	} else if len(apiError.Errors) != 4 {
		t.Errorf("Expected %s to be %s, got %s", "apiError.Errors", "containing 4 errors", strconv.Itoa(len(apiError.Errors))+" errors")
	}
}

func TestDelete(t *testing.T) {

	// init test variable
	publicId := profilePublicId

	// print test variable for easy debug
	t.Log("testing with following parameters:")
	t.Log(publicId)

	// call api
	var profileDTO rest.ResponseDTO
	resp, _ := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method:        "DELETE",
		URL:           privateBaseUrl + "/profiles/" + publicId,
		ContentType:   "application/json",
		Authorization: "Bearer XXX",
	}, nil, &profileDTO)
	defer resp.Body.Close()

	// test response
	if resp.StatusCode != 200 {
		t.Errorf("Expected %s to be %s, got %s", "status", "200", resp.Status)
	}
}
