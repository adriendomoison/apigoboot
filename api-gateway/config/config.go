// Package config generate the environment of the API
package config

import (
	"log"
	"os"
)

var devPort = "4200"
var devAppUrl = "http://api.go.boot"
var prodAppUrl = "https://apigoboot.herokuapp.com"

// GPort is the application current port
var GPort string

// GAppUrl is the application url
var GAppUrl string

// init initialize the default environment
func init() {
	GPort = os.Getenv("PORT")
	if GPort == "" {
		GPort = devPort
		GAppUrl = devAppUrl + ":" + GPort
		log.Println("Dev Environnement detected")
	} else {
		GAppUrl = prodAppUrl
		log.Println("Heroku Environement detected")
	}
}
