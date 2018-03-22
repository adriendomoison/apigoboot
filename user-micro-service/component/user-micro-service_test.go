package main_test

import (
	"github.com/adriendomoison/apigoboot/api-tool/apitool"
	"github.com/adriendomoison/apigoboot/user-micro-service/component/user"
	"github.com/adriendomoison/apigoboot/user-micro-service/component/user/repo"
	"github.com/adriendomoison/apigoboot/user-micro-service/component/user/rest"
	"github.com/adriendomoison/apigoboot/user-micro-service/component/user/service"
	"github.com/adriendomoison/apigoboot/user-micro-service/config"
	"github.com/adriendomoison/apigoboot/user-micro-service/database/dbconn"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"testing"
)

var publicBaseUrl = config.GAppUrl + "/api/v1"
var privateBaseUrl = config.GAppUrl + "/api/private-v1"

func getAccessTokenOwnerUserIdMock(c *gin.Context) {
	accessToken := c.Param("accessToken")
	if accessToken == "XXX" {
		c.JSON(http.StatusOK, gin.H{"user_id": 1})
	} else if accessToken == "YYY" {
		c.JSON(http.StatusOK, gin.H{"user_id": 2})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{})
	}
}

func getUserProfileMock(c *gin.Context) {
	c.JSON(http.StatusOK, struct {
		PublicId  string `json:"profile_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Birthday  string `json:"birthday"`
	}{
		FirstName: "John",
		LastName:  "Doe",
		Birthday:  "1980-10-20",
		Email:     c.Param("email"),
		PublicId:  "12345678",
	})
}

func postUserProfileMock(c *gin.Context) {
	c.JSON(http.StatusCreated, struct {
		PublicId  string `json:"profile_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Birthday  string `json:"birthday"`
	}{
		FirstName: "John",
		LastName:  "Doe",
		Birthday:  "1980-10-20",
		Email:     "test00@example.dev",
		PublicId:  "12345678",
	})
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
	userComponent := user.New(rest.New(service.New(repo.New())))
	userComponent.AttachPublicAPI(router.Group("/api/v1"))
	userComponent.AttachPrivateAPI(router.Group("/api/private-v1"))

	// Add mocked other micro-services called by this service
	router.GET("/api/private-v1/access-token/:accessToken/get-owner", getAccessTokenOwnerUserIdMock)
	router.POST("/api/private-v1/profiles", postUserProfileMock)
	router.GET("/api/private-v1/profiles/:email", getUserProfileMock)

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
	email := "test00@example.dev"
	password := "mySecretPassword#123"

	// build JSON request body
	requestBody := rest.RequestDTO{
		Email:    email,
		Password: password,
	}

	// call api
	userDTO := rest.ResponseDTO{}
	resp, _ := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method:      "POST",
		URL:         publicBaseUrl + "/users/",
		ContentType: "application/x-www-form-urlencoded",
	}, requestBody, &userDTO)
	defer resp.Body.Close()

	// test response
	if resp.StatusCode != 201 {
		t.Errorf("Expected %v to be %v, got %v", "status", "201", resp.Status)
	} else if userDTO.Email != email {
		t.Errorf("Expected %v to be %v, got %v", "email", email, userDTO.Email)
	}

}

func TestGet(t *testing.T) {

	// init test variable
	email := "test00@example.dev"

	// print test variable for easy debug
	t.Log("testing with following parameters:")
	t.Log(email)

	// call api
	userDTO := rest.ResponseDTO{}
	resp, _ := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method:        "GET",
		URL:           publicBaseUrl + "/users/" + email,
		Authorization: "Bearer XXX",
	}, nil, &userDTO)
	defer resp.Body.Close()

	// test response
	if resp.StatusCode != 200 {
		t.Errorf("Expected %s to be %s, got %s", "status", "200", resp.Status)
	} else if userDTO.Email != email {
		t.Errorf("Expected %s to be %s, got %s", "email", email, userDTO.Email)
	}
}

