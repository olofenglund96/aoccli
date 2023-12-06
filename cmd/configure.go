package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aoccli/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
			term := helpers.NewInteractiveTerminal([]string{"domain", "year", "day", "session-token", "root-dir", "python-exec"})
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

		if day != dayOfMonth {
			fmt.Printf("Day in config '%s' is not today (%s), do you wish to change day to today? (Y/n)\n", day, dayOfMonth)
			var choice string
			fmt.Scanln(&choice)

			if strings.Contains("yY", choice) {
				fmt.Printf("Updating day to %s\n", dayOfMonth)
				day = dayOfMonth
				viper.Set("day", day)
			}
		}
		fmt.Println("== Current Configuration ==")
		fmt.Printf("Domain: %s\n", helpers.GetViperValueEnsureSet("domain"))
		fmt.Printf("Year: %s\n", helpers.GetViperValueEnsureSet("year"))
		fmt.Printf("Day: %s\n", day)
		fmt.Printf("Session token: %s\n", helpers.GetViperValueEnsureSet("session-token"))
		fmt.Printf("Root directory: %s\n", helpers.GetViperValueEnsureSet("root-dir"))
		fmt.Printf("Python executable path: %s\n", helpers.GetViperValueEnsureSet("python-exec"))

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
}
