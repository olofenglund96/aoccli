/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os/exec"

	"github.com/aoccli/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open the url for the specified day in a web browser",
	Long:  "Opens the url for the specified day in a web browser using xdg-open",
	Run: func(cmd *cobra.Command, args []string) {
		year := viper.GetInt("year")
		day := viper.GetInt("day")
		domain := viper.GetString("domain")

		url := helpers.GetDayUrl(domain, year, day)
		browserCmd := exec.Command("xdg-open", url)
		cobra.CheckErr(browserCmd.Run())
	},
}

func init() {
	rootCmd.AddCommand(openCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// openCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// openCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
