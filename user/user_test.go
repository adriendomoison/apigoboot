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
	"github.com/adriendomoison/gobootapi/oauth"
	"github.com/adriendomoison/gobootapi/apicore/config"
	"github.com/adriendomoison/gobootapi/database/dbconn"
	"github.com/adriendomoison/gobootapi/user/rest/jsonmodel"
	"github.com/adriendomoison/gobootapi/user/repo/dbmodel"
	"github.com/adriendomoison/gobootapi/apicore/helpers/apihelper"
	oauthdbmodel "github.com/adriendomoison/gobootapi/oauth/repo/dbmodel"
	profilejsonmodel "github.com/adriendomoison/gobootapi/profile/rest/jsonmodel"
)

var URL = config.GAppUrl + "/api/v1/users"
var URLOAuth = config.GAppUrl + "/authentication"

// Generate CORS config for router
func getCORSConfig() cors.Config {
	CORSConfig := cors.DefaultConfig()
	CORSConfig.AllowCredentials = true
	CORSConfig.AllowAllOrigins = true
	CORSConfig.AllowHeaders = []string{"*", "Origin", "Content-Type", "Authorization", "Cookie"}
	CORSConfig.AllowMethods = []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"}
	return CORSConfig
}

func createClient() {
	dbconn.DB.AutoMigrate(&oauthdbmodel.Access{})
	dbconn.DB.AutoMigrate(&oauthdbmodel.Client{})
	dbconn.DB.AutoMigrate(&oauthdbmodel.Authorize{})
	dbconn.DB.AutoMigrate(&oauthdbmodel.Refresh{})
	dbconn.DB.Create(&oauthdbmodel.Client{
		Id:          "apigoboot",
		Secret:      "apigoboot",
		RedirectUri: "http://api.apigoboot.dev:4200/authentication/oauth2/code",
		UserId:      1,
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
	router.Use(cors.New(getCORSConfig()))

	// Append routes to server
	apiV1 := router.Group("/api/v1")
	oauth.New().Attach(router)
	New().Attach(apiV1)

	createClient()

	go router.Run(":" + config.GPort)

	code := m.Run()

	dbconn.DB.DropTable(&dbmodel.Entity{})
	dbconn.DB.DropTable(&oauthdbmodel.Authorize{})
	dbconn.DB.DropTable(&oauthdbmodel.Client{})
	dbconn.DB.DropTable(&oauthdbmodel.Access{})
	dbconn.DB.DropTable(&oauthdbmodel.Refresh{})

	os.Exit(code)
}

func TestPost(t *testing.T) {

	jsonValue, err := json.Marshal(map[string]string{
		"email":      "test@example.dev",
		"password":   "mySecretPassword#123",
		"first_name": "John",
		"last_name":  "Doe",
		"birthday":   "1990-10-20",
	})
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

	profile := profilejsonmodel.ResponseDTO{}
	apiErrors := apihelper.ApiErrors{}
	if json.Unmarshal(body, &profile) != nil {
		json.Unmarshal(body, &apiErrors)
		t.Errorf("Did not expect any error at this point: %s", apiErrors)
	}

	if resp.Status != "201 Created" {
		t.Errorf("Status %s was not expected, user was not created", resp.Status)
	} else if profile.Email != "test@example.dev" {
		t.Errorf("Email was different than what was send: %s", profile.Email)
	}
}

func TestProfileWasCreated(t *testing.T) {

	accessToken := getAccessToken(map[string]string{
		"client_id":     "apigoboot",
		"client_secret": "apigoboot",
		"method":        "password",
		"username":      "test@example.dev",
		"password":      "mySecretPassword#123",
	})

	if len(accessToken) == 0 {
		t.Error("No access token")
		return
	}

	req, err := http.NewRequest("GET", URL+"/test@example.dev", nil)
	req.Header.Set("authorization", "Bearer "+accessToken)

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
	}
}

func TestGet(t *testing.T) {

	accessToken := getAccessToken(map[string]string{
		"client_id":     "apigoboot",
		"client_secret": "apigoboot",
		"method":        "password",
		"username":      "test@example.dev",
		"password":      "mySecretPassword#123",
	})

	if len(accessToken) == 0 {
		t.Error("No access token")
		return
	}

	req, err := http.NewRequest("GET", URL+"/test@example.dev", nil)
	req.Header.Set("authorization", "Bearer "+accessToken)

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
	}
}

