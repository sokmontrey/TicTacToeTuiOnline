package page

import tea "github.com/charmbracelet/bubbletea"

type PageManager struct {
	currentPage tea.Model
}

func NewPageManager(currentPage tea.Model) PageManager {
	return PageManager{
		currentPage: currentPage,
	}
}

func (m PageManager) Init() tea.Cmd {
	return nil
}

func (m PageManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.currentPage, cmd = m.currentPage.Update(msg)
	return m, cmd
}

func (m PageManager) View() string {
	return m.currentPage.View()
}
