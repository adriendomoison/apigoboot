package oauth

import (
	"os"
	"bytes"
	"testing"
	"strings"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/elithrar/simple-scrypt"
	"github.com/adriendomoison/gobootapi/oauth/repo/dbmodel"
	"github.com/adriendomoison/gobootapi/apicore/config"
	"github.com/adriendomoison/gobootapi/apicore/helpers/apihelper"
	"github.com/adriendomoison/gobootapi/database/dbconn"
	userrepomodel "github.com/adriendomoison/gobootapi/user/repo/dbmodel"
)

var URL = config.GAppUrl + "/authentication"

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
	config.SetToTestingEnv()
	dbconn.Connect()
	defer dbconn.DB.Close()

	// Init router
	router := gin.Default()
	router.Use(cors.New(getCORSConfig()))

	// Append routes to server
	New().Attach(router)

	// Start server in a routine
	go router.Run(":" + config.GPort)

	createClient()

	code := m.Run()

	dbconn.DB.DropTable(&userrepomodel.Entity{})
	dbconn.DB.DropTable(&dbmodel.Authorize{})
	dbconn.DB.DropTable(&dbmodel.Client{})
	dbconn.DB.DropTable(&dbmodel.Access{})
	dbconn.DB.DropTable(&dbmodel.Refresh{})

	os.Exit(code)
}

func createClient() {
	dbconn.DB.AutoMigrate(&userrepomodel.Entity{})
	hashedPassword, _ := scrypt.GenerateFromPassword([]byte("password123"), scrypt.DefaultParams)
	dbconn.DB.Create(&userrepomodel.Entity{
		Email:    "adrien@example.dev",
		Password: string(hashedPassword[:]),
		Username: "adrien",
	})
	dbconn.DB.Create(&dbmodel.Client{
		Id:          "apigoboot",
		Secret:      "apigoboot",
		RedirectUri: "http://api.apigoboot.dev:4200/authentication/oauth2/code",
		UserId:      1,
	})
}

func TestPasswordAuthentication(t *testing.T) {
	values := map[string]string{
		"client_id":     "apigoboot",
		"client_secret": "apigoboot",
		"method":        "password",
		"username":      "adrien@example.dev",
		"password":      "password123",
	}

	jsonValue, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", URL+"/oauth2/password", bytes.NewBuffer(jsonValue))
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
	apiErrors := apihelper.ApiErrors{}
	if json.Unmarshal(body, &access) != nil {
		json.Unmarshal(body, &apiErrors)
	}

	if resp.Status != "200 OK" {
		t.Errorf("Status %s was not expected, user was not updated", resp.Status)
	} else if access.AccessToken == "" {
		t.Error("Access token is empty")
	}
}

func TestCodeAuthentication(t *testing.T) {

	code := requestCode(t)

	// TEST USAGE OF CODE
	req, err := http.NewRequest("GET", URL+"/oauth2/code?code="+code+"&state=xyz&client_id=apigoboot&client_secret=apigoboot&parse=yes", nil)
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
	apiErrors := apihelper.ApiErrors{}
	if json.Unmarshal(body, &access) != nil {
		json.Unmarshal(body, &apiErrors)
	}

	if resp.Status != "200 OK" {
		t.Errorf("Status %s was not expected, user was not updated", resp.Status)
	} else if access.AccessToken == "" {
		t.Error("Access token is empty")
	}
}

func requestCode(t *testing.T) string {
	form := url.Values{}
	form.Add("username", "adrien@example.dev")
	form.Add("password", "password123")

	req, err := http.NewRequest("POST", URL+"/authorize?response_type=code&client_id=apigoboot&client_secret=apigoboot&state=xyz&scope=everything&redirect_uri=http://api.apigoboot.dev:4200/authentication/oauth2/code", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	access := struct {
		Code string `json:"code"`
	}{}
	apiErrors := apihelper.ApiErrors{}
	if json.Unmarshal(body, &access) != nil {
		json.Unmarshal(body, &apiErrors)
	}

	if resp.Status != "200 OK" {
		t.Errorf("Status %s was not expected, user was not updated", resp.Status)
	} else if access.Code == "" {
		t.Error("Authorization code is empty")
	}
	return access.Code
}
