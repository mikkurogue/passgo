package main

import (
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"passgo/db"
	"passgo/ui"
)

// entry point
func main() {

	reset := flag.Bool("reset", false, "Reset the data store settings and all data stored.")
	flag.Parse()

	var database db.Database
	// bootstrap the app on first run.
	if *reset || !db.CheckIfStoreExists() {
		db.Bootstrap(database)
	} else {
		database.CreateInitialConnection()
	}

	if _, err := tea.NewProgram(ui.CreateTableModel()).Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	//
	//
	//
	// s := db.Service{
	// 	Username: "test",
	// 	Password: "test",
	// 	Service:  "youtube.com",
	// }
	//
	// database.InsertService(s)
	//
	// database.GetAllServices()
	//
	// fmt.Println("-------- post delete ----------")
	// database.DeleteService(1)
	// database.GetAllServices()
	//
	// database.CloseConnection()
}
