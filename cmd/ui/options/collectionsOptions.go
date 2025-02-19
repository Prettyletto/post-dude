package options

import (
	"log"

	"github.com/Prettyletto/post-dude/cmd/ui/client"
	editcollection "github.com/Prettyletto/post-dude/cmd/ui/editCollection"
	"github.com/Prettyletto/post-dude/cmd/ui/views"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	optionsState = iota
	editState
)

type Model struct {
	State      int
	list       list.Model
	collection *views.Collection
	editModel  *editcollection.EditModel
}

type OptionItem struct {
	title       string
	description string
}

func (o OptionItem) Title() string       { return o.title }
func (o OptionItem) Description() string { return o.description }
func (o OptionItem) FilterValue() string { return "" }

type BackMsg struct{}

type DeletedMsg struct{}

type UpdatedMsg struct {
	Updated string
}

type ErrorMsg struct{}

func deleteCollectionCmd(id int) tea.Cmd {
	return func() tea.Msg {
		err := client.DeleteCollection(id)
		if err != nil {
			log.Fatal(err)
			return ErrorMsg{}
		}
		return DeletedMsg{}
	}
}

func updateCollectionCmd(id int, collection views.Collection) tea.Cmd {
	return func() tea.Msg {
		err := client.UpdateCollection(id, collection)
		if err != nil {
			log.Fatal(err)
			return ErrorMsg{}
		}
		return UpdatedMsg{Updated: collection.Name}
	}

}

func New(collection *views.Collection) *Model {
	items := []list.Item{
		OptionItem{title: "Entities", description: "Create a new entity"},
		OptionItem{title: "Edit Collection", description: "Edit Current Collection"},
		OptionItem{title: "Delete Collection", description: "Delete current Collection"},
		OptionItem{title: "Back", description: "Goes back to collections menu"},
	}

	const listWidth = 30
	l := list.New(items, list.NewDefaultDelegate(), listWidth, 35)
	l.Title = collection.Name
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	l.Styles.HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	l.Styles.PaginationStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	return &Model{
		collection: collection,
		list:       l,
		State:      optionsState,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.State {
	case optionsState:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				switch sel := m.list.SelectedItem().(OptionItem); sel.Title() {
				case "Edit Collection":
					m.State = editState
					newModel := editcollection.New(m.collection)
					m.editModel = &newModel
					return m, nil
				case "Delete Collection":
					return m, deleteCollectionCmd(m.collection.ID)
				case "Back":
					return m, func() tea.Msg { return BackMsg{} }
				}
			case "esc", "backspace":
				return m, func() tea.Msg { return BackMsg{} }
			}
		}
	case editState:
		newEditModel, cmd := m.editModel.Update(msg)

		if em, ok := newEditModel.(editcollection.EditModel); ok {
			m.editModel = &em
		}
		switch msg := msg.(type) {
		case editcollection.DoneMsg:
			m.collection.Name = string(msg)
			m.editModel = nil
			m.State = optionsState
			m.list.Title = m.collection.Name
			return m, updateCollectionCmd(m.collection.ID, *m.collection)
		}

		return m, cmd
	}
	m.list, cmd = m.list.Update(msg)
	return m, cmd

}

func (m *Model) View() string {
	switch m.State {
	case optionsState:
		return m.list.View()
	case editState:
		return m.editModel.View()
	}
	return "Unkown state"
}
