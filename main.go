package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PackageJsonFields struct {
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

type Script struct {
	name, command string
}

func (s Script) Title() string       { return s.name }
func (s Script) Description() string { return s.command }
func (s Script) FilterValue() string { return s.name }

type CommandExecuted struct{}

var packageManager string

func installDependencies(packageManager string) {
	contents, err := os.ReadDir("node_modules")
	if err != nil || len(contents) == 0 {
		if os.IsPermission(err) {
			printErrorFatal("Unable to read node_modules directory", err)
		}

		command := exec.Command(packageManager, "install")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		fmt.Printf("Installing packages using %v\n", packageManager)
		err := command.Run()
		if err != nil {
			printErrorFatal("Error running command", err)
		}
	}
}

func runScript(packageManager, scriptName string) tea.Cmd {
	command := exec.Command(packageManager, "run", scriptName)
	return tea.ExecProcess(command, func(err error) tea.Msg {
		if err != nil {
			return tea.Quit()
		}
		return CommandExecuted{}
	})
}

func main() {
	jsonData, err := os.ReadFile("package.json")
	if err != nil {
		printErrorFatal("package.json not found", err)
	}

	packageManager = detectPackageManager()
	var parsedJson PackageJsonFields

	parseErr := json.Unmarshal(jsonData, &parsedJson)
	if parseErr != nil {
		printErrorFatal("Unable to parse package.json", parseErr)
	}

	scriptsList := parsedJson.Scripts
	if len(scriptsList) == 0 {
		printErrorFatal("No scripts to run", nil)
	}

	if len(parsedJson.Dependencies) > 0 || len(parsedJson.DevDependencies) > 0 {
		installDependencies(packageManager)
	}

	var items []list.Item

	for name, command := range scriptsList {
		items = append(items, Script{name, command})
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].(Script).name < items[j].(Script).name
	})

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Scripts to Run"
	m.list.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fff")).
		Background(lipgloss.Color("#bc54c4")).
		Padding(0, 1)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		printErrorFatal("Error running program", err)
	}
}
