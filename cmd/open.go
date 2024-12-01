package cmd

import (
	"os/exec"
	"strings"

	"github.com/aoccli/helpers"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open the url for the specified day in a web browser",
	Long:  "Opens the url for the specified day in a web browser using xdg-open",
	Run: func(cmd *cobra.Command, args []string) {
		year := helpers.GetViperValueEnsureSet("year")
		day := helpers.GetViperValueEnsureSet("day")

		url := helpers.GetDayUrl("https://adventofcode.com", year, day)
		helpers.HandleErr(openBrowser(url))
	},
}

func openBrowser(url string) error {
	unameOutput, err := exec.Command("uname", "-s").Output()
	if err != nil {
		return err
	}

	openExecutable := "xgd-open"
	if strings.Contains(string(unameOutput), "Darwin") {
		openExecutable = "open"
	}

	browserCmd := exec.Command(openExecutable, url)
	return browserCmd.Run()
}

func init() {
	rootCmd.AddCommand(openCmd)
}
