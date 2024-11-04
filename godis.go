package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// Model represents the state of the application
type Model struct {
	disassembly1 string
	disassembly2 string
	cursorY1     int
	cursorY2     int
	width        int
	height       int
	viewport     viewport.Model
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	m.viewport = viewport.New(m.width, m.height-4)
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			if m.cursorY1 < len(strings.Split(m.disassembly1, "\n"))-1 {
				m.cursorY1++
			}
		case "k":
			if m.cursorY1 > 0 {
				m.cursorY1--
			}
		case "tab":
			if m.cursorY2 < len(strings.Split(m.disassembly2, "\n"))-1 {
				m.cursorY2++
			}
		case "shift+tab":
			if m.cursorY2 > 0 {
				m.cursorY2--
			}
		case "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = m.width
		m.viewport.Height = m.height - 4
	}
	return m, nil
}

// View renders the model
func (m Model) View() string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(m.disassembly1, m.disassembly2, false)

	var leftColumn, rightColumn strings.Builder
	equalStyle := lipgloss.NewStyle()
	deleteStyle := lipgloss.NewStyle().Background(lipgloss.Color("9"))  // Red background
	insertStyle := lipgloss.NewStyle().Background(lipgloss.Color("10")) // Green background

	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffEqual:
			leftColumn.WriteString(equalStyle.Render(diff.Text))
			rightColumn.WriteString(equalStyle.Render(diff.Text))
		case diffmatchpatch.DiffDelete:
			leftColumn.WriteString(deleteStyle.Render(diff.Text))
		case diffmatchpatch.DiffInsert:
			rightColumn.WriteString(insertStyle.Render(diff.Text))
		}
	}

	halfWidth := m.width / 2
	leftStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Width(halfWidth - 2).Height(m.viewport.Height)
	rightStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Width(halfWidth - 2).Height(m.viewport.Height)

	header := lipgloss.NewStyle().
		Bold(true).
		Border(lipgloss.NormalBorder()).
		Align(lipgloss.Center).
		Width(halfWidth).
		MarginLeft((m.width - halfWidth) / 2).
		Render("Godis")

	content := lipgloss.JoinHorizontal(lipgloss.Top, leftStyle.Render(leftColumn.String()), rightStyle.Render(rightColumn.String()))
	m.viewport.SetContent(content)
	m.viewport.GotoBottom() // Scroll to the bottom

	return lipgloss.JoinVertical(lipgloss.Top, header, m.viewport.View())
}

func main() {
	disassembly1 := "push ebp\nmov ebp, esp\nsub esp, 0x10\nmov eax, 0x2\nmov ebx, 0x1\nadd eax, ebx\nmov [ebp-0x4], eax\nmov esp, ebp\npop ebp\nret"
	disassembly2 := "push ebp\nmov ebp, esp\nsub esp, 0x10\nmov eax, 0x5\nmov ebx, 0x3\nadd eax, ebx\nmov [ebp-0x4], eax\nmov esp, ebp\npop ebp\n"

	model := Model{
		disassembly1: disassembly1,
		disassembly2: disassembly2,
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
