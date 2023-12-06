package helpers

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

type interactiveTerminal struct {
	program *tea.Program
}

type (
	errMsg error
)

func NewInteractiveTerminal(configs []string) interactiveTerminal {
	return interactiveTerminal{
		program: tea.NewProgram(initialModel(configs)),
	}
}

func (ic interactiveTerminal) Run() error {
	_, err := ic.program.Run()
	return err
}

type model struct {
	configMap map[string]*textinput.Model
	configs   []string
	cursor    int
	err       error
}

func createTextInput(placeholder string) *textinput.Model {
	ti := textinput.New()
	ti.SetValue(placeholder)

	return &ti
}

func initialModel(configs []string) model {
	configMap := map[string]*textinput.Model{}
	for _, config := range configs {
		configMap[config] = createTextInput(viper.GetString(config))
	}

	return model{
		configMap: configMap,
		configs:   configs,
	}
}

func (m model) saveCurrentState() error {
	for config, ti := range m.configMap {
		viper.Set(config, ti.Value())
	}

	return viper.WriteConfig()
}

func (m model) Init() tea.Cmd {
	m.configMap[m.configs[0]].Focus()
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+q", "ctrl+c", "ctrl+z":
			return m, tea.Quit

		// These keys should save and then exit the program.
		case "ctrl+s":
			err := m.saveCurrentState()
			if err != nil {
				m.err = err
			}
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.configMap[m.configs[m.cursor]].Blur()
				m.cursor--
				m.configMap[m.configs[m.cursor]].Focus()
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.configMap)-1 {
				m.configMap[m.configs[m.cursor]].Blur()
				m.cursor++
				m.configMap[m.configs[m.cursor]].Focus()
			}
		}
	case errMsg:
		m.err = msg
		return m, nil
	}
	cm := m.configMap[m.configs[m.cursor]]
	newCm, cmd := cm.Update(msg)
	m.configMap[m.configs[m.cursor]] = &newCm

	// Return the updated model to the Bubble Tea runtime for processing.
	return m, cmd
}

func (m model) View() string {
	// The header
	s := "Update configs\n\n"

	// Iterate over our configs
	for _, config := range m.configs {
		ti := m.configMap[config]

		// Render the row
		s += fmt.Sprintf("[%s]: %s\n", config, ti.View())
	}

	// The footer
	s += "\nPress ctrl+[q/c/z] to quit, ctrl+s to save.\n"

	// Send the UI for rendering
	return s
}
