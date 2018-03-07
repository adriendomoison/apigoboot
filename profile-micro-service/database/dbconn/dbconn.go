package dbconn

import (
	"os"
	"log"
	"syscall"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/adriendomoison/gobootapi/profile-micro-service/config"
)

var DB *gorm.DB

// Connect connect to database depending of the env
func Connect() {
	if config.GUnitTestingEnv {
		connectToDB(config.GAppName+"_test", config.GAppName+"_test", config.GAppName+"_test", "localhost")
	} else if _, ok := syscall.Getenv("DYNO"); ok {
		connectToDB(os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"))
	} else if config.GDevEnv {
		connectToDB(config.GAppName, config.GAppName, config.GAppName, "localhost")
	}
}

// connectToDB do the connection request to the database depending on provided parameters
func connectToDB(username string, dbName string, password string, host string) (err error) {
	log.Println("CONNECTING TO [" + dbName + "] DB...")
	DB, err = gorm.Open("postgres", "host="+host+" user="+username+" dbname="+dbName+" sslmode=disable password="+password)
	if err != nil {
		log.Panic("Database status: [Failed to connect]", err)
	} else {
		log.Println("Database status: [Connected]")
	}
	DB.SingularTable(true)
	return
}