func TestPutEmail(t *testing.T) {

	accessToken := getAccessToken(map[string]string{
		"client_id":     "apigoboot",
		"client_secret": "apigoboot",
		"method":        "password",
		"username":      "test@example.dev",
		"password":      "mySecretPassword#123",
	})

	if len(accessToken) == 0 {
		t.Error("No access token")
		return
	}

	jsonValue, err := json.Marshal(map[string]string{
		"client_id":     "apigoboot",
		"client_secret": "apigoboot",
		"method":        "password",
		"email":         "test@example.dev",
		"new_email":     "test2@example.dev",
		"password":      "mySecretPassword#123",
	})

	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("PUT", URL+"/test@example.dev"+"/email", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("authorization", "Bearer "+accessToken)

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
	}
}

func TestPutEmailWithEmailThatIsNotAnEmail(t *testing.T) {

	accessToken := getAccessToken(map[string]string{
		"client_id":     "apigoboot",
		"client_secret": "apigoboot",
		"method":        "password",
		"username":      "test2@example.dev",
		"password":      "mySecretPassword#123",
	})

	if len(accessToken) == 0 {
		t.Error("No access token")
		return
	}

	jsonValue, err := json.Marshal(map[string]string{
		"email":     "test2@example.dev",
		"new_email": "test2",
		"password":  "mySecretPassword#123",
	})

	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("PUT", URL+"/test2@example.dev"+"/email", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("authorization", "Bearer "+accessToken)

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

	accessToken := getAccessToken(map[string]string{
		"client_id":     "apigoboot",
		"client_secret": "apigoboot",
		"method":        "password",
		"username":      "test2@example.dev",
		"password":      "mySecretPassword#123",
	})

	if len(accessToken) == 0 {
		t.Error("No access token")
		return
	}

	jsonValue, err := json.Marshal(map[string]string{
		"email":        "test2@example.dev",
		"new_password": "Password#123",
		"password":     "mySecretPassword#123",
	})
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("PUT", URL+"/test2@example.dev"+"/password", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("authorization", "Bearer "+accessToken)

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

	accessToken := getAccessToken(map[string]string{
		"client_id":     "apigoboot",
		"client_secret": "apigoboot",
		"method":        "password",
		"username":      "test2@example.dev",
		"password":      "Password#123",
	})

	if len(accessToken) == 0 {
		t.Error("No access token")
		return
	}

	jsonValue, err := json.Marshal(map[string]string{
		"email":        "test2@example.dev",
		"password":     "password#123",
		"new_password": "Password#123321",
	})
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("PUT", URL+"/test2@example.dev"+"/password", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("authorization", "Bearer "+accessToken)

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

	accessToken := getAccessToken(map[string]string{
		"client_id":     "apigoboot",
		"client_secret": "apigoboot",
		"method":        "password",
		"username":      "test2@example.dev",
		"password":      "Password#123",
	})

	if len(accessToken) == 0 {
		t.Error("No access token")
		return
	}

	req, err := http.NewRequest("DELETE", URL+"/test2@example.dev", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("authorization", "Bearer "+accessToken)

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

func getAccessToken(values map[string]string) string {
	jsonValue, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", URLOAuth+"/oauth2/password", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	access := struct {
		AccessToken string `json:"access_token"`
	}{}
	json.Unmarshal(body, &access)

	return access.AccessToken
}
