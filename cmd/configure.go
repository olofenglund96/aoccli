/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"
	"time"

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

// configureCmd represents the configure command
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

		day := viper.GetInt("day")
		currentTime := time.Now()

		// Get the day of the month
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
		fmt.Printf("Domain: %s\n", viper.GetString("domain"))
		fmt.Printf("Year: %d\n", viper.GetInt("year"))
		fmt.Printf("Day: %d\n", day)
		fmt.Printf("Session token: %s\n", viper.GetString("session-token"))
		fmt.Printf("Root directory: %s\n", viper.GetString("root-dir"))
		fmt.Printf("Python virtualenv path: %s\n", viper.GetString("python-venv"))

		cobra.CheckErr(viper.WriteConfig())
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
