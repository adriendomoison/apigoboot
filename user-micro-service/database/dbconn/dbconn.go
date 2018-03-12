package dbconn

import (
	"github.com/adriendomoison/apigoboot/user-micro-service/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
	"syscall"
	"time"
)

var DB *gorm.DB

// Connect connect to database depending of the env
func Connect() (err error) {
	if config.GUnitTestingEnv {
		err = connectToDB(config.GAppName+"_test", config.GAppName+"_test", config.GAppName+"_test", "localhost")
	} else if _, ok := syscall.Getenv("DYNO"); ok {
		err = connectToDB(os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"))
	} else if config.GDevEnv {
		err = connectToDB(config.GAppName, config.GAppName, config.GAppName, "db")
	}
	if err != nil {
		log.Panic("Database status: [Failed to connect]", err)
	}
	return
}

// connectToDB do the connection request to the database depending on provided parameters
func connectToDB(username string, dbName string, password string, host string) (err error) {
	log.Println("CONNECTING TO [" + dbName + "] DB...")
	for i := 0; i < 5; i++ {
		DB, err = gorm.Open("postgres", "host="+host+" user="+username+" dbname="+dbName+" sslmode=disable password="+password)
		if err != nil {
			log.Println("Still trying...")
		} else {
			DB.SingularTable(true)
			log.Println("Database status: [Connected]")
			break
		}
		time.Sleep(5 * time.Second)
	}
	return
}
