package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/RangelReale/osin"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/config"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/oauth2component"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/database/dbconn"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/oauth2component/rest"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/oauth2component/repo"
	"github.com/adriendomoison/apigoboot/oauth2-micro-service/oauth2component/service"
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
	router.Use(cors.New(getCORSConfig()))

	// Init statics
	if !config.GUnitTestingEnv {
		router.LoadHTMLGlob("oauth2-micro-service/rest/statics/templates/*")
		router.Static("authentication/styles", "oauth2-micro-service/rest/statics/styles")
	}

	// Oauth2 components
	oauth2Component := oauth2component.New(rest.New(initOAuthServer(), service.New(repo.New())))
	oauth2Component.AttachPublicAPI(router.Group("/authentication"))
	oauth2Component.AttachPrivateAPI(router.Group("/api/private-v1/authentication"))

	// Start router
	go log.Println("Service oauth2 started: Navigate to " + config.GAppUrl)
	router.Run(":" + config.GPort)
}

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