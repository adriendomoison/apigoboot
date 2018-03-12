package rest

import (
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/adriendomoison/apigoboot/errorhandling/apihelper"
	"github.com/adriendomoison/apigoboot/errorhandling/servicehelper"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/config"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"log"
	"net/http"
	"net/url"
)

type ServiceInterface interface {
	AskUserServiceToCheckCredentials(username string, password string, method string) (ResponseDTOUserInfo, *apihelper.ApiErrors)
	GetResourceOwnerId(token string) (ResponseDTOUserInfo, *servicehelper.Error)
}

type rest struct {
	server  *osin.Server
	service ServiceInterface
}

type clientCredential struct {
	Method       string `json:"method"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type RequestDTOUserCredentials struct {
	Method   string `json:"method"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseDTOUserInfo struct {
	UserId uint   `json:"user_id"`
	Email  string `json:"email"`
}

// New return a new rest instance
func New(server *osin.Server, service ServiceInterface) *rest {
	return &rest{server, service}
}

// Authorization code endpoint
func (r *rest) AppAuthorize(c *gin.Context) {
	resp := r.server.NewResponse()
	defer resp.Close()
	if ar := r.server.HandleAuthorizeRequest(resp, c.Request); ar != nil {
		if userId, ok := HandleLoginPage(r, ar, c); !ok {
			return
		} else {
			ar.UserData = userId
			ar.Authorized = true
			r.server.FinishAuthorizeRequest(resp, c.Request, ar)
		}
	}
	if resp.IsError && resp.InternalError != nil {
		log.Printf("ERROR: %s\n", resp.InternalError)
	}
	if !resp.IsError {
		//resp.Output["custom_parameter"] = 42
	}
	osin.OutputJSON(resp, c.Writer, c.Request)
}

// Access token endpoint
func (r *rest) AppToken(c *gin.Context) {
	resp := r.server.NewResponse()
	defer resp.Close()
	if ar := r.server.HandleAccessRequest(resp, c.Request); ar != nil {
		ar.UserData = uint(0)
		switch ar.Type {
		case osin.AUTHORIZATION_CODE:
			ar.Authorized = true
		case osin.REFRESH_TOKEN:
			ar.Authorized = true
		case osin.PASSWORD:
			if userInfo, err := r.service.AskUserServiceToCheckCredentials(ar.Username, ar.Password, c.Query("method")); err == nil {
				ar.Authorized = true
				ar.UserData = userInfo.UserId
			}
		case osin.CLIENT_CREDENTIALS:
			ar.Authorized = true
		case osin.ASSERTION:
			if ar.AssertionType == "urn:osin.example.complete" && ar.Assertion == "osin.data" {
				ar.Authorized = true
			}
		}
		r.server.FinishAccessRequest(resp, c.Request, ar)
	}
	if resp.IsError && resp.InternalError != nil {
		log.Printf("ERROR: %s\n", resp.InternalError)
	}
	if !resp.IsError {
		//resp.Output["custom_parameter"] = 42
	}
	osin.OutputJSON(resp, c.Writer, c.Request)
}

// Information endpoint
func (r *rest) AppInfo(c *gin.Context) {
	resp := r.server.NewResponse()
	defer resp.Close()

	if ir := r.server.HandleInfoRequest(resp, c.Request); ir != nil {
		r.server.FinishInfoRequest(resp, c.Request, ir)
	}
	osin.OutputJSON(resp, c.Writer, c.Request)
}

// Application destination - CODE
func (r *rest) AppAuthCode(c *gin.Context) {

	// build credentials
	var cc clientCredential
	cc.ClientId = c.Query("client_id")
	cc.ClientSecret = c.Query("client_secret")

	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no code was provided"})
		return
	}

	jr := make(map[string]interface{})

	// build access code url
	authURL := fmt.Sprintf(
		"%s/authentication/token?grant_type=authorization_code&client_id=%s&client_secret=%s&state=xyz&redirect_uri=%s&code=%s",
		config.GAppUrl,
		cc.ClientId,
		cc.ClientSecret,
		url.QueryEscape(
			fmt.Sprintf("%s/authentication/oauth2/code", config.GAppUrl),
		), url.QueryEscape(code),
	)

	// build app credentials
	auth := osin.BasicAuth{Username: cc.ClientId, Password: cc.ClientSecret}

	// if parse, download and parse json
	if c.Query("parse") == "yes" {
		if err := DownloadAccessToken(authURL, &auth, jr); err != nil {
			c.JSON(apihelper.BuildRequestError(err))
			return
		}
	}

	// show json error
	if err, ok := jr["error"]; ok {
		c.JSON(apihelper.BuildRequestError(errors.New(err)))
		return
	}

	if _, ok := jr["access_token"]; ok {
		c.JSON(http.StatusOK, jr)
		return
	}

	cururl := *c.Request.URL
	curq := cururl.Query()
	curq.Add("parse", "yes")
	cururl.RawQuery = curq.Encode()
	c.JSON(http.StatusOK, gin.H{
		"state": c.Query("state"),
		"code":  c.Query("code"),
	})
}

// Application destination - TOKEN
func (r *rest) AppAuthToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Response data in fragment - not acessible via server - Nothing to do"})
}

