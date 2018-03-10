package config

import (
	"os"
	"log"
)

var GAppName = "apigoboot"
var GDevPort = "4202"
var GDevAppUrl = "http://apigoboot"
var GProdAppUrl = "https://apigoboot.herokuapp.com"

var GDevEnv bool
var GUnitTestingEnv bool
var GPort string
var GAppUrl string

// init initialize the default environment
func init() {
	GPort = os.Getenv("PORT")
	if GPort == "" {
		GDevEnv = true
		GPort = GDevPort
		GAppUrl = GDevAppUrl + ":" + GPort
		log.Println("Dev Environnement detected")
	} else {
		GDevEnv = false
		GAppUrl = GProdAppUrl
		log.Println("Heroku Environement detected")
	}
}

// SetToTestingEnv set the test environment, this need to be called before testing to prevent the development database to be used
func SetToTestingEnv() {
	GDevEnv = false
	GUnitTestingEnv = true
}