package ui

import (
	"passgo/db"
	"strconv"

	// "strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type TableModel struct {
	Table table.Model
}

func CreateTableModel() TableModel {
	columns := []table.Column{
		{Title: "Id", Width: 10},
		{Title: "Service", Width: 15},
		{Title: "Username", Width: 15},
	}

	var database db.Database
	database.CreateInitialConnection()
	serviceList := database.GetAllServices()
	database.CloseConnection()

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

	m := TableModel{t}

	return m
}

func (m TableModel) Init() tea.Cmd { return nil }

func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.Table.Focused() {
				m.Table.Blur()
			} else {
				m.Table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "n":

			return InitialCreateFormModal(), nil
		case "d":
			var db db.Database

			db.CreateInitialConnection()
			// db.DeleteService()

			id, _ := strconv.Atoi(m.Table.SelectedRow()[0])

			return m, tea.Batch(
				tea.Printf(db.DeleteService(id)),
			)
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.Table.SelectedRow()[0]),
			)
		}
	}
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m TableModel) View() string {

	s := baseStyle.Render(m.Table.View())

	s += "\n"
	s += m.Table.HelpView()
	s += "\n"
	s += "Press q to exit"

	return s
}
