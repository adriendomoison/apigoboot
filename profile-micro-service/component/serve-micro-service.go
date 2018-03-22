// Package main
package main

import (
	"github.com/adriendomoison/apigoboot/api-tool/apitool"
	"github.com/adriendomoison/apigoboot/profile-micro-service/component/profile"
	"github.com/adriendomoison/apigoboot/profile-micro-service/component/profile/repo"
	"github.com/adriendomoison/apigoboot/profile-micro-service/component/profile/rest"
	"github.com/adriendomoison/apigoboot/profile-micro-service/component/profile/service"
	"github.com/adriendomoison/apigoboot/profile-micro-service/config"
	"github.com/adriendomoison/apigoboot/profile-micro-service/database/dbconn"
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

	// Profile component
	profileComponent := profile.New(rest.New(service.New(repo.New())))
	profileComponent.AttachPublicAPI(router.Group("/api/v1"))
	profileComponent.AttachPrivateAPI(router.Group("/api/private-v1"))

	// Start router
	go log.Println("Service profile started: Navigate to " + config.GAppUrl)
	router.Run(":" + config.GPort)
}
