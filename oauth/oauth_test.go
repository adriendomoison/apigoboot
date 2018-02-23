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
	"github.com/adriendomoison/go-boot-api/oauth/repo/model"
	"github.com/adriendomoison/go-boot-api/apicore/config"
	"github.com/adriendomoison/go-boot-api/apicore/helpers/apihelper"
	"github.com/adriendomoison/go-boot-api/database/dbconn"
	userrepomodel "github.com/adriendomoison/go-boot-api/user/repo/dbmodel"
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
	dbconn.DB.DropTable(&model.Authorize{})
	dbconn.DB.DropTable(&model.Client{})
	dbconn.DB.DropTable(&model.Access{})
	dbconn.DB.DropTable(&model.Refresh{})

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
	dbconn.DB.Create(&model.Client{
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
