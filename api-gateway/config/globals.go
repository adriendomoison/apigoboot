/*
	config generate the environment of the API
*/
package config

import (
	"os"
	"log"
)

var GDevPort = "4000"
var GDevAppUrl = "http://api.go.boot"
var GProdAppUrl = "https://apigoboot.herokuapp.com"

var GPort string
var GAppUrl string

// init initialize the default environment
func init() {
	GPort = os.Getenv("PORT")
	if GPort == "" {
		GPort = GDevPort
		GAppUrl = GDevAppUrl + ":" + GPort
		log.Println("Dev Environnement detected")
	} else {
		GAppUrl = GProdAppUrl
		log.Println("Heroku Environement detected")
	}
}