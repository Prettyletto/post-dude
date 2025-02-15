package collections

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	MenuState = iota
	CollectionOptionsState
)

type CollectionItem struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (c CollectionItem) Title() string       { return c.Name }
func (c CollectionItem) Description() string { return "" }
func (c CollectionItem) FilterValue() string { return c.Name }

type CollectionOptionItem struct {
	title, description string
}

func (o CollectionOptionItem) Title() string       { return o.title }
func (o CollectionOptionItem) Description() string { return o.description }
func (o CollectionOptionItem) FilterValue() string { return "" }

type BackMsg struct{}

type Model struct {
	State             int
	list              list.Model
	collection        CollectionItem
	collectionOptions list.Model
}

func New(collections []CollectionItem) *Model {
	items := make([]list.Item, len(collections))
	for i, collection := range collections {
		items[i] = collection
	}

	const listWidth = 30
	l := list.New(items, list.NewDefaultDelegate(), listWidth, 35)
	l.Title = "Collections"
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	l.Styles.HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	l.Styles.PaginationStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	return &Model{
		State: MenuState,
		list:  l,
	}
}

func newOptionsList(collection CollectionItem) list.Model {
	options := []list.Item{
		CollectionOptionItem{title: "Add new request", description: "Create a new request"},
		CollectionOptionItem{title: "Edit", description: "Create a new request"},
		CollectionOptionItem{title: "Requests", description: "Create a new request"},
		CollectionOptionItem{title: "Back", description: "Create a new request"},
	}
	const optionsWidth = 30
	optsList := list.New(options, list.NewDefaultDelegate(), optionsWidth, 30)
	optsList.Title = collection.Name
	optsList.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))

	return optsList
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.State {
	case MenuState:
		return m.updateMainMenu(msg)
	case CollectionOptionsState:
		return m.updateCollectionOptions(msg)
	default:
		return m, nil
	}
}

func (m *Model) updateMainMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.collection = m.list.SelectedItem().(CollectionItem)
			m.collectionOptions = newOptionsList(m.collection)
			m.State = CollectionOptionsState
			return m, nil
		case "esc", "backspace":
			return m, func() tea.Msg { return BackMsg{} }
		}

	}
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *Model) updateCollectionOptions(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.collectionOptions, cmd = m.collectionOptions.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace":
			m.State = MenuState
			return m, nil
		}
	}
	return m, cmd
}

func (m *Model) View() string {
	switch m.State {
	case MenuState:
		return m.list.View()
	case CollectionOptionsState:
		return m.collectionOptions.View()
	default:
		return "Unkown State"
	}
}
