package main

import (
	"flag"
	"fmt"
	"os"
	"passgo/db"
	"passgo/ui"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// entry point
func main() {

	reset := flag.Bool("reset", false, "Reset the data store settings and all data stored.")
	flag.Parse()

	columns := []table.Column{
		{Title: "Id", Width: 10},
		{Title: "Service", Width: 15},
		{Title: "Username", Width: 15},
	}

	var database db.Database

	// bootstrap the app on first run.
	if *reset || !db.CheckIfStoreExists() {
		db.Bootstrap(database)
	} else {
		database.CreateInitialConnection()
	}

	serviceList := database.GetAllServices()

	rows := []table.Row{}

	for _, srv := range serviceList {
		rows = append(rows, table.Row{strconv.Itoa(srv.Id), srv.Service, srv.Username})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := ui.Model{t}

	if _, err := tea.NewProgram(m).Run(); err != nil {
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
