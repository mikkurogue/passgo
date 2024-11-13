package ui

import (
	"log"
	"passgo/db"
	"passgo/pkg"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type TableModel struct {
	Table           table.Model
	showModal       bool
	isEmpty         bool
	selectedService *db.Service
}

var (
	modalBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2).
			Width(50)

	modalTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("69")).
			Bold(true)
)

func CreateTableModel() TableModel {
	// clear the screen for consistent ui experience
	tea.ClearScreen()

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

	isEmpty := len(rows) == 0

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

	m := TableModel{t, isEmpty, false, nil}

	return m
}

func (m TableModel) Init() tea.Cmd { return nil }

func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	copier := &pkg.ClipboardCopier{}

	var cmd tea.Cmd
	var database db.Database

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.showModal {
				// make sure we dont just unfocus the table but close the modal IF its open
				m.showModal = false
				m.selectedService = nil
			} else if m.Table.Focused() {
				m.Table.Blur()
			} else {
				m.Table.Focus()
			}
		case "q", "ctrl+c":

			return m, tea.Quit
		case "n":
			return InitialCreateFormModal(), nil
		case "v":
			if selectedRow := m.Table.SelectedRow(); len(selectedRow) > 0 {

				database.CreateInitialConnection()

				id, err := strconv.Atoi(selectedRow[0])
				if err != nil {
					m.showModal = false
					m.selectedService = nil
					log.Fatal("Somehow could not select the row properly for look up")
				}

				service, err := database.FindServiceById(id)
				if err != nil {
					log.Fatal(err)
				}

				m.selectedService = &db.Service{
					Id:       service.Id,
					Username: service.Username,
					Password: service.Password,
					Service:  service.Service,
				}

				m.showModal = true
			}
			return m, nil
		case "m":
			copy(copier, m)
			return m, nil
		case "d":
			if len(m.Table.Rows()) == 0 {
				// do not allow deleting here as it will result in panic
				return m, nil
			}

			database.CreateInitialConnection()
			id, _ := strconv.Atoi(m.Table.SelectedRow()[0])
			database.DeleteService(id)
			tea.ClearScreen()
			list := database.GetAllServices()
			rows := []table.Row{}

			for _, srv := range list {
				rows = append(rows, table.Row{strconv.Itoa(srv.Id), srv.Service, srv.Username})
			}

			m.isEmpty = len(rows) == 0
			m.Table.SetRows(rows)

			// seems counter intuitive atm but this essentially
			// refreshes the model with a clear screen
			// causing the ui to re-render instead of overlap tables like
			// it did before
			return m, tea.ClearScreen
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.Table.SelectedRow()[0]),
			)
		}
	}

	if !m.showModal {
		m.Table, cmd = m.Table.Update(msg)
	}

	return m, cmd
}

func (m TableModel) View() string {
	if m.showModal && m.selectedService != nil {
		return lipgloss.PlaceHorizontal(80, lipgloss.Center, renderModal(m.selectedService))
	}

	if m.isEmpty {
		return lipgloss.PlaceHorizontal(80, lipgloss.Center, "No data available.\nPress 'n' to add a new entry!")
	}

	s := baseStyle.Render(m.Table.View())

	s += "\n"
	s += m.Table.HelpView()
	s += "\n"
	s += "Press q to exit"

	return lipgloss.PlaceHorizontal(80, lipgloss.Center, s)
}

func copy(copier *pkg.ClipboardCopier, m TableModel) {

	var index = 2

	if err := copier.Copy(m.Table.SelectedRow()[index]); err != nil {
		tea.Printf("Uh oh, something went wrong copying to clipboard! Err: %v", err)
	} else {
		tea.Printf("Copied %s to clipboard!", m.Table.SelectedRow()[index])
	}

}

func renderModal(service *db.Service) string {
	if service == nil {
		return ""
	}

	// Content of the modal
	content := modalTextStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			"Service Details",
			"",
			"ID: "+strconv.Itoa(service.Id),
			"Username: "+service.Username,
			"Password: "+service.Password,
			"Service: "+service.Service,
		),
	)

	// Wrap content with the border style
	return modalBorder.Render(content)
}