// Application destination - PASSWORD
func (r *rest) AppAuthPassword(c *gin.Context) {

	// build credentials
	var cc clientCredential
	if err := c.BindJSON(&cc); err != nil {
		c.JSON(apihelper.BuildRequestError(err))
		return
	}

	// Encode credentials to escape special characters
	cc.Username = url.QueryEscape(cc.Username)
	cc.Password = url.QueryEscape(cc.Password)

	jr := make(map[string]interface{})

	// build authentication URL
	authURL := fmt.Sprintf(
		"%s/authentication/token?grant_type=password&scope=everything&username=%s&password=%s&method=%s",
		config.GAppUrl, cc.Username, cc.Password, cc.Method,
	)

	// build App credentials
	auth := osin.BasicAuth{Username: cc.ClientId, Password: cc.ClientSecret}

	// download token
	if err := DownloadAccessToken(authURL, &auth, jr); err != nil {
		c.JSON(apihelper.BuildRequestError(err))
		return
	}

	// show json error
	if err, ok := jr["error"]; ok {
		c.JSON(apihelper.BuildRequestError(errors.New(err)))
		return
	}

	// show json access token
	c.JSON(http.StatusOK, jr)
}

// Application destination - CLIENT_CREDENTIALS
func (r *rest) AppAuthClientCredentials(c *gin.Context) {

	// build credentials
	var cc clientCredential
	if err := c.BindJSON(&cc); err != nil {
		c.JSON(apihelper.BuildRequestError(err))
		return
	}

	// Encode credentials to escape special characters
	cc.Username = url.QueryEscape(cc.Username)
	cc.Password = url.QueryEscape(cc.Password)

	jr := make(map[string]interface{})

	// build access code url
	authURL := fmt.Sprintf("%s/authentication/token?grant_type=client_credentials", config.GAppUrl)

	// build app credentials
	auth := osin.BasicAuth{Username: cc.ClientId, Password: cc.ClientSecret}

	// download token
	err := DownloadAccessToken(authURL, &auth, jr)
	if err != nil {
		c.JSON(apihelper.BuildRequestError(err))
		return
	}

	// show json error
	if err, ok := jr["error"]; ok {
		c.JSON(apihelper.BuildRequestError(errors.New(err)))
		return
	}

	// show json access token
	c.JSON(http.StatusOK, jr)
}

// Application destination - ASSERTION
func (r *rest) AppAuthAssertion(c *gin.Context) {

	// build credentials
	var cc clientCredential
	if err := c.BindJSON(&cc); err != nil {
		c.JSON(apihelper.BuildRequestError(err))
		return
	}

	// Encode credentials to escape special characters
	cc.Username = url.QueryEscape(cc.Username)
	cc.Password = url.QueryEscape(cc.Password)

	jr := make(map[string]interface{})

	// build access code url
	authURL := fmt.Sprintf(
		"%s/authentication/token?grant_type=assertion&assertion_type=urn:osin.example.complete&assertion=osin.data",
		config.GAppUrl,
	)

	// build app credentials
	auth := osin.BasicAuth{Username: cc.ClientId, Password: cc.ClientSecret}

	// download token
	err := DownloadAccessToken(authURL, &auth, jr)
	if err != nil {
		c.JSON(apihelper.BuildRequestError(err))
		return
	}

	// show json error
	if err, ok := jr["error"]; ok {
		c.JSON(apihelper.BuildRequestError(errors.New(err)))
		return
	}

	// show json access token
	c.JSON(http.StatusOK, jr)
}

// Application destination - REFRESH
func (r *rest) AppAuthRefresh(c *gin.Context) {

	// build credentials
	var cc clientCredential
	if err := c.BindJSON(&cc); err != nil {
		c.JSON(apihelper.BuildRequestError(err))
		return
	}

	// Encode credentials to escape special characters
	cc.Username = url.QueryEscape(cc.Username)
	cc.Password = url.QueryEscape(cc.Password)

	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no code was provided"})
		return
	}

	jr := make(map[string]interface{})

	// build access code url
	authURL := fmt.Sprintf(
		"%s/authentication/token?grant_type=refresh_token&refresh_token=%s",
		config.GAppUrl, url.QueryEscape(code),
	)

	// build app credentials
	auth := osin.BasicAuth{Username: cc.ClientId, Password: cc.ClientSecret}

	// download token
	err := DownloadAccessToken(authURL, &auth, jr)
	if err != nil {
		c.JSON(apihelper.BuildRequestError(err))
		return
	}

	// show json error
	if err, ok := jr["error"]; ok {
		c.JSON(apihelper.BuildRequestError(errors.New(err)))
		return
	}

	// show json access token
	c.JSON(http.StatusOK, jr)

}

// Application destination - INFO
func (r *rest) AppAuthInfo(c *gin.Context) {

	// build credentials
	var cc clientCredential
	if err := c.BindJSON(&cc); err != nil {
		c.JSON(apihelper.BuildRequestError(err))
		return
	}

	// Encode credentials to escape special characters
	cc.Username = url.QueryEscape(cc.Username)
	cc.Password = url.QueryEscape(cc.Password)

	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no code was provided"})
		return
	}

	jr := make(map[string]interface{})

	// build access code url
	authURL := fmt.Sprintf("%s/authentication/info?code=%s", config.GAppUrl, url.QueryEscape(code))

	// build app credentials
	auth := osin.BasicAuth{Username: cc.ClientId, Password: cc.ClientSecret}

	// download token
	err := DownloadAccessToken(authURL, &auth, jr)
	if err != nil {
		c.JSON(apihelper.BuildRequestError(err))
		return
	}

	// show json error
	if err, ok := jr["error"]; ok {
		c.JSON(apihelper.BuildRequestError(errors.New(err)))
		return
	}

	// show json access token
	c.JSON(http.StatusOK, jr)
}

func (r *rest) GetAccessTokenOwnerUserId(c *gin.Context) {
	if resDTO, err := r.service.GetResourceOwnerId(c.Param("accessToken")); err != nil {
		c.JSON(apihelper.BuildResponseError(err))
	} else {
		c.JSON(http.StatusOK, resDTO)
	}
}
