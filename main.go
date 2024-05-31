package main

import (
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if os.Getenv("DEBUG") != "" {
		_, _ = tea.LogToFile("debug.log", "debug")
	}

	_, _ = tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	).Run()
}

type model struct {
	content  string
	ready    bool
	viewport viewport.Model
	steps    Steps
}

func initialModel() model {
	var steps []Step //nolint:prealloc
	for _, command := range parseConfig()["steps"] {
		steps = append(steps, newStep(command))
	}

	return model{
		content:  "",
		ready:    false,
		viewport: viewport.New(0, 0),
		steps: Steps{
			currentStep:   0,
			steps:         steps,
			viewportWidth: 0,
		},
	}
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
	cmds = append(cmds, cmd)

	m.viewport.SetContent(m.steps.View())
	m.steps, cmd = m.steps.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return m.viewport.View()
}
