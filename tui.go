package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	title, fullName, shortTitle, desc, cloneUrl string
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)
var items = []list.Item{}

func (i item) Title() string       { return i.title + "   " }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	selecedItem list.Item
	list        list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == " " {
			println("quit")
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func DisplayTui(repos []Repository) item {
	for _, repo := range repos {
		items = append(items,
			item{
				title:      fmt.Sprintf("%s    %d     %d     %d     %s", repo.FullName, repo.Stars, repo.Forks, repo.Issues, repo.Language),
				desc:       repo.Description,
				cloneUrl:   repo.CloneUrl,
				shortTitle: repo.Name,
				fullName:   repo.FullName,
			})
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "<SPACE> to clone the selected repository"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return m.list.SelectedItem().(item)
}
