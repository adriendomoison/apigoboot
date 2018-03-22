// Package main
package main

import (
	"github.com/RangelReale/osin"
	"github.com/adriendomoison/apigoboot/api-tool/apitool"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/component/oauth2"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/component/oauth2/repo"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/component/oauth2/rest"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/component/oauth2/service"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/config"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/database/dbconn"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

// startAPI start the API and keep it alive
func main() {
	// Init DB and plan to close it at the end of the programme
	dbconn.Connect()
	defer dbconn.DB.Close()

	// Set GIN in production mode if run in production
	if !config.GDevEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	// Init router
	router := gin.Default()
	router.Use(cors.New(apitool.DefaultCORSConfig()))

	// Init statics
	if !config.GUnitTestingEnv {
		router.LoadHTMLGlob("component/oauth2/rest/statics/templates/*")
		router.Static("authentication/styles", "component/oauth2/rest/statics/styles")
	}

	// Oauth2 components
	oauth2Component := oauth2.New(rest.New(initOAuthServer(), service.New(repo.New())))
	oauth2Component.AttachPublicAPI(router.Group("/authentication"))
	oauth2Component.AttachPrivateAPI(router.Group("/api/private-v1/authentication"))

	// Start router
	go log.Println("Service oauth2 started: Navigate to " + config.GAppUrl)
	router.Run(":" + config.GPort)
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
