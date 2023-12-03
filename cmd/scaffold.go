package cmd

import (
	"fmt"

	"github.com/aoccli/client"
	"github.com/aoccli/helpers"
	"github.com/spf13/cobra"
)

var scaffoldCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Scaffold a day",
	Long:  `Scaffold all the files for a specific day`,
	Run: func(cmd *cobra.Command, args []string) {
		year := helpers.GetViperValueEnsureSet[int]("year")
		day := helpers.GetViperValueEnsureSet[int]("day")

		domain := helpers.GetViperValueEnsureSet[string]("domain")
		sessionToken := helpers.GetViperValueEnsureSet[string]("session-token")
		rootDir := helpers.GetViperValueEnsureSet[string]("root-dir")
		aocClient, err := client.NewAOCClient(domain, sessionToken)
		cobra.CheckErr(err)

		dayInput, err := aocClient.GetDayInput(year, day)
		cobra.CheckErr(err)

		scaffolder, err := client.NewFileClient(rootDir, year, day)
		cobra.CheckErr(err)
		cobra.CheckErr(scaffolder.ScaffoldDay(year, day))
		cobra.CheckErr(scaffolder.WriteInput([]byte(dayInput)))
		fmt.Printf("Successfully scaffolded year %d, day %d\n", year, day)
	},
}

func init() {
	rootCmd.AddCommand(scaffoldCmd)
}
