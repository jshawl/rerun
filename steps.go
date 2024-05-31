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
	duration  time.Duration
	state     State
	startedAt time.Time
	output    string
}

type tickMsg struct {
}

type startMsg struct {
	id int
}

type exitMsg struct {
	id     int
	output string
	err    error
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

		return exitMsg{
			id:     index,
			output: string(output),
			err:    err,
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
	case exitMsg:
		if msg.err != nil {
			m.steps[m.currentStep].state = Exited1
			m.steps[m.currentStep].output = msg.output

			for i := m.currentStep + 1; i < len(m.steps); i++ {
				m.steps[i].state = Skipped
			}

			return m, nil
		} else {
			m.steps[m.currentStep].state = Exited0
		}
		if m.currentStep < len(m.steps)-1 {
			m.steps[m.currentStep].output = msg.output
			m.currentStep++
			m, cmd := m.start(m.currentStep)
			return m, tea.Batch(tick(), cmd)
		}
	}

	return m, nil
}

func (m Steps) ViewOne(index int) string {
	var (
		space   string
		content strings.Builder
	)

	step := m.steps[index]
	style := lipgloss.NewStyle().
		Margin(0, 1).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		Foreground(lipgloss.Color("#666")).BorderForeground(lipgloss.Color("#aaa"))

	command := fmt.Sprintf("%s %s", step.state, step.command)
	lastly := step.duration.Round(time.Millisecond).String()

	if m.viewportWidth > 0 {
		space = strings.Repeat(" ", m.viewportWidth-lipgloss.Width(command)-lipgloss.Width(lastly)-6)
	} else {
		space = ""
	}
	content.WriteString(lipgloss.JoinHorizontal(lipgloss.Center, command, space, lastly))
	if step.state == Exited1 {
		content.WriteString("\n")
		content.WriteString(strings.TrimSpace(step.output))
	}
	return style.Render(content.String())
}

func (m Steps) View() string {
	var content strings.Builder
	for index := range m.steps {
		content.WriteString(m.ViewOne(index) + "\n")
	}
	return content.String()
}