func TestPutEmailWithEmailThatIsNotAnEmail(t *testing.T) {

	// init test variable
	email := "test00@example.dev"
	newEmail := "test01"
	password := "mySecretPassword#123"

	// build JSON request body
	requestBody := rest.RequestDTOPutEmail{
		Email:    email,
		NewEmail: newEmail,
		Password: password,
	}

	// call api
	userDTO := rest.ResponseDTO{}
	resp, apiError := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method:        "PUT",
		URL:           publicBaseUrl + "/users/" + email + "/email",
		ContentType:   "application/x-www-form-urlencoded",
		Authorization: "Bearer XXX",
	}, requestBody, &userDTO)
	defer resp.Body.Close()

	// test response
	if resp.StatusCode != 400 {
		t.Errorf("Expected %s to be %s, got %s", "status", "400", resp.Status)
	} else if userDTO.Email == newEmail {
		t.Errorf("Expected %s to be %s, got %s", "email", userDTO.Email, newEmail)
	} else if apiError.Errors == nil {
		t.Errorf("Expected %s to be %s, got %s", "an error", "returned", "nothing")
	}
}

func TestPutEmail(t *testing.T) {

	// init test variable
	email := "test00@example.dev"
	newEmail := "test01@example.dev"
	password := "mySecretPassword#123"

	// build JSON request body
	requestBody := rest.RequestDTOPutEmail{
		Email:    email,
		NewEmail: newEmail,
		Password: password,
	}

	// call api
	userDTO := rest.ResponseDTO{}
	resp, apiError := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method:        "PUT",
		URL:           publicBaseUrl + "/users/" + email + "/email",
		ContentType:   "application/x-www-form-urlencoded",
		Authorization: "Bearer XXX",
	}, requestBody, &userDTO)
	defer resp.Body.Close()

	// test response
	if resp.StatusCode != 200 {
		t.Errorf("Expected %s to be %s, got %s", "status", "200", resp.Status)
	}
	if userDTO.Email != newEmail {
		t.Errorf("Expected %s to be %s, got %s", "email", newEmail, userDTO.Email)
	}
	if len(apiError.Errors) == 1 {
		t.Errorf("Expected %v to be %v, got %v", "len(apiError.Errors)", 0, 1)
	}
}

func TestPutPasswordWithoutKnowingPassword(t *testing.T) {

	// init test variable
	email := "test01@example.dev"
	password := "WrongPassword"
	newPassword := "Password123"

	// build JSON request body
	requestBody := rest.RequestDTOPutPassword{
		Email:       email,
		NewPassword: newPassword,
		Password:    password,
	}

	// call api
	resp, apiError := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method:        "PUT",
		URL:           publicBaseUrl + "/users/" + email + "/password",
		ContentType:   "application/x-www-form-urlencoded",
		Authorization: "Bearer XXX",
	}, requestBody, &rest.ResponseDTO{})
	defer resp.Body.Close()

	// test response
	if resp.StatusCode != 400 {
		t.Errorf("Expected %v to be %v, got %v", "status", "400", resp.Status)
	}
	if len(apiError.Errors) == 0 {
		t.Errorf("Expected %v to be %v, got %v", "len(apiError.Errors)", 1, 0)
	}
}

func TestPutPassword(t *testing.T) {

	// init test variable
	email := "test01@example.dev"
	password := "mySecretPassword#123"
	newPassword := "Password123"

	// build JSON request body
	requestBody := rest.RequestDTOPutPassword{
		Email:       email,
		NewPassword: newPassword,
		Password:    password,
	}

	// call api
	resp, _ := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method:        "PUT",
		URL:           publicBaseUrl + "/users/" + email + "/password",
		ContentType:   "application/x-www-form-urlencoded",
		Authorization: "Bearer XXX",
	}, requestBody, &rest.ResponseDTO{})
	defer resp.Body.Close()

	// test response
	if resp.StatusCode != 200 {
		t.Errorf("Expected %s to be %s, got %s", "status", "200", resp.Status)
	}
}

func TestDeleteWithWrongAccessToken(t *testing.T) {

	// init test variable
	email := "test01@example.dev"

	// print test variable for easy debug
	t.Log("testing with following parameters:")
	t.Log(email)

	// call api
	resp, _ := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method:        "DELETE",
		URL:           publicBaseUrl + "/users/" + email,
		Authorization: "Bearer YYY",
	}, nil, &rest.ResponseDTO{})
	defer resp.Body.Close()

	// test response
	if resp.StatusCode != 403 {
		t.Errorf("Expected %s to be %s, got %s", "status", "403", resp.Status)
	}
}

func TestDelete(t *testing.T) {

	// init test variable
	email := "test01@example.dev"

	// print test variable for easy debug
	t.Log("testing with following parameters:")
	t.Log(email)

	// call api
	resp, _ := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method:        "DELETE",
		URL:           publicBaseUrl + "/users/" + email,
		Authorization: "Bearer XXX",
	}, nil, &rest.ResponseDTO{})
	defer resp.Body.Close()

	// test response
	if resp.StatusCode != 200 {
		t.Errorf("Expected %s to be %s, got %s", "status", "200", resp.Status)
	}
}

