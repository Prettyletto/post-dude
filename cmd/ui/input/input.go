package input

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type State int

const (
	Ready State = iota
	Success
)

type InputModel struct {
	State       State
	TextInput   textinput.Model
	SuccessMsg  string
	SuccessTime time.Time
}

func New() InputModel {
	ti := textinput.New()
	ti.Placeholder = "Enter Collection Name"
	ti.Focus()
	ti.CharLimit = 154
	ti.Width = 30

	return InputModel{
		State:     Ready,
		TextInput: ti,
	}

}

func (m InputModel) Init() tea.Cmd {
	return nil
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.State {
	case Ready:
		var cmd tea.Cmd
		m.TextInput, cmd = m.TextInput.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				if m.TextInput.Value() != "" {
					m.SuccessMsg = fmt.Sprintf("Success! Added Collection: %s", m.TextInput.Value())
					m.SuccessTime = time.Now().Add(1 * time.Second)
					m.State = Success
				}
				return m, nil
			case "esc":
				return m, tea.Quit
			}
		}
		return m, cmd
	case Success:
		if time.Now().After(m.SuccessTime) {
			return m, func() tea.Msg {
				return DoneMsg(m.TextInput.Value())
			}
		}
		return m, nil
	default:
		return m, nil
	}
}

func (m InputModel) View() string {
	switch m.State {
	case Ready:
		return "Enter collection name:\n\n" + m.TextInput.View() + "\n\n(Enter to submit, Esc to cancel)"
	case Success:
		return m.SuccessMsg
	default:
		return ""
	}
}

type DoneMsg string
