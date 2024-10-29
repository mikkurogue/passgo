package main

import "passgo/db"

// entry point
func main() {
	// err := godotenv.Load()
	//
	// if err != nil {
	// 	fmt.Println("No env file found... exiting...")
	// 	os.Exit(0)
	// }

	var d db.Database

	d.CreateInitialConnection()
	d.CreateStoreTable()

	s := db.Service{
		Username: "test",
		Password: "test",
		Service:  "youtube.com",
	}

	d.InsertService(s)
}
