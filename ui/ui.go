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

type Model struct {
	Table table.Model
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m Model) View() string {

	s := baseStyle.Render(m.Table.View())

	s += "\n"
	s += m.Table.HelpView()
	s += "\n"
	s += "Press q to exit"

	return s
}
