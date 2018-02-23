/*
	Authentication with OAuth2.0 package
*/
package oauth

import (
	"github.com/gin-gonic/gin"
	"github.com/RangelReale/osin"
	"github.com/adriendomoison/gobootapi/database/dbconn"
	"github.com/adriendomoison/gobootapi/oauth/repo"
	"github.com/adriendomoison/gobootapi/oauth/rest"
	"github.com/adriendomoison/gobootapi/oauth/rest/model"
	"github.com/adriendomoison/gobootapi/apicore/config"
)

// MicroService
type MicroService struct {
	rest model.Interface
}

// New return a new micro service instance
func New() (ms *MicroService) {
	return &MicroService{rest.New(initOAuthServer())}
}

// initOAuthServer Init OSIN OAuth server
func initOAuthServer() *osin.Server {
	serverConfig := osin.NewServerConfig()
	serverConfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}
	serverConfig.AllowedAccessTypes = osin.AllowedAccessType{osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN, osin.PASSWORD, osin.CLIENT_CREDENTIALS, osin.ASSERTION}
	serverConfig.AllowGetAccessRequest = true
	serverConfig.AllowClientSecretInParams = true
	return osin.NewServer(serverConfig, repo.New(dbconn.DB.DB()))
}

// Attach link the oauth micro-service with its dependencies to the system
func (ms *MicroService) Attach(router *gin.Engine) {
	if !config.GUnitTestingEnv {
		router.LoadHTMLGlob("oauth/rest/statics/templates/*")
		router.Static("authentication/styles", "oauth/rest/statics/styles")
	}
	oauth2 := router.Group("/authentication")
	oauth2.GET("/authorize", ms.rest.AppAuthorize)
	oauth2.POST("/authorize", ms.rest.AppAuthorize)
	oauth2.POST("/token", ms.rest.AppToken)
	oauth2.POST("/info", ms.rest.AppInfo)
	oauth2.GET("/oauth2/code", ms.rest.AppAuthCode)
	oauth2.GET("/oauth2/token", ms.rest.AppAuthToken)
	oauth2.POST("/oauth2/password", ms.rest.AppAuthPassword)
	oauth2.GET("/oauth2/client_credentials", ms.rest.AppAuthClientCredentials)
	oauth2.GET("/oauth2/assertion", ms.rest.AppAuthAssertion)
	oauth2.GET("/oauth2/refresh", ms.rest.AppAuthRefresh)
	oauth2.GET("/oauth2/info", ms.rest.AppAuthInfo)
}