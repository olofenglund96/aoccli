package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/aoccli/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func saveIfStringChanged(cmd *cobra.Command, cmdString string) {
	if cmd.LocalFlags().Changed(cmdString) {
		value, err := cmd.LocalFlags().GetString(cmdString)
		cobra.CheckErr(err)

		viper.Set(cmdString, value)
		fmt.Printf("Set %s to '%s'\n", cmdString, value)
	}
}

func saveIfIntChanged(cmd *cobra.Command, cmdString string) {
	if cmd.LocalFlags().Changed(cmdString) {
		value, err := cmd.LocalFlags().GetInt(cmdString)
		cobra.CheckErr(err)

		viper.Set(cmdString, value)
		fmt.Printf("Set %s to %d\n", cmdString, value)
	}
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure the CLI",
	Long:  `Set the domain, year and session token to talk to AOC`,
	Run: func(cmd *cobra.Command, args []string) {
		saveIfStringChanged(cmd, "domain")
		saveIfIntChanged(cmd, "year")
		saveIfIntChanged(cmd, "day")
		saveIfStringChanged(cmd, "session-token")
		saveIfStringChanged(cmd, "root-dir")
		saveIfStringChanged(cmd, "python-venv")

		day := helpers.GetViperValueEnsureSet[int]("day")
		currentTime := time.Now()

		dayOfMonth := currentTime.Day()

		if day != dayOfMonth {
			fmt.Printf("Day in config '%d' is not today (%d), do you wish to update?(Y/n)\n", day, dayOfMonth)
			var choice string
			fmt.Scanln(&choice)

			if strings.Contains("yY", choice) {
				fmt.Printf("Updating day to %d\n", dayOfMonth)
				day = dayOfMonth
				viper.Set("day", day)
			}
		}

		fmt.Println("== Current Configuration ==")
		fmt.Printf("Domain: %s\n", helpers.GetViperValueEnsureSet[string]("domain"))
		fmt.Printf("Year: %d\n", helpers.GetViperValueEnsureSet[int]("year"))
		fmt.Printf("Day: %d\n", day)
		fmt.Printf("Session token: %s\n", helpers.GetViperValueEnsureSet[string]("session-token"))
		fmt.Printf("Root directory: %s\n", helpers.GetViperValueEnsureSet[string]("root-dir"))
		fmt.Printf("Python virtualenv path: %s\n", helpers.GetViperValueEnsureSet[string]("python-venv"))

		cobra.CheckErr(viper.WriteConfig())
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

	domain := viper.GetString("domain")
	year := viper.GetInt("year")
	day := viper.GetInt("day")
	sessionToken := viper.GetString("session-token")
	rootDir := viper.GetString("root-dir")
	pythonVenv := viper.GetString("python-venv")

	configureCmd.Flags().StringP("domain", "d", domain, "Domain of AOC")
	configureCmd.Flags().IntP("year", "y", year, "Selected year")
	configureCmd.Flags().Int("day", day, "Selected day")
	configureCmd.Flags().StringP("session-token", "t", sessionToken, "Session token copied from AOC")
	configureCmd.Flags().StringP("root-dir", "r", rootDir, "Root directory for the problems")
	configureCmd.Flags().StringP("python-venv", "p", pythonVenv, "Path to the python virtual env")
}
