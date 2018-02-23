package core

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/adriendomoison/go-boot-api/oauth"
	"github.com/adriendomoison/go-boot-api/user"
	"github.com/adriendomoison/go-boot-api/profile"
	"github.com/adriendomoison/go-boot-api/apicore/config"
	"github.com/adriendomoison/go-boot-api/apicore/rest"
)

// startAPI start the API and keep it alive
func StartAPI() {

	// Init router
	router := gin.Default()
	router.Use(cors.New(getCORSConfig()))

	// Plug software micro-services
	plugMicroServices(router)

	// Start router
	go log.Println("Platform started: Navigate to " + config.GAppUrl)
	router.Run(":" + config.GPort)
}

// plugMicroServices attach all micro-services of the API
func plugMicroServices(router *gin.Engine) {
	// Add a root path to get a quick overview of the server status
	router.GET("/", rest.AppInfo)

	// Authentication micro-service
	oauth.New().Attach(router)

	// Add all the following route under the same root "/api/v1"
	apiV1 := router.Group("/api/v1")

	// User micro-service
	user.New().Attach(apiV1)
	// Profile micro-service
	profile.New().Attach(apiV1)
	// Your micro-service
	// ...
}

// getCORSConfig Generate CORS config for router
func getCORSConfig() cors.Config {
	CORSConfig := cors.DefaultConfig()
	CORSConfig.AllowCredentials = true
	CORSConfig.AllowAllOrigins = true
	CORSConfig.AllowHeaders = []string{"*", "Origin", "Content-Type", "Authorization", "Cookie"}
	CORSConfig.AllowMethods = []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"}
	return CORSConfig
}
