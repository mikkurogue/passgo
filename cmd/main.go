package main

import (
	"flag"
	"passgo/db"
)

// entry point
func main() {

	reset := flag.Bool("reset", false, "Reset the data store settings and all data stored.")
	flag.Parse()

	var database db.Database

	// bootstrap the app on first run.
	if *reset || !db.CheckIfStoreExists() {
		db.Bootstrap(database)
	}

	s := db.Service{
		Username: "test",
		Password: "test",
		Service:  "youtube.com",
	}

	database.InsertService(s)

	database.GetAllServices()

	database.CloseConnection()
}
