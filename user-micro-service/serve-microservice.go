package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/adriendomoison/gobootapi/user-micro-service/config"
	"github.com/adriendomoison/gobootapi/user-micro-service/usercomponent"
	"github.com/adriendomoison/gobootapi/user-micro-service/database/dbconn"
	"github.com/adriendomoison/gobootapi/user-micro-service/usercomponent/repo"
	"github.com/adriendomoison/gobootapi/user-micro-service/usercomponent/service"
	"github.com/adriendomoison/gobootapi/user-micro-service/usercomponent/rest"
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

	// User component
	userComponent := usercomponent.New(rest.New(service.New(repo.New())))
	userComponent.AttachPublicAPI(router.Group("/api/v1"))
	userComponent.AttachPrivateAPI(router.Group("/api/private-v1"))

	// Start router
	go log.Println("Service user started: Navigate to " + config.GAppUrl)
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