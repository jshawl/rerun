package main

// A simple program that counts down from 5 and then exits.

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	logfilePath := os.Getenv("DEBUG")
	if logfilePath != "" {
		if _, err := tea.LogToFile(logfilePath, "debug"); err != nil {
			log.Fatal(err)
		}
	}

	p := tea.NewProgram(
		model{
			content: "hello",
			steps:   Steps{},
		},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	content  string
	ready    bool
	viewport viewport.Model
	steps    Steps
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.viewport.YPosition = 0
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	m.viewport.SetContent(m.steps.View())
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return m.viewport.View()
}
