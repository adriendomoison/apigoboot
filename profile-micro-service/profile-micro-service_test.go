package main_test

import (
	"io"
	"os"
	"bytes"
	"strconv"
	"testing"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/adriendomoison/gobootapi/errorhandling/apihelper"
	"github.com/adriendomoison/gobootapi/profile-micro-service/config"
	"github.com/adriendomoison/gobootapi/profile-micro-service/database/dbconn"
	"github.com/adriendomoison/gobootapi/profile-micro-service/profilecomponent/rest"
	"github.com/adriendomoison/gobootapi/profile-micro-service/profilecomponent/service"
	"github.com/adriendomoison/gobootapi/profile-micro-service/profilecomponent/repo"
	"github.com/adriendomoison/gobootapi/profile-micro-service/profilecomponent"
)

var PublicBaseUrl = config.GAppUrl + "/api/v1"
var PrivateBaseUrl = config.GAppUrl + "/api/private-v1"
var ProfilePublicId = ""

// Generate CORS config for router
func getCORSConfig() cors.Config {
	CORSConfig := cors.DefaultConfig()
	CORSConfig.AllowCredentials = true
	CORSConfig.AllowAllOrigins = true
	CORSConfig.AllowHeaders = []string{"*"}
	CORSConfig.AllowMethods = []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"}
	return CORSConfig
}

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

func encodeRequestBody(t *testing.T, reqBody interface{}) io.Reader {
	t.Log("testing with following parameters:")
	t.Log(reqBody)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqBody)
	return b
}

func TestMain(m *testing.M) {
	// Init Env
	config.SetToTestingEnv()

	// Init DB
	dbconn.Connect()
	defer dbconn.DB.Close()

	// Init router
	router := gin.Default()
	router.Use(cors.New(getCORSConfig()))

	// Append routes to server
	profileComponent := profilecomponent.New(rest.New(service.New(repo.New())))
	profileComponent.AttachPublicAPI(router.Group("/api/v1"))
	profileComponent.AttachPrivateAPI(router.Group("/api/private-v1"))

	// Add mocked other micro-services called by this service
	router.GET("/api/private-v1/access-token/:accessToken/get-owner", getAccessTokenOwnerUserIdMock)
	router.GET("/api/private-v1/user/id/:userId", getUserById)
	router.GET("/api/private-v1/user/email/:email", getUserByEmail)

	// Start service
	go router.Run(":" + config.GPort)

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
	req, err := http.NewRequest("POST", PrivateBaseUrl+"/profiles", encodeRequestBody(t, requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer XXX")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	profile := rest.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	json.Unmarshal(body, &profile)
	json.Unmarshal(body, &apiErrors)

	t.Log(profile)
	t.Log(apiErrors)

	// test response
	if resp.StatusCode != 201 {
		t.Errorf("Expected %s to be %s, got %s", "status", "201", resp.Status)
	} else if profile.Email != email {
		t.Errorf("Expected %s to be %s, got %s", "email", email, profile.Email)
	} else if profile.FirstName != firstName {
		t.Errorf("Expected %s to be %s, got %s", "first name", firstName, profile.FirstName)
	} else if profile.Birthday != birthday {
		t.Errorf("Expected %s to be %s, got %s", "birthday", birthday, profile.Birthday)
	} else if profile.ProfilePictureUrl == "" {
		t.Errorf("Expected %s to be %s, got %s", "profile pircture", "filled with default value", "nothing")
	} else if profile.PublicId == "" {
		t.Errorf("Expected %s to be %s, got %s", "profile public id", "generated", "nothing")
	}
	ProfilePublicId = profile.PublicId
}

func TestGet(t *testing.T) {

	// init test variable
	publicId := ProfilePublicId
	birthday := "1980-10-20"
	email := "test00@example.dev"

	// print test variable for easy debug
	t.Log("testing with following parameters:")
	t.Log(publicId)

	// call api
	req, err := http.NewRequest("GET", PublicBaseUrl+"/profiles/"+publicId, nil)
	req.Header.Set("Authorization", "Bearer XXX")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	profile := rest.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	json.Unmarshal(body, &profile)
	json.Unmarshal(body, &apiErrors)

	t.Log(profile)
	t.Log(apiErrors)

	// test response
	if resp.StatusCode != 200 {
		t.Errorf("Expected %s to be %s, got %s", "status", "200", resp.Status)
	} else if profile.Email != email {
		t.Errorf("Expected %s to be %s, got %s", "email", email, profile.Email)
	} else if profile.Birthday != birthday {
		t.Errorf("Expected %s to be %s, got %s", "first name", birthday, profile.Birthday)
	}
}

func TestPut(t *testing.T) {

	// init test variable
	publicId := ProfilePublicId
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
	req, err := http.NewRequest("PUT", PublicBaseUrl+"/profiles/"+publicId, encodeRequestBody(t, requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer XXX")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	profile := rest.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	json.Unmarshal(body, &profile)
	json.Unmarshal(body, &apiErrors)

	t.Log(profile)
	t.Log(apiErrors)

	// test response
	if resp.StatusCode != 200 {
		t.Errorf("Expected %s to be %s, got %s", "status", "200", resp.Status)
	} else if profile.Email != email {
		t.Errorf("Expected %s to be %s, got %s", "email", email, profile.Email)
	} else if profile.FirstName != firstName {
		t.Errorf("Expected %s to be %s, got %s", "first name", firstName, profile.FirstName)
	} else if profile.Birthday != birthday {
		t.Errorf("Expected %s to be %s, got %s", "birthday", birthday, profile.Birthday)
	}
}

func TestPutMissingParam(t *testing.T) {

	// init test variable
	publicId := ProfilePublicId

	// build JSON request body
	requestBody := rest.RequestDTO{
		PublicId:          publicId,
	}

	// call api
	req, err := http.NewRequest("PUT", PublicBaseUrl+"/profiles/"+publicId, encodeRequestBody(t, requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer XXX")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	profile := rest.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	json.Unmarshal(body, &profile)
	json.Unmarshal(body, &apiErrors)

	t.Log(profile)
	t.Log(apiErrors)

	// test response
	if resp.StatusCode != 400 {
		t.Errorf("Expected %s to be %s, got %s", "status", "400", resp.Status)
	} else if len(apiErrors.Errors) == 0 {
		t.Errorf("Expected %s to be %s, got %s", "apiErrors.Errors", "containing errors", "no errors")
	} else if len(apiErrors.Errors) != 4 {
		t.Errorf("Expected %s to be %s, got %s", "apiErrors.Errors", "containing 4 errors", strconv.Itoa(len(apiErrors.Errors)) + " errors")
	}
}

func TestDelete(t *testing.T) {

	// init test variable
	publicId := ProfilePublicId

	// print test variable for easy debug
	t.Log("testing with following parameters:")
	t.Log(publicId)

	// call api
	req, err := http.NewRequest("DELETE", PrivateBaseUrl+"/profiles/"+publicId, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer XXX")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	profile := rest.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	json.Unmarshal(body, &profile)
	json.Unmarshal(body, &apiErrors)

	t.Log(profile)
	t.Log(apiErrors)

	// test response
	if resp.StatusCode != 200 {
		t.Errorf("Expected %s to be %s, got %s", "status", "200", resp.Status)
	}
}
