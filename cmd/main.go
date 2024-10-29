package main

import (
	"fmt"
	"passgo/db"
)

// entry point
func main() {
	// err := godotenv.Load()
	//
	// if err != nil {
	// 	fmt.Println("No env file found... exiting...")
	// 	os.Exit(0)
	// }

	var d db.Database

	err := d.CreateInitialConnection()
	if err != nil {
		fmt.Println(err)
	}

	createErr := d.CreateStoreTable()
	if createErr != nil {
		fmt.Println(createErr.Error())
	}

	s := db.Service{
		Username: "test",
		Password: "test",
		Service:  "youtube.com",
	}

	d.InsertService(s)

	d.GetAllServices()

	d.CloseConnection()
}
