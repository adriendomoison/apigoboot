package config

import (
	"log"
	"os"
)

// GAppName define the app name
var GAppName = "apigoboot"

var devPort = "4200"
var devAppUrl = "http://api.go.boot"
var prodAppUrl = "https://apigoboot.herokuapp.com"

// GDevEnv define if environment is in dev mode
var GDevEnv bool
// GUnitTestingEnv define if environment is in testing mode
var GUnitTestingEnv bool
// GPort is the application current port
var GPort string
// GAppUrl is the application url
var GAppUrl string

// init initialize the default environment
func init() {
	GPort = os.Getenv("PORT")
	if GPort == "" {
		GDevEnv = true
		GPort = devPort
		GAppUrl = devAppUrl + ":" + GPort
		log.Println("Dev Environnement detected")
	} else {
		GDevEnv = false
		GAppUrl = prodAppUrl
		log.Println("Heroku Environement detected")
	}
}

// SetToTestingEnv set the test environment, this need to be called before testing to prevent the development database to be used
func SetToTestingEnv() {
	GDevEnv = false
	GUnitTestingEnv = true
}
