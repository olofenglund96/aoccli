package cmd

import (
	"os/exec"

	"github.com/aoccli/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
}
