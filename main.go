package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Script struct {
	name, command string
}

func (s Script) Title() string       { return s.name }
func (s Script) Description() string { return s.command }
func (s Script) FilterValue() string { return s.name }

func runScript(scriptName string) tea.Cmd {
	command := exec.Command("npm", "run", scriptName)
	return tea.ExecProcess(command, func(err error) tea.Msg {
		return tea.Quit
	})
}

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("np-run")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if msg.String() == "enter" {
			script, _ := m.list.SelectedItem().(Script)
			return m, runScript(script.name)
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

func main() {
	packageJsonContent, readFileErr := os.ReadFile("package.json")

	if readFileErr != nil {
		fmt.Println("package.json not found")
		os.Exit(1)
	}

	var parsedJson map[string]interface{}

	parseErr := json.Unmarshal(packageJsonContent, &parsedJson)
	if parseErr != nil {
		fmt.Println(parseErr.Error())
		os.Exit(1)
	}

	var items []list.Item

	for scriptKey, scriptValue := range parsedJson["scripts"].(map[string]interface{}) {
		items = append(items, Script{name: scriptKey, command: scriptValue.(string)})
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Scripts to Run"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
