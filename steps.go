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

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(_ time.Time) tea.Msg {
		return tickMsg{}
	})
}

func newStep(command string) Step {
	return Step{
		command:   command,
		duration:  0,
		output:    "",
		startedAt: time.Now(),
		state:     Pending,
	}
}

func (m Steps) start(index int) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("bash", "-c", m.steps[index].command) //nolint:gosec

		output, err := cmd.Output()

		return exitMsg{
			output: string(output),
			err:    err,
		}
	}
}

func (m Steps) reset() (Steps, tea.Cmd) {
	for i := range m.steps {
		m.currentStep = 0
		m.steps[i].state = Pending
		m.steps[i].duration = time.Second * 0
	}

	return m, m.Init()
}

func (m Steps) Init() tea.Cmd {
	return func() tea.Msg {
		return startMsg{id: 0}
	}
}

func (m Steps) Update(msg tea.Msg) (Steps, tea.Cmd) { //nolint:cyclop
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tickMsg:
		if m.steps[m.currentStep].state == Started {
			m.steps[m.currentStep].duration = time.Since(m.steps[m.currentStep].startedAt).Round(time.Millisecond)

			return m, tick()
		}
	case startMsg:
		m.steps[msg.id].state = Started
		m.steps[msg.id].startedAt = time.Now()
		cmd := m.start(msg.id)

		return m, tea.Batch(tick(), cmd)
	case exitMsg:
		if msg.err != nil {
			m.steps[m.currentStep].state = Exited1

			m.steps[m.currentStep].output = msg.output

			if len(msg.output) == 0 {
				m.steps[m.currentStep].output = msg.err.Error()
			}

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
			start := func() tea.Msg {
				return startMsg{id: m.currentStep}
			}

			return m, tea.Batch(tick(), start)
		}
	case tea.KeyMsg:
		k := msg.String()
		if k == "r" {
			m, cmd = m.reset()

			return m, cmd
		}
	}

	return m, nil
}

func (m Steps) ViewOne(index int) string {
	var (
		space   string
		lastly  string
		content strings.Builder
	)

	step := m.steps[index]
	style := lipgloss.NewStyle().
		Margin(0, 1).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		Foreground(lipgloss.Color("#666")).BorderForeground(lipgloss.Color("#aaa"))

	command := fmt.Sprintf("%s %s", step.state, step.command)

	if step.state == Skipped {
		lastly = "(skipped)"
	} else {
		lastly = step.duration.Round(time.Millisecond).String()
	}

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
