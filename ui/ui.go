package ui

import (
	"fmt"
	"log"
	"passgo/db"
	"passgo/pkg"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type TableModel struct {
	Table             table.Model
	searchInput       textinput.Model
	showModal         bool
	showSearch        bool
	isEmpty           bool
	selectedService   *db.Service
	Notification      string
	NotificationTimer *time.Timer
}

type NotificationTimeoutMsg struct{}

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

	ti := textinput.New()
	ti.Placeholder = "[Search for a service...]"
	ti.Focus()
	ti.Prompt = ""
	ti.CharLimit = 256
	ti.Width = 40

	ti.TextStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true).
		BorderForeground(lipgloss.Color("5")).
		Align(lipgloss.Center)

	m := TableModel{t, ti, isEmpty, false, false, nil, "", nil}

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
			} else if m.showSearch {
				m.showSearch = false
			} else if m.Table.Focused() {
				m.Table.Blur()
			} else if m.Notification != "" {
				m.Notification = ""
			} else {
				m.Table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "n":
			return InitialCreateFormModal(), nil
		case "/":
			m.showSearch = !m.showSearch
			return m, nil
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
		case "c":
			id, err := strconv.Atoi(m.Table.SelectedRow()[0])
			if err != nil {
				m.Notification = "Invalid row selection"
				return m, nil
			}
			copy(&m, copier, id)
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
	case NotificationTimeoutMsg:
		m.Notification = ""
		return m, nil
	}

	if m.showSearch {
		var inputCmd tea.Cmd
		m.searchInput, inputCmd = m.searchInput.Update(msg)
		return m, inputCmd
	}

	if !m.showModal {
		m.Table, cmd = m.Table.Update(msg)
	}

	return m, cmd
}

func (m TableModel) View() string {
	var view string

	// Always start with the table (or modal/search as necessary)
	if m.showModal && m.selectedService != nil {
		view = lipgloss.PlaceHorizontal(80, lipgloss.Center, renderModal(m.selectedService))
	} else if m.isEmpty {
		view = lipgloss.PlaceHorizontal(80, lipgloss.Center, "No data available.\nPress 'n' to add a new entry!")
	} else if m.showSearch {
		view = lipgloss.PlaceHorizontal(80, lipgloss.Center, renderSearchInput(m))
	} else {
		view = baseStyle.Render(m.Table.View()) + "\n"
		view += m.Table.HelpView() + "\n"
		view += "Press q to exit"
		view = lipgloss.PlaceHorizontal(80, lipgloss.Center, view)
	}

	// Add the notification at the bottom if it exists
	if m.Notification != "" {
		notificationStyle := lipgloss.NewStyle().
			Background(lipgloss.Color("2")).
			Foreground(lipgloss.Color("0")).
			Padding(1).
			Align(lipgloss.Center)

		notification := notificationStyle.Render(m.Notification)
		view += "\n\n" + notification
	}

	return view
}

func copy(m *TableModel, copier *pkg.ClipboardCopier, rowId int) tea.Cmd {
	var database db.Database

	database.CreateInitialConnection()
	srv, err := database.FindServiceById(rowId)
	if err != nil {
		m.Notification = "No service found for the selected row"
		return nil
	}

	decrypted, err := pkg.Decrypt(srv.Password, pkg.Key)
	if err != nil {
		log.Fatal(err)
	}

	if err := copier.Copy(decrypted); err != nil {
		m.Notification = fmt.Sprintf("Could not copy password. Error: %v", err)
	} else {
		m.Notification = fmt.Sprintf("Copied %s to clipboard", srv.Service)
	}

	// Stop previous timer if active
	if m.NotificationTimer != nil {
		m.NotificationTimer.Stop()
	}

	// Start a new timer
	m.NotificationTimer = time.NewTimer(3 * time.Second)

	return func() tea.Msg {
		<-m.NotificationTimer.C
		return NotificationTimeoutMsg{}
	}
}

func renderSearchInput(m TableModel) string {
	var b strings.Builder

	b.WriteString(lipgloss.NewStyle().
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder(), true).
		BorderForeground(lipgloss.Color("5")).
		Render(m.searchInput.View()))

	return b.String()
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
