package main_test

import (
	"encoding/json"
	"github.com/RangelReale/osin"
	"github.com/adriendomoison/apigoboot/api-tool/apitool"
	"github.com/adriendomoison/apigoboot/api-tool/errorhandling/apihelper"
	"github.com/adriendomoison/apigoboot/api-tool/errorhandling/servicehelper"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/component/oauth2"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/component/oauth2/repo"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/component/oauth2/rest"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/component/oauth2/service"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/config"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/database/dbconn"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

var publicBaseUrl = config.GAppUrl + "/authentication"
var privateBaseUrl = config.GAppUrl + "/api/private-v1/authentication"
var accessToken = ""

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
	form.Add("username", "test00@example.dev")
	form.Add("password", "password123")

	req, err := http.NewRequest("POST", publicBaseUrl+"/authorize?response_type=code&client_id=apigoboot&client_secret=apigoboot&state=xyz&scope=everything&redirect_uri=http://api.go.boot:4200/authentication/oauth2/code", strings.NewReader(form.Encode()))
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
		if reqDTO.Username == "test00@example.dev" && reqDTO.Password == "password123" {
			c.JSON(http.StatusOK, rest.ResponseDTOUserInfo{
				UserId: 1,
				Email:  "test00@example.dev",
			})
		} else {
			c.JSON(apihelper.BuildResponseError(&servicehelper.Error{
				Detail: errors.New("bad username or password"),
				Code:   servicehelper.BadRequest,
			}))
		}
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
	router.Use(cors.New(apitool.DefaultCORSConfig()))

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

	// Wait and check if the http server is running
	apitool.WaitForServerToStart(publicBaseUrl + "/")

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
		RedirectUri: "http://api.go.boot:4200/authentication/oauth2/code",
		UserId:      1,
	})
}

func TestPasswordAuthentication(t *testing.T) {

	// init test variable
	requestBody := map[string]string{
		"client_id":     "apigoboot",
		"client_secret": "apigoboot",
		"method":        "password",
		"username":      "test00@example.dev",
		"password":      "password123",
	}

	// call api
	access := struct {
		AccessToken string `json:"access_token"`
	}{}
	resp, _ := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method:      "POST",
		URL:         publicBaseUrl + "/oauth2/password",
		ContentType: "application/x-www-form-urlencoded",
	}, requestBody, &access)
	defer resp.Body.Close()

	// test response
	if resp.Status != "200 OK" {
		t.Errorf("Status %s was not expected, user was not updated", resp.Status)
	} else if access.AccessToken == "" {
		t.Error("Access token is empty")
	}

	accessToken = access.AccessToken
}

func TestCodeAuthentication(t *testing.T) {

	// init test variable
	code := requestCode(t)

	// call api
	access := struct {
		AccessToken string `json:"access_token"`
	}{}
	resp, _ := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method:      "GET",
		URL:         publicBaseUrl + "/oauth2/code?code=" + code + "&state=xyz&client_id=apigoboot&client_secret=apigoboot&parse=yes",
		ContentType: "application/x-www-form-urlencoded",
	}, nil, &access)
	defer resp.Body.Close()

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
	t.Log(accessToken)

	// call api
	var user rest.ResponseDTOUserInfo
	resp, _ := apitool.HttpRequestHandlerForUnitTesting(t, apitool.RequestHeader{
		Method:        "GET",
		URL:           privateBaseUrl + "/access-token/" + accessToken + "/get-owner",
		Authorization: "Bearer " + accessToken,
	}, nil, &user)
	defer resp.Body.Close()

	// test response
	if resp.StatusCode != 200 {
		t.Errorf("Expected %s to be %s, got %s", "status", "200", resp.Status)
	} else if user.UserId != uint(userId) {
		t.Errorf("Expected %v to be %v, got %v", "user id", userId, user.UserId)
	}
}
