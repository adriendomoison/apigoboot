package main_test

import (
	"io"
	"os"
	"bytes"
	"testing"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/adriendomoison/gobootapi/errorhandling/apihelper"
	"github.com/adriendomoison/gobootapi/user-micro-service/config"
	"github.com/adriendomoison/gobootapi/user-micro-service/usercomponent"
	"github.com/adriendomoison/gobootapi/user-micro-service/usercomponent/rest"
	"github.com/adriendomoison/gobootapi/user-micro-service/usercomponent/service"
	"github.com/adriendomoison/gobootapi/user-micro-service/usercomponent/repo"
	"github.com/adriendomoison/gobootapi/user-micro-service/database/dbconn"
)

var BaseUrl = config.GAppUrl + "/api/v1"

// Generate CORS config for router
func getCORSConfig() cors.Config {
	CORSConfig := cors.DefaultConfig()
	CORSConfig.AllowCredentials = true
	CORSConfig.AllowAllOrigins = true
	CORSConfig.AllowHeaders = []string{"*", "Origin", "Content-Type", "Authorization", "Cookie"}
	CORSConfig.AllowMethods = []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"}
	return CORSConfig
}

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
		Email:     "test00@example.dev",
		PublicId:  "12345678",
	})
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
	userComponent := usercomponent.New(rest.New(service.New(repo.New())))
	userComponent.AttachPublicAPI(router.Group("/api/v1"))

	// Add mocked other micro-services called by this service
	router.GET("/api/private-v1/access-token/:accessToken/get-owner", getAccessTokenOwnerUserIdMock)
	router.POST("/api/private-v1/profiles", postUserProfileMock)
	router.GET("/api/private-v1/profiles/:email", getUserProfileMock)

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
	email := "test00@example.dev"
	password := "mySecretPassword#123"

	// build JSON request body
	requestBody := rest.RequestDTO{
		Email:    email,
		Password: password,
	}

	// call api
	req, err := http.NewRequest("POST", BaseUrl+"/users", encodeRequestBody(t, requestBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	user := rest.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	json.Unmarshal(body, &user)
	json.Unmarshal(body, &apiErrors)

	t.Log(user)
	t.Log(apiErrors)

	// test response
	if resp.StatusCode != 201 {
		t.Errorf("Expected %s to be %s, got %s", "status", "201", resp.Status)
	} else if user.Email != email {
		t.Errorf("Expected %s to be %s, got %s", "email", email, user.Email)
	}

}

func TestGet(t *testing.T) {

	// init test variable
	email := "test00@example.dev"

	// print test variable for easy debug
	t.Log("testing with following parameters:")
	t.Log(email)

	// call api
	req, err := http.NewRequest("GET", BaseUrl+"/users/"+email, nil)
	req.Header.Set("Authorization", "Bearer XXX")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	user := rest.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	json.Unmarshal(body, &user)
	json.Unmarshal(body, &apiErrors)

	t.Log(user)
	t.Log(apiErrors)

	// test response
	if resp.StatusCode != 200 {
		t.Errorf("Expected %s to be %s, got %s", "status", "200", resp.Status)
	} else if user.Email != email {
		t.Errorf("Expected %s to be %s, got %s", "email", email, user.Email)
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
	req, err := http.NewRequest("PUT", BaseUrl+"/users/"+email+"/email", encodeRequestBody(t, requestBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer XXX")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	user := rest.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	json.Unmarshal(body, &user)
	json.Unmarshal(body, &apiErrors)

	t.Log(user)
	t.Log(apiErrors)

	// test response
	if resp.StatusCode != 400 {
		t.Errorf("Expected %s to be %s, got %s", "status", "400", resp.Status)
	} else if user.Email == newEmail {
		t.Errorf("Expected %s to be %s, got %s", "email", user.Email, newEmail)
	} else if apiErrors.Errors == nil {
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
	req, err := http.NewRequest("PUT", BaseUrl+"/users/"+email+"/email", encodeRequestBody(t, requestBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer XXX")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	user := rest.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	json.Unmarshal(body, &user)
	json.Unmarshal(body, &apiErrors)

	t.Log(user)
	t.Log(apiErrors)

	// test response
	if resp.StatusCode != 200 {
		t.Errorf("Expected %s to be %s, got %s", "status", "200", resp.Status)
	} else if user.Email != newEmail {
		t.Errorf("Expected %s to be %s, got %s", "email", newEmail, user.Email)
	}
}

func TestPutPasswordWithoutKnowingPassword(t *testing.T) {

	// init test variable
	email := "test01@example.dev"
	password := "WrongPassword"
	newPassword := "Password123"

	// build JSON request body
	requestBody := rest.RequestDTOPutPassword{
		Email:    email,
		NewPassword: newPassword,
		Password: password,
	}

	// call api
	req, err := http.NewRequest("PUT", BaseUrl+"/users/"+email+"/password", encodeRequestBody(t, requestBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer XXX")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	user := rest.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	json.Unmarshal(body, &user)
	json.Unmarshal(body, &apiErrors)

	t.Log(user)
	t.Log(apiErrors)

	// test response
	if resp.StatusCode != 400 {
		t.Errorf("Expected %s to be %s, got %s", "status", "400", resp.Status)
	}
}

func TestPutPassword(t *testing.T) {

	// init test variable
	email := "test01@example.dev"
	password := "mySecretPassword#123"
	newPassword := "Password123"

	// build JSON request body
	requestBody := rest.RequestDTOPutPassword{
		Email:    email,
		NewPassword: newPassword,
		Password: password,
	}

	// call api
	req, err := http.NewRequest("PUT", BaseUrl+"/users/"+email+"/password", encodeRequestBody(t, requestBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer XXX")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	user := rest.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	json.Unmarshal(body, &user)
	json.Unmarshal(body, &apiErrors)

	t.Log(user)
	t.Log(apiErrors)

	// test response
	if resp.StatusCode != 200 {
		t.Errorf("Expected %s to be %s, got %s", "status", "200", resp.Status)
	}
}

func TestDeleteWithWringAccessToken(t *testing.T) {

	// init test variable
	email := "test01@example.dev"

	// print test variable for easy debug
	t.Log("testing with following parameters:")
	t.Log(email)

	// call api
	req, err := http.NewRequest("DELETE", BaseUrl+"/users/"+email, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer YYY")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	user := rest.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	json.Unmarshal(body, &user)
	json.Unmarshal(body, &apiErrors)

	t.Log(user)
	t.Log(apiErrors)

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
	req, err := http.NewRequest("DELETE", BaseUrl+"/users/"+email, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer XXX")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	user := rest.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	json.Unmarshal(body, &user)
	json.Unmarshal(body, &apiErrors)

	t.Log(user)
	t.Log(apiErrors)

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
		Email:    email,
		Password: password,
		FirstName: firstName,
		LastName: lastName,
		Birthday: birthday,
	}

	// call api
	req, err := http.NewRequest("POST", BaseUrl+"/users?createprofile=true", encodeRequestBody(t, requestBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	user := rest.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	json.Unmarshal(body, &user)
	json.Unmarshal(body, &apiErrors)

	t.Log(user)
	t.Log(apiErrors)

	// test response
	if resp.StatusCode != 201 {
		t.Errorf("Expected %s to be %s, got %s", "status", "201", resp.Status)
	} else if user.Email != email {
		t.Errorf("Expected %s to be %s, got %s", "email", email, user.Email)
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
	req, err := http.NewRequest("GET", BaseUrl+"/users/"+email+"?getprofile=true", nil)
	req.Header.Set("Authorization", "Bearer YYY")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	userWithProfile := rest.ResponseDTOWithProfile{}
	apiErrors := apihelper.ApiErrors{}
	json.Unmarshal(body, &userWithProfile)
	json.Unmarshal(body, &apiErrors)

	t.Log(userWithProfile)
	t.Log(apiErrors)

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
