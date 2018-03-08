package core

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/adriendomoison/apigoboot/apicore/config"
	"github.com/adriendomoison/apigoboot/apicore/rest"
)

// startAPI start the API and keep it alive
func StartAPI() {

	// Init router
	router := gin.Default()
	router.Use(cors.New(getCORSConfig()))

	// Add a root path to get a quick overview of the server status
	router.GET("/", rest.AppInfo)

	// Start router
	go log.Println("Platform started: Navigate to " + config.GAppUrl)
	router.Run(":" + config.GPort)
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
