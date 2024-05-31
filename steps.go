package main

import (
	"fmt"
	"log"
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
	command string
	counter int
}

type tickMsg struct {
}

type startMsg struct {
	id int
}

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

func (m Steps) start(index int) tea.Cmd {
	return func() tea.Msg {
		return startMsg{id: index}
	}
}

func (m Steps) Init() tea.Cmd {
	return tea.Batch(tick(true), m.start(0))
}

func (m Steps) Update(msg tea.Msg) (Steps, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		log.Println("update currentstep", m.currentStep)
		m.steps[m.currentStep].counter++
		return m, tick()
	case startMsg:
		m.currentStep = msg.id
	}
	return m, nil
}

func (m Steps) ViewOne(index int) string {
	var (
		icon  string
		space string
	)
	step := m.steps[index]
	style := lipgloss.NewStyle().
		Margin(0, 1).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		Foreground(lipgloss.Color("#666")).BorderForeground(lipgloss.Color("#aaa"))

	if index == m.currentStep {
		icon = "🟡"
	} else {
		icon = "🔜"
	}
	content := fmt.Sprintf("%s %s", icon, step.command)
	lastly := fmt.Sprintf("%d", step.counter)

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