func TestPostWithProfileCreation(t *testing.T) {

	// init test variable
	firstName := "John"
	lastName := "Doe"
	birthday := "1990-10-20"
	email := "test00@example.dev"
	password := "mySecretPassword#123"

	// build JSON request body
	requestBody := rest.RequestDTOWithProfile{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Birthday:  birthday,
	}

	// call api
	var userDTO rest.ResponseDTO
	resp, _ := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method:      "POST",
		URL:         publicBaseUrl + "/users?create_profile=true",
		ContentType: "application/x-www-form-urlencoded",
	}, requestBody, &userDTO)
	defer resp.Body.Close()

	// test response
	if resp.StatusCode != 201 {
		t.Errorf("Expected %s to be %s, got %s", "status", "201", resp.Status)
	} else if userDTO.Email != email {
		t.Errorf("Expected %s to be %s, got %s", "email", email, userDTO.Email)
	}
}

func TestProfileWasCreated(t *testing.T) {

	// init test variable
	email := "test00@example.dev"
	firstName := "John"
	publicId := "12345678"

	// print test variable for easy debug
	t.Log("testing with following parameters:")
	t.Log(email)

	// call api
	var userWithProfile rest.ResponseDTOWithProfile
	resp, _ := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method:        "GET",
		URL:           publicBaseUrl + "/users/" + email + "?get_profile=true",
		Authorization: "Bearer YYY",
	}, nil, &userWithProfile)
	defer resp.Body.Close()

	// test response
	if resp.StatusCode != 200 {
		t.Errorf("Expected %s to be %s, got %s", "status", "200", resp.Status)
	} else if userWithProfile.Email != email {
		t.Errorf("Expected %s to be %s, got %s", "email", email, userWithProfile.Email)
	} else if userWithProfile.FirstName != firstName {
		t.Errorf("Expected %s to be %s, got %s", "frist name", firstName, userWithProfile.FirstName)
	} else if userWithProfile.PublicId != publicId {
		t.Errorf("Expected %s to be %s, got %s", "profile ID", publicId, userWithProfile.PublicId)
	}
}

func TestCheckCredentials(t *testing.T) {

	// init test variable
	email := "test00@example.dev"
	password := "mySecretPassword#123"
	authType := "password"

	// build JSON request body
	requestBody := rest.RequestDTOCheckCredentials{
		Username: email,
		Password: password,
		AuthType: authType,
	}

	// call api
	var userInfo rest.ResponseDTOUserInfo
	resp, _ := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method: "POST",
		URL:    privateBaseUrl + "/user/check-credentials",
	}, requestBody, &userInfo)
	defer resp.Body.Close()

	// test response
	if resp.StatusCode != 200 {
		t.Errorf("Expected %s to be %s, got %s", "status", "200", resp.Status)
	} else if userInfo.Email != email {
		t.Errorf("Expected %s to be %s, got %s", "email", email, userInfo.Email)
	}
}

func TestGetByEmail(t *testing.T) {

	// init test variable
	email := "test00@example.dev"
	userId := 2

	// print test variable for easy debug
	t.Log(email)

	// call api
	var userInfo rest.ResponseDTOUserInfo
	resp, _ := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method: "GET",
		URL:    privateBaseUrl + "/user/email/" + email,
	}, nil, &userInfo)
	defer resp.Body.Close()

	// test response
	if resp.StatusCode != 200 {
		t.Errorf("Expected %s to be %s, got %s", "status", "200", resp.Status)
	} else if userInfo.UserId != uint(userId) {
		t.Errorf("Expected %s to be %v, got %vR", "user id", userId, userInfo.UserId)
	}
}

func TestGetByUserId(t *testing.T) {

	// init test variable
	email := "test00@example.dev"
	userId := 2

	// print test variable for easy debug
	t.Log(email)

	// call api
	var userInfo rest.ResponseDTOUserInfo
	resp, _ := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method: "GET",
		URL:    privateBaseUrl + "/user/id/" + strconv.Itoa(userId),
	}, nil, &userInfo)
	defer resp.Body.Close()

	// test response
	if resp.StatusCode != 200 {
		t.Errorf("Expected %s to be %s, got %s", "status", "200", resp.Status)
	} else if userInfo.Email != email {
		t.Errorf("Expected %s to be %s, got %s", "email", email, userInfo.Email)
	}
}
