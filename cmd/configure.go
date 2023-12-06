package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aoccli/helpers"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configKeys = []string{"domain", "year", "day", "session-token", "root-dir", "python-exec"}

const (
	purple    = lipgloss.Color("99")
	gray      = lipgloss.Color("245")
	lightGray = lipgloss.Color("241")
)

func printConfigurationTable() {
	fmt.Println("== Current Configuration ==")
	re := lipgloss.NewRenderer(os.Stdout)

	var (
		// HeaderStyle is the lipgloss style used for the table headers.
		HeaderStyle = re.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
		// CellStyle is the base lipgloss style used for the table rows.
		CellStyle = re.NewStyle().Padding(0, 1)
		// OddRowStyle is the lipgloss style used for odd-numbered table rows.
		OddRowStyle = CellStyle.Copy().Foreground(gray)
		// EvenRowStyle is the lipgloss style used for even-numbered table rows.
		EvenRowStyle = CellStyle.Copy().Foreground(lightGray)
		// BorderStyle is the lipgloss style used for the table border.
		BorderStyle = lipgloss.NewStyle().Foreground(purple)
	)

	t := table.New().Border(lipgloss.ThickBorder()).
		BorderStyle(BorderStyle).
		StyleFunc(func(row, col int) lipgloss.Style {
			var style lipgloss.Style

			switch {
			case row == 0:
				return HeaderStyle
			case row%2 == 0:
				style = EvenRowStyle
			default:
				style = OddRowStyle
			}

			return style
		}).
		Headers("Config", "Value")

	for _, ck := range configKeys {
		t.Row(ck, helpers.GetViperValueEnsureSet(ck))
	}

	fmt.Println(t)
}

func saveIfChanged(cmd *cobra.Command, cmdString string) {
	if cmd.LocalFlags().Changed(cmdString) {
		value, err := cmd.LocalFlags().GetString(cmdString)
		cobra.CheckErr(err)

		viper.Set(cmdString, value)
		fmt.Printf("Set %s to '%s'\n", cmdString, value)
	}
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure the CLI",
	Long:  `Set the domain, year and session token to talk to AOC`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 {
			term := helpers.NewInteractiveTerminal(configKeys)
			err := term.Run()
			cobra.CheckErr(err)
			return
		}

		saveIfChanged(cmd, "domain")
		saveIfChanged(cmd, "year")
		saveIfChanged(cmd, "day")
		saveIfChanged(cmd, "session-token")
		saveIfChanged(cmd, "root-dir")
		saveIfChanged(cmd, "python-exec")

		day := helpers.GetViperValueEnsureSet("day")
		currentTime := time.Now()

		dayOfMonth := strconv.Itoa(currentTime.Day())
		rd, err := cmd.LocalFlags().GetBool("refresh-day")
		cobra.CheckErr(err)

		if rd {
			if day != dayOfMonth {
				fmt.Printf("Day in config '%s' is not today (%s), changing day to today\n", day, dayOfMonth)
				day = dayOfMonth
				viper.Set("day", day)
			} else {
				fmt.Printf("Day in config '%s' is already set to today (%s), not changing day..\n", day, dayOfMonth)
			}
		}
		printConfigurationTable()
		cobra.CheckErr(viper.WriteConfig())
	},
}

func parseInt(value string) int {
	var err error
	intValue := 0
	if value != "" {
		intValue, err = strconv.Atoi(value)
		cobra.CheckErr(err)
	}

	return intValue
}

func init() {
	rootCmd.AddCommand(configureCmd)

	domain := viper.GetString("domain")
	sessionToken := viper.GetString("session-token")
	rootDir := viper.GetString("root-dir")
	pythonExecutable := viper.GetString("python-exec")

	year := parseInt(viper.GetString("year"))
	day := parseInt(viper.GetString("day"))

	configureCmd.Flags().StringP("domain", "d", domain, "Domain of AOC")
	configureCmd.Flags().IntP("year", "y", year, "Selected year")
	configureCmd.Flags().Int("day", day, "Selected day")
	configureCmd.Flags().StringP("session-token", "t", sessionToken, "Session token copied from AOC")
	configureCmd.Flags().StringP("root-dir", "r", rootDir, "Root directory for the problems")
	configureCmd.Flags().StringP("python-exec", "p", pythonExecutable, "Path to the python executable")
	configureCmd.Flags().BoolP("refresh-day", "u", false, "Simply refresh the day")
}
