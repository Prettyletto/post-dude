package collections

import (
	"github.com/Prettyletto/post-dude/cmd/ui/options"
	"github.com/Prettyletto/post-dude/cmd/ui/views"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	CollectionsState = iota
	OptionsState
)

type CollectionItem views.Collection

func (c CollectionItem) Title() string       { return c.Name }
func (c CollectionItem) Description() string { return "" }
func (c CollectionItem) FilterValue() string { return c.Name }

type Model struct {
	State        int
	list         list.Model
	optionsModel *options.Model
	collection   *views.Collection
}

type BackMsg struct{}

func New(collections []views.Collection) *Model {
	items := make([]list.Item, len(collections))
	for i, collection := range collections {
		items[i] = CollectionItem(collection)
	}

	const listWidth = 30
	l := list.New(items, list.NewDefaultDelegate(), listWidth, 35)
	l.Title = "Collections"
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	l.Styles.HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	l.Styles.PaginationStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	return &Model{
		list:  l,
		State: CollectionsState,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.State {
	case CollectionsState:
		return m.updateCollectionsState(msg)
	case OptionsState:
		if m.optionsModel == nil {
			return m, nil
		}
		var cmd tea.Cmd
		var newModel tea.Model
		newModel, cmd = m.optionsModel.Update(msg)

		if optionModel, ok := newModel.(*options.Model); ok {
			m.optionsModel = optionModel
		}

		if _, ok := msg.(options.BackMsg); ok {
			m.optionsModel = nil
			m.State = CollectionsState
		}

		if _, ok := msg.(options.DeletedMsg); ok {
			m.optionsModel = nil
			m.State = CollectionsState
			m.list.RemoveItem(m.list.Index())
		}

		if updated, ok := msg.(options.UpdatedMsg); ok {
			m.list.SetItem(m.list.Index(), CollectionItem{Name: updated.Updated})
		}

		return m, cmd
	}
	return m, nil
}

func (m *Model) updateCollectionsState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.State = OptionsState
			sel := m.list.SelectedItem()

			if collItem, ok := sel.(CollectionItem); ok {
				m.collection = (*views.Collection)(&collItem)
				m.optionsModel = options.New(m.collection)
			}

			return m, nil
		case "esc", "backspace":
			return m, func() tea.Msg { return BackMsg{} }
		}
	}
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	switch m.State {
	case CollectionsState:
		return m.list.View()
	case OptionsState:
		return m.optionsModel.View()
	}
	return "Invalid View"
}
