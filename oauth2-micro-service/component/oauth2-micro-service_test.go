package main_test

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
	"github.com/adriendomoison/apigoboot/errorhandling/apihelper"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/config"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/database/dbconn"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/component/oauth2"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/component/oauth2/rest"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/component/oauth2/repo"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/component/oauth2/service"
	"github.com/RangelReale/osin"
)

var PublicBaseUrl = config.GAppUrl + "/authentication"
var PrivateBaseUrl = config.GAppUrl + "/api/private-v1/authentication"
var AccessToken = ""

// getCORSConfig Generate CORS config for router
func getCORSConfig() cors.Config {
	CORSConfig := cors.DefaultConfig()
	CORSConfig.AllowCredentials = true
	CORSConfig.AllowAllOrigins = true
	CORSConfig.AllowHeaders = []string{"*"}
	CORSConfig.AllowMethods = []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"}
	return CORSConfig
}

// initOAuthServer Init OSIN OAuth server
func initOAuthServer() *osin.Server {
	serverConfig := osin.NewServerConfig()
	serverConfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}
	serverConfig.AllowedAccessTypes = osin.AllowedAccessType{osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN, osin.PASSWORD, osin.CLIENT_CREDENTIALS, osin.ASSERTION}
	serverConfig.AllowGetAccessRequest = true
	serverConfig.AllowClientSecretInParams = true
	return osin.NewServer(serverConfig, repo.NewStorage(dbconn.DB.DB()))
}

func requestCode(t *testing.T) string {
	form := url.Values{}
	form.Add("username", "adrien@example.dev")
	form.Add("password", "password123")

	req, err := http.NewRequest("POST", PublicBaseUrl+"/authorize?response_type=code&client_id=apigoboot&client_secret=apigoboot&state=xyz&scope=everything&redirect_uri=http://api.go.boot:4202/authentication/oauth2/code", strings.NewReader(form.Encode()))
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

func CheckCredentialsMock(c *gin.Context) {
	var reqDTO rest.RequestDTOUserCredentials
	if err := c.BindJSON(&reqDTO); err != nil {
		c.JSON(apihelper.BuildRequestError(err))
	} else {
		c.JSON(http.StatusOK, rest.ResponseDTOUserInfo{
			UserId: 1,
			Email: "test00@example.dev",
		})
	}
	return
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
	oauth2Component := oauth2.New(rest.New(initOAuthServer(), service.New(repo.New())))
	oauth2Component.AttachPublicAPI(router.Group("/authentication"))
	oauth2Component.AttachPrivateAPI(router.Group("/api/private-v1/authentication"))

	// Add mocked other micro-services called by this service
	router.POST("/api/private-v1/user/check-credentials", CheckCredentialsMock)

	// Start server in a routine
	go router.Run(":" + config.GPort)

	// Set up items in DB
	createClient()

	// Start tests
	code := m.Run()

	// Drop test tables
	dbconn.DB.DropTable(&service.Authorize{})
	dbconn.DB.DropTable(&service.Client{})
	dbconn.DB.DropTable(&service.Access{})
	dbconn.DB.DropTable(&service.Refresh{})

	// Stop tests
	os.Exit(code)
}

func createClient() {
	dbconn.DB.Create(&service.Client{
		Id:          "apigoboot",
		Secret:      "apigoboot",
		RedirectUri: "http://api.go.boot:4202/authentication/oauth2/code",
		UserId:      1,
	})
}

func TestPasswordAuthentication(t *testing.T) {

	// init test variable
	values := map[string]string{
		"client_id":     "apigoboot",
		"client_secret": "apigoboot",
		"method":        "password",
		"username":      "test00@example.dev",
		"password":      "password123",
	}

	jsonValue, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}

	// call api
	req, err := http.NewRequest("POST", PublicBaseUrl+"/oauth2/password", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	access := struct {
		AccessToken string `json:"access_token"`
	}{}
	apiErrors := apihelper.ApiErrors{}
	if json.Unmarshal(body, &access) != nil {
		json.Unmarshal(body, &apiErrors)
	}

	// test response
	if resp.Status != "200 OK" {
		t.Errorf("Status %s was not expected, user was not updated", resp.Status)
	} else if access.AccessToken == "" {
		t.Error("Access token is empty")
	}

	AccessToken = access.AccessToken
}

func TestCodeAuthentication(t *testing.T) {

	// init test variable
	code := requestCode(t)

	// call api
	req, err := http.NewRequest("GET", PublicBaseUrl+"/oauth2/code?code="+code+"&state=xyz&client_id=apigoboot&client_secret=apigoboot&parse=yes", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	access := struct {
		AccessToken string `json:"access_token"`
	}{}
	apiErrors := apihelper.ApiErrors{}
	if json.Unmarshal(body, &access) != nil {
		json.Unmarshal(body, &apiErrors)
	}

	// test response
	if resp.Status != "200 OK" {
		t.Errorf("Status %s was not expected, user was not updated", resp.Status)
	} else if access.AccessToken == "" {
		t.Error("Access token is empty")
	}
}

func TestGetAccessTokenOwnerUserId(t *testing.T) {

	// init test variable
	userId := 1

	// print test variable for easy debug
	t.Log("testing with following parameters:")
	t.Log(AccessToken)

	// call api
	req, err := http.NewRequest("GET", PrivateBaseUrl+"/access-token/"+AccessToken+"/get-owner", nil)
	req.Header.Set("Authorization", "Bearer XXX")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	body, _ := ioutil.ReadAll(resp.Body)

	user := rest.ResponseDTOUserInfo{}
	apiErrors := apihelper.ApiErrors{}
	json.Unmarshal(body, &user)
	json.Unmarshal(body, &apiErrors)

	t.Log(user)
	t.Log(apiErrors)

	// test response
	if resp.StatusCode != 200 {
		t.Errorf("Expected %s to be %s, got %s", "status", "200", resp.Status)
	} else if user.UserId != uint(userId) {
		t.Errorf("Expected %v to be %v, got %v", "user id", userId, user.UserId)
	}
}
