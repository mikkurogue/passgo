package ui

import (
	"fmt"
	"log"
	"passgo/db"
	"passgo/pkg"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	inputs  []textinput.Model
	focused int
	err     error
}

type (
	errMsg error
)

const (
	usr = iota
	pw
	srv
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

func InitialCreateFormModal() model {
	var inputs []textinput.Model = make([]textinput.Model, 3)
	inputs[usr] = textinput.New()
	inputs[usr].Placeholder = "mikkurogue"
	inputs[usr].Focus()
	inputs[usr].CharLimit = 20
	inputs[usr].Width = 30
	inputs[usr].Prompt = ""

	inputs[pw] = textinput.New()
	inputs[pw].Placeholder = "p$55w0rD"
	inputs[pw].CharLimit = 100
	inputs[pw].Width = 20
	inputs[pw].Prompt = ""

	inputs[srv] = textinput.New()
	inputs[srv].Placeholder = "github.com"
	inputs[srv].CharLimit = 50
	inputs[srv].Width = 20
	inputs[srv].Prompt = ""

	return model{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				m.addService()
				return CreateTableModel(), nil
			}
			m.nextInput()
		case tea.KeyEsc:
			return CreateTableModel(), nil
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()

	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return fmt.Sprintf(
		`
 %s
 %s

 %s
 %s
 %s
 %s

 %s
`,
		inputStyle.Width(30).Render("Username"),
		m.inputs[usr].View(),
		inputStyle.Width(30).Render("Password"),
		m.inputs[pw].View(),
		inputStyle.Width(30).Render("Service"),
		m.inputs[srv].View(),
		continueStyle.Render("Create ->"),
	) + "\n"
}

func (m *model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *model) prevInput() {
	m.focused--
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

func (m *model) addService() {
	var database db.Database
	database.CreateInitialConnection()

	encrypted, err := pkg.Encrypt(m.inputs[pw].Value(), pkg.KEY)
	if err != nil {
		log.Fatal(err)
	}

	database.InsertService(db.Service{
		Username: m.inputs[usr].Value(),
		Password: encrypted,
		Service:  m.inputs[srv].Value(),
	})
	database.CloseConnection()
}
