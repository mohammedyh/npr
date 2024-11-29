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

var (
	warningStyle = lipgloss.NewStyle().Margin(1, 2).Foreground(lipgloss.Color("222"))
	errorStyle   = lipgloss.NewStyle().Margin(1, 2).Foreground(lipgloss.Color("161"))
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

var lockfilesToPackageManagers = map[string]string{
	"pnpm-lock.yaml":    "pnpm",
	"package-lock.json": "npm",
	"bun.lockb":         "bun",
	"yarn.lock":         "yarn",
}
var packageManager string

func detectPackageManager() string {
	var lockfiles []string

	cwd, err := os.Getwd()

	if err != nil {
		printErrorFatal("Unable to get current directory", err)
	}

	files, err := os.ReadDir(cwd)

	if err != nil {
		printErrorFatal("Unable to read contents of current directory", err)
	}

	for _, file := range files {
		if _, setInMap := lockfilesToPackageManagers[file.Name()]; setInMap {
			lockfiles = append(lockfiles, file.Name())
		}

		if !file.IsDir() {
			switch file.Name() {
			case "pnpm-lock.yaml":
				packageManager = "pnpm"
			case "package-lock.json":
				packageManager = "npm"
			case "bun.lockb":
				packageManager = "bun"
			case "yarn.lock":
				packageManager = "yarn"
			}
		}
	}

	if len(lockfiles) > 1 {
		multipeLockfilesErr := errors.New("- " + strings.Join(lockfiles, "\n- "))
		printErrorFatal("Found multiple lockfiles", multipeLockfilesErr)
	}

	if packageManager == "" {
		packageManager = "npm"
	}
	return packageManager
}

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

type model struct {
	list list.Model
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

		if msg.String() == "enter" {
			script, _ := m.list.SelectedItem().(Script)
			return m, runScript(packageManager, script.name)
		}
	case tea.WindowSizeMsg:
		h, v := warningStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case CommandExecuted:
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return lipgloss.NewStyle().Margin(1, 2).Render(m.list.View())
}

func main() {
	jsonData, err := os.ReadFile("package.json")

	if err != nil {
		printErrorFatal("package.json not found", err)
	}

	detectPackageManager()

	var parsedJson PackageJsonFields

	parseErr := json.Unmarshal(jsonData, &parsedJson)

	if parseErr != nil {
		printErrorFatal("Unable to parse package.json", parseErr)
	}

	scriptsList := parsedJson.Scripts
	depsList := parsedJson.Dependencies
	devDepsList := parsedJson.DevDependencies

	if len(scriptsList) == 0 {
		fmt.Println(errorStyle.Render("No scripts to run"))
		os.Exit(1)
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
