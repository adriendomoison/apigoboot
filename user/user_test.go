package user

import (
	"os"
	"bytes"
	"testing"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/adriendomoison/go-boot-api/user/rest/jsonmodel"
	"github.com/adriendomoison/go-boot-api/user/repo/dbmodel"
	"github.com/adriendomoison/go-boot-api/apicore/config"
	"github.com/adriendomoison/go-boot-api/database/dbconn"
	"github.com/adriendomoison/go-boot-api/apicore/helpers/apihelper"
)

var URL = config.GAppUrl + "/api/v1/users"

// Generate CORS config for router
func getCORSConfig() cors.Config {
	CORSConfig := cors.DefaultConfig()
	CORSConfig.AllowCredentials = true
	CORSConfig.AllowAllOrigins = true
	CORSConfig.AllowHeaders = []string{"*", "Origin", "Content-Type", "Authorization", "Cookie"}
	CORSConfig.AllowMethods = []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"}
	return CORSConfig
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
	apiV1 := router.Group("/api/v1")
	New().Attach(apiV1)

	go router.Run(":" + config.GPort)

	code := m.Run()

	dbconn.DB.DropTable(&dbmodel.Entity{})

	os.Exit(code)
}

func TestPost(t *testing.T) {

	values := map[string]string{
		"email":    "test@example.dev",
		"username": "test",
		"password": "mySecretPassword123#",
	}

	jsonValue, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	user := jsonmodel.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	if json.Unmarshal(body, &user) != nil {
		json.Unmarshal(body, &apiErrors)
	}

	if resp.Status != "201 Created" {
		t.Errorf("Status %s was not expected, user was not created", resp.Status)
	} else if user.Email != "test@example.dev" {
		t.Error("Email was different than what was send")
	} else if user.Username != "test" {
		t.Error("Username was different than what was send")
	}
}

func TestGet(t *testing.T) {

	req, err := http.NewRequest("GET", URL+"/test@example.dev", nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	user := jsonmodel.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	if json.Unmarshal(body, &user) != nil {
		json.Unmarshal(body, &apiErrors)
	}

	if resp.Status != "200 OK" {
		t.Errorf("Status %s was not expected, user was not retrieved", resp.Status)
	} else if user.Email != "test@example.dev" {
		t.Error("Email was different than what was send")
	} else if user.Username != "test" {
		t.Error("Username was different than what was send")
	}
}

func TestPutEmail(t *testing.T) {

	values := map[string]string{
		"email":     "test@example.dev",
		"new_email": "test2@example.dev",
		"password":  "mySecretPassword123#",
	}

	jsonValue, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("PUT", URL+"/test@example.dev"+"/email", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	user := jsonmodel.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	if json.Unmarshal(body, &user) != nil {
		json.Unmarshal(body, &apiErrors)
	}

	if resp.Status != "200 OK" {
		t.Errorf("Status %s was not expected, user was not updated", resp.Status)
	} else if user.Email != "test2@example.dev" {
		t.Error("Email was different than what was send")
	} else if user.Username != "test" {
		t.Error("Username was different than what was send")
	}
}

func TestPutEmailWithEmailThatIsNotAnEmail(t *testing.T) {

	values := map[string]string{
		"email":     "test2@example.dev",
		"new_email": "test2",
		"password":  "mySecretPassword123#",
	}

	jsonValue, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("PUT", URL+"/test2@example.dev"+"/email", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	user := &jsonmodel.ResponseDTO{}
	apiErrors := &apihelper.ApiErrors{}
	if json.Unmarshal(body, user) != nil {
		json.Unmarshal(body, apiErrors)
	}

	if resp.Status != "400 Bad Request" {
		t.Errorf("Status %s was not expected", resp.Status)
	} else if user.Email == "test2" {
		t.Error("Email was not an email and this passed!?")
	} else if apiErrors == nil {
		t.Error("No error was raised")
	}
}

func TestPutPassword(t *testing.T) {

	values := map[string]string{
		"email":        "test2@example.dev",
		"password":     "mySecretPassword123#",
		"new_password": "Password123#",
	}

	jsonValue, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("PUT", URL+"/test2@example.dev"+"/password", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	user := jsonmodel.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	if json.Unmarshal(body, &user) != nil {
		json.Unmarshal(body, &apiErrors)
	}

	if resp.Status != "200 OK" {
		t.Errorf("Status %s was not expected, user was not updated", resp.Status)
	}
}

func TestPutPasswordWithoutKnowingPassword(t *testing.T) {

	values := map[string]string{
		"email":        "test2@example.dev",
		"password":     "mySecretPassword123#",
		"new_password": "Password123#321",
	}

	jsonValue, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("PUT", URL+"/test2@example.dev"+"/password", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	user := jsonmodel.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	if json.Unmarshal(body, &user) != nil {
		json.Unmarshal(body, &apiErrors)
	}

	if resp.Status != "400 Bad Request" {
		t.Error("Was able to update password without authentication")
	}
}

func TestDelete(t *testing.T) {

	req, err := http.NewRequest("DELETE", URL+"/test2@example.dev", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	user := jsonmodel.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	if json.Unmarshal(body, &user) != nil {
		json.Unmarshal(body, &apiErrors)
	}

	if resp.Status != "200 OK" {
		t.Error(apiErrors.Errors)
	}
}
