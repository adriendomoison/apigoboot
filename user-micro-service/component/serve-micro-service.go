package main

import (
	"github.com/adriendomoison/apigoboot/user-micro-service/component/user"
	"github.com/adriendomoison/apigoboot/user-micro-service/component/user/repo"
	"github.com/adriendomoison/apigoboot/user-micro-service/component/user/rest"
	"github.com/adriendomoison/apigoboot/user-micro-service/component/user/service"
	"github.com/adriendomoison/apigoboot/user-micro-service/config"
	"github.com/adriendomoison/apigoboot/user-micro-service/database/dbconn"
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
	router.Use(cors.New(getCORSConfig()))

	// User component
	userComponent := user.New(rest.New(service.New(repo.New())))
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
