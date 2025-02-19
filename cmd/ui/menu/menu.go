package menu

import (
	"log"

	"github.com/Prettyletto/post-dude/cmd/ui/client"
	"github.com/Prettyletto/post-dude/cmd/ui/collections"
	"github.com/Prettyletto/post-dude/cmd/ui/input"
	"github.com/Prettyletto/post-dude/cmd/ui/views"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type State int

const (
	MainMenuState = iota
	InputState
	CollectionsState
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type MenuStyle struct {
	titleStyle        lipgloss.Style
	itemStyle         lipgloss.Style
	selectedItemStyle lipgloss.Style
	paginationStyle   lipgloss.Style
	helpStyle         lipgloss.Style
	quitTextStyle     lipgloss.Style
}

func newMenuStyle() MenuStyle {
	return MenuStyle{
		titleStyle:        lipgloss.NewStyle().MarginLeft(2),
		itemStyle:         lipgloss.NewStyle().PaddingLeft(4),
		selectedItemStyle: lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")),
		paginationStyle:   list.DefaultStyles().PaginationStyle.PaddingLeft(4),
		helpStyle:         list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1),
		quitTextStyle:     lipgloss.NewStyle().Margin(1, 0, 2, 4),
	}

}

type MenuItem struct {
	title, description string
}

func (m MenuItem) Title() string { return m.title }

func (m MenuItem) Description() string { return m.description }

func (m MenuItem) FilterValue() string { return "" }

type CollectionsPostMsg struct {
	Success bool
	Error   error
}

type CollectionsFetchedMsg struct {
	Collections []views.Collection
}

func postCollectionsCmd(newCollection views.Collection) tea.Cmd {
	return func() tea.Msg {
		err := client.PostCollection(newCollection)
		if err != nil {
			log.Fatal(err)
			return CollectionsPostMsg{Success: false, Error: err}
		}
		return CollectionsPostMsg{Success: true}
	}
}

func fetchCollectionsCmd() tea.Cmd {
	return func() tea.Msg {
		collections, err := client.FetchCollections()
		if err != nil {
			log.Fatal(err)
			return CollectionsFetchedMsg{Collections: nil}
		}
		return CollectionsFetchedMsg{Collections: collections}
	}
}

type MenuModel struct {
	State            State
	collections      []views.Collection
	menuStyle        MenuStyle
	MainMenu         list.Model
	InputModel       tea.Model
	CollectionsModel *collections.Model
}

func New() *MenuModel {
	items := []list.Item{
		MenuItem{title: "Create a collection", description: "Create a collection of new requests (this should be your entire application)"},
		MenuItem{title: "Collections", description: "Your applications collections should be here"},
		MenuItem{title: "Quit", description: "Quit the application "},
	}
	const defaultWidth = 50
	l := list.New(items, list.NewDefaultDelegate(), defaultWidth, 24)
	l.Title = "Post Dude Menu"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	styles := newMenuStyle()
	l.Styles.Title = styles.titleStyle
	l.Styles.PaginationStyle = styles.paginationStyle
	l.Styles.HelpStyle = styles.helpStyle

	return &MenuModel{
		State:     MainMenuState,
		MainMenu:  l,
		menuStyle: styles,
	}
}

func (m *MenuModel) Init() tea.Cmd {
	return nil
}

func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.State {
	case MainMenuState:
		return m.updateMainMenu(msg)
	case InputState:
		return m.updateInputScreen(msg)
	case CollectionsState:
		switch msg := msg.(type) {
		case CollectionsFetchedMsg:
			m.collections = msg.Collections
			m.CollectionsModel = collections.New(m.collections)
			return m, nil

		default:
			var cmd tea.Cmd
			var newModel tea.Model
			newModel, cmd = m.CollectionsModel.Update(msg)

			if collectionsModel, ok := newModel.(*collections.Model); ok {
				m.CollectionsModel = collectionsModel
			}
			switch msg.(type) {
			case collections.BackMsg:
				m.State = MainMenuState
				m.CollectionsModel = nil
			}
			return m, cmd
		}

	default:
		return m, nil
	}
}

func (m *MenuModel) updateMainMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			switch sel := m.MainMenu.SelectedItem().(MenuItem); sel.title {
			case "Create a collection":
				m.State = InputState
				inputModel := input.New()
				m.InputModel = &inputModel
				return m, nil
			case "Collections":
				m.State = CollectionsState
				return m, fetchCollectionsCmd()
			case "Quit":
				return m, tea.Quit
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	m.MainMenu, cmd = m.MainMenu.Update(msg)
	return m, cmd
}

func (m *MenuModel) updateInputScreen(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.InputModel == nil {
		return m, nil
	}
	switch msg := msg.(type) {
	case input.DoneMsg:
		newCollection := views.Collection{Name: string(msg)}
		m.State = MainMenuState
		m.InputModel = nil
		return m, postCollectionsCmd(newCollection)

	}
	var cmd tea.Cmd
	m.InputModel, cmd = m.InputModel.Update(msg)
	return m, cmd
}

func (m MenuModel) View() string {
	switch m.State {
	case MainMenuState:
		return m.MainMenu.View() + "\n"
	case InputState:
		return m.InputModel.View()
	case CollectionsState:
		if m.CollectionsModel == nil {
			return "Loading..."
		}
		return m.CollectionsModel.View()
	default:
		return "Unknown state"
	}
}
