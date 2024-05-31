package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Steps struct {
	viewportWidth int
	steps         []Step
	currentStep   int
}

type Step struct {
	command   string
	counter   int
	duration  time.Duration
	state     State
	startedAt time.Time
}

type tickMsg struct {
}

type startMsg struct {
	id int
}

type exitMsg struct {
	id     int
	state  State
	output string
}

type State string

const (
	Pending State = "🔜"
	Started State = "🟡"
	Exited0 State = "🟢"
	Exited1 State = "❌"
	Skipped State = "🙈"
)

func tick(immediately ...bool) tea.Cmd {
	if len(immediately) > 0 {
		return func() tea.Msg {
			return tickMsg{}
		}
	}
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func newStep(command string) Step {
	return Step{
		command: command,
		state:   Pending,
	}
}

func (m Steps) start(index int) (Steps, tea.Cmd) {
	m.steps[index].state = Started
	m.steps[index].startedAt = time.Now()
	return m, func() tea.Msg {
		command := strings.Split(m.steps[index].command, " ")
		cmd := exec.Command(command[0], command[1:]...) //nolint:gosec
		output, err := cmd.Output()
		m.steps[index].duration = time.Since(m.steps[index].startedAt).Round(time.Millisecond)

		if err != nil {
			m.steps[index].state = Exited1
		} else {
			m.steps[index].state = Exited0
		}

		return exitMsg{
			id:     index,
			state:  m.steps[index].state,
			output: string(output),
		}
	}
}

func (m Steps) Init() tea.Cmd {
	return func() tea.Msg {
		return startMsg{id: 0}
	}
}

func (m Steps) Update(msg tea.Msg) (Steps, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.steps[m.currentStep].state == Started {
			m.steps[m.currentStep].duration = time.Since(m.steps[m.currentStep].startedAt).Round(time.Millisecond)
			return m, tick()
		}
	case startMsg:
		m, cmd := m.start(msg.id)
		return m, tea.Batch(tick(), cmd)
	}
	return m, nil
}

func (m Steps) ViewOne(index int) string {
	var (
		space string
	)
	step := m.steps[index]
	style := lipgloss.NewStyle().
		Margin(0, 1).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		Foreground(lipgloss.Color("#666")).BorderForeground(lipgloss.Color("#aaa"))

	content := fmt.Sprintf("%s %s", step.state, step.command)
	lastly := step.duration.Round(time.Millisecond).String()

	if m.viewportWidth > 0 {
		space = strings.Repeat(" ", m.viewportWidth-lipgloss.Width(content)-lipgloss.Width(lastly)-6)
	} else {
		space = ""
	}
	return style.Render(lipgloss.JoinHorizontal(lipgloss.Center, content, space, lastly))
}

func (m Steps) View() string {
	var content strings.Builder
	for index := range m.steps {
		content.WriteString(m.ViewOne(index) + "\n")
	}
	return content.String()
}
