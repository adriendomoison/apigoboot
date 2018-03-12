package core

import (
	"github.com/adriendomoison/apigoboot/api-gateway/config"
	"github.com/adriendomoison/apigoboot/api-gateway/rest"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

// StartAPIGateway start the API and keep it alive
func StartAPIGateway() {

	// Init router
	router := gin.Default()
	router.Use(cors.New(getCORSConfig()))

	// Add a root path to get a quick overview of the server status
	attachRoutes(router)

	// Start router
	go log.Println("Platform started: Navigate to " + config.GAppUrl)
	router.Run(":" + config.GPort)
}

func attachRoutes(router *gin.Engine) {
	router.GET("/", rest.AppInfo)
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
