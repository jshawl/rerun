package main

// A simple program that counts down from 5 and then exits.

import (
	"log"
	"os"

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
		model{},
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	default:
		return m, nil
	}
}

func (m model) View() string {
	return "hey"
}
