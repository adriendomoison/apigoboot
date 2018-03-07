package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/adriendomoison/gobootapi/profile-micro-service/config"
	"github.com/adriendomoison/gobootapi/profile-micro-service/profilecomponent"
	"github.com/adriendomoison/gobootapi/profile-micro-service/database/dbconn"
	"github.com/adriendomoison/gobootapi/profile-micro-service/profilecomponent/repo"
	"github.com/adriendomoison/gobootapi/profile-micro-service/profilecomponent/service"
	"github.com/adriendomoison/gobootapi/profile-micro-service/profilecomponent/rest"
)

// startAPI start the API and keep it alive
func main() {
	// Init DB and plan to close it at the end of the programme
	dbconn.Connect()
	defer dbconn.DB.Close()

	// Init router
	router := gin.Default()
	router.Use(cors.New(getCORSConfig()))

	// User component
	profileComponent := profilecomponent.New(rest.New(service.New(repo.New())))
	profileComponent.AttachPublicAPI(router.Group("/api/v1"))
	profileComponent.AttachPrivateAPI(router.Group("/api/private-v1"))

	// Start router
	go log.Println("Service profile started: Navigate to " + config.GAppUrl)
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