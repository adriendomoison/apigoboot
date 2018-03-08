package main

import (
	"github.com/adriendomoison/apigoboot/apicore/core"
	"github.com/adriendomoison/apigoboot/apicore/database/dbconn"
)

func main() {

	// Init DB and plan to close it at the end of the programme
	dbconn.Connect()
	defer dbconn.DB.Close()

	// Start API
	core.StartAPI()

}
