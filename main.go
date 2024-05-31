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
		if _, err := tea.LogToFile("debug.log", "debug"); err != nil {
			log.Fatal(err)
		}
	}

	p := tea.NewProgram(
		model{
			content: "hello",
			steps: Steps{
				steps: []Step{
					newStep("sleep 2"),
					newStep("go test -v ./..."),
					newStep("go build"),
				},
			},
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
	return m.steps.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		if k == "q" {
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
		m.steps.viewportWidth = msg.Width
	}

	m.viewport, cmd = m.viewport.Update(msg)
	m.viewport.SetContent(m.steps.View())
	cmds = append(cmds, cmd)

	m.steps, cmd = m.steps.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return m.viewport.View()
}
