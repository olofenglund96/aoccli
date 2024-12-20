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

var configKeys = []string{"year", "day", "session-token", "root-dir", "python-exec", "leaderboard"}

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
		helpers.HandleErr(err)

		viper.Set(cmdString, value)
		fmt.Printf("Set %s to '%s'\n", cmdString, value)
	}
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure the CLI",
	Long:  `Set the year and session token to talk to AOC`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 {
			term := helpers.NewInteractiveTerminal(configKeys)
			err := term.Run()
			helpers.HandleErr(err)
			return
		}

		for _, configKey := range configKeys {
			saveIfChanged(cmd, configKey)
		}

		day := helpers.GetViperValueEnsureSet("day")
		currentTime := time.Now()

		dayOfMonth := strconv.Itoa(currentTime.Day())
		rd, err := cmd.LocalFlags().GetBool("refresh-day")
		helpers.HandleErr(err)

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
		helpers.HandleErr(viper.WriteConfig())
	},
}


func init() {
	rootCmd.AddCommand(configureCmd)

	viper.SetDefault("year", fmt.Sprintf("%d", time.Now().Local().Year()))
	viper.SetDefault("day", fmt.Sprintf("%02d", time.Now().Local().Day())) // Ensure it's formatted as a string

	sessionToken := viper.GetString("session-token")
	rootDir := viper.GetString("root-dir")
	pythonExecutable := viper.GetString("python-exec")
	leaderboard := viper.GetString("leaderboard")

	year := viper.GetString("year")
	day := viper.GetString("day")

	configureCmd.Flags().StringP("year", "y", year, "Selected year")
	configureCmd.Flags().String("day", day, "Selected day")
	configureCmd.Flags().StringP("session-token", "t", sessionToken, "Session token copied from AOC")
	configureCmd.Flags().StringP("root-dir", "r", rootDir, "Root directory for the problems")
	configureCmd.Flags().StringP("python-exec", "p", pythonExecutable, "Path to the python executable")
	configureCmd.Flags().BoolP("refresh-day", "u", false, "Simply refresh the day")
	configureCmd.Flags().StringP("leaderboard", "l", leaderboard, "Private leaderboard id")
}