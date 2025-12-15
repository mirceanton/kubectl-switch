package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SelectModel represents a selection list component
type SelectModel struct {
	message         string
	options         []string
	filteredOptions []string
	filter          string
	current         string
	cursor          int
	pageSize        int
	offset          int
	width           int
	selected        string
	quitting        bool
	aborted         bool
}

// Styles for the select component
var (
	cursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))   // cyan
	normalStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("250")) // light gray
	currentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))   // magenta

	// Prompt styles
	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7")).
			Background(lipgloss.Color("4")).
			Bold(true).
			Padding(0, 1)
	filterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("3")).
			Italic(true)
	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))
)

// fuzzyMatch performs a simple fuzzy match - checks if all characters in pattern
// appear in str in order (case-insensitive)
func fuzzyMatch(str, pattern string) bool {
	if pattern == "" {
		return true
	}
	str = strings.ToLower(str)
	pattern = strings.ToLower(pattern)

	patternIdx := 0
	for i := 0; i < len(str) && patternIdx < len(pattern); i++ {
		if str[i] == pattern[patternIdx] {
			patternIdx++
		}
	}
	return patternIdx == len(pattern)
}

// Key bindings
type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Right  key.Binding
	Enter  key.Binding
	Quit   key.Binding
	Escape key.Binding
	PgUp   key.Binding
	PgDown key.Binding
	Home   key.Binding
	End    key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("up", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("down", "move down"),
	),
	Right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("right", "use as filter"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "clear filter/quit"),
	),
	PgUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "page up"),
	),
	PgDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdown", "page down"),
	),
	Home: key.NewBinding(
		key.WithKeys("home"),
		key.WithHelp("home", "go to start"),
	),
	End: key.NewBinding(
		key.WithKeys("end"),
		key.WithHelp("end", "go to end"),
	),
}

// NewSelectModel creates a new selection model
func NewSelectModel(message string, options []string, current string, pageSize int) SelectModel {
	if pageSize <= 0 {
		pageSize = 10
	}
	// Initialize filteredOptions as a copy of options
	filteredOptions := make([]string, len(options))
	copy(filteredOptions, options)

	return SelectModel{
		message:         message,
		options:         options,
		filteredOptions: filteredOptions,
		filter:          "",
		current:         current,
		cursor:          0,
		pageSize:        pageSize,
		offset:          0,
		width:           80,
	}
}

// Init implements tea.Model
func (m SelectModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m SelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.aborted = true
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, keys.Escape):
			if m.filter != "" {
				m.filter = ""
				m.updateFilter()
			} else {
				m.aborted = true
				m.quitting = true
				return m, tea.Quit
			}

		case key.Matches(msg, keys.Enter):
			if len(m.filteredOptions) > 0 {
				m.selected = m.filteredOptions[m.cursor]
				m.quitting = true
				return m, tea.Quit
			}

		case key.Matches(msg, keys.Right):
			if len(m.filteredOptions) > 0 {
				m.filter = m.filteredOptions[m.cursor]
				m.updateFilter()
			}

		case key.Matches(msg, keys.Up):
			if len(m.filteredOptions) > 0 {
				if m.cursor > 0 {
					m.cursor--
				} else {
					m.cursor = len(m.filteredOptions) - 1
				}
				m.adjustOffset()
			}

		case key.Matches(msg, keys.Down):
			if len(m.filteredOptions) > 0 {
				if m.cursor < len(m.filteredOptions)-1 {
					m.cursor++
				} else {
					m.cursor = 0
				}
				m.adjustOffset()
			}

		case key.Matches(msg, keys.PgUp):
			m.cursor -= m.pageSize
			if m.cursor < 0 {
				m.cursor = 0
			}
			m.adjustOffset()

		case key.Matches(msg, keys.PgDown):
			m.cursor += m.pageSize
			if m.cursor >= len(m.filteredOptions) {
				m.cursor = len(m.filteredOptions) - 1
			}
			if m.cursor < 0 {
				m.cursor = 0
			}
			m.adjustOffset()

		case key.Matches(msg, keys.Home):
			m.cursor = 0
			m.adjustOffset()

		case key.Matches(msg, keys.End):
			if len(m.filteredOptions) > 0 {
				m.cursor = len(m.filteredOptions) - 1
			}
			m.adjustOffset()

		default:
			// Handle character input for filtering
			keyStr := msg.String()
			if keyStr == "backspace" {
				if len(m.filter) > 0 {
					m.filter = m.filter[:len(m.filter)-1]
					m.updateFilter()
				}
			} else if len(keyStr) == 1 && keyStr[0] >= 32 && keyStr[0] < 127 {
				// Printable ASCII character
				m.filter += keyStr
				m.updateFilter()
			}
		}
	}

	return m, nil
}

// updateFilter updates the filtered options based on the current filter
func (m *SelectModel) updateFilter() {
	m.filteredOptions = nil
	for _, opt := range m.options {
		if fuzzyMatch(opt, m.filter) {
			m.filteredOptions = append(m.filteredOptions, opt)
		}
	}
	// Reset cursor and offset when filter changes
	m.cursor = 0
	m.offset = 0
}

// adjustOffset ensures the cursor is visible within the page
func (m *SelectModel) adjustOffset() {
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+m.pageSize {
		m.offset = m.cursor - m.pageSize + 1
	}
}

// View implements tea.Model
func (m SelectModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Build left side of header
	leftSide := promptStyle.Render(m.message) + " "
	if m.filter != "" {
		leftSide += filterStyle.Render(m.filter)
	} else {
		leftSide += hintStyle.Render("Type to filter...")
	}

	// Build right side (counter)
	rightSide := hintStyle.Render(fmt.Sprintf("(%d/%d)", m.cursor+1, len(m.filteredOptions)))

	// Calculate padding for right alignment
	leftLen := lipgloss.Width(leftSide)
	rightLen := lipgloss.Width(rightSide)
	padding := m.width - leftLen - rightLen
	if padding < 1 {
		padding = 1
	}

	b.WriteString(leftSide)
	b.WriteString(strings.Repeat(" ", padding))
	b.WriteString(rightSide)
	b.WriteString("\n")

	// Handle empty filtered results
	if len(m.filteredOptions) == 0 {
		b.WriteString(normalStyle.Render("  No matches found"))
		b.WriteString("\n")
		return b.String()
	}

	// Calculate visible range
	end := m.offset + m.pageSize
	if end > len(m.filteredOptions) {
		end = len(m.filteredOptions)
	}

	// Display options
	for i := m.offset; i < end; i++ {
		option := m.filteredOptions[i]
		isCurrent := option == m.current

		if i == m.cursor {
			b.WriteString(cursorStyle.Render("> "))
			b.WriteString(cursorStyle.Render(option))
		} else {
			b.WriteString("  ")
			if isCurrent {
				b.WriteString(currentStyle.Render(option))
			} else {
				b.WriteString(normalStyle.Render(option))
			}
		}
		if isCurrent {
			b.WriteString(currentStyle.Render(" (current)"))
		}
		b.WriteString("\n")
	}

	return b.String()
}

// Selected returns the selected option
func (m SelectModel) Selected() string {
	return m.selected
}

// Aborted returns true if the user aborted the selection
func (m SelectModel) Aborted() bool {
	return m.aborted
}

// Select runs an interactive selection prompt and returns the selected option
func Select(message string, options []string, current string, pageSize int) (string, error) {
	model := NewSelectModel(message, options, current, pageSize)
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run selection: %w", err)
	}

	result := finalModel.(SelectModel)
	if result.Aborted() {
		return "", fmt.Errorf("selection aborted")
	}

	return result.Selected(), nil
}
