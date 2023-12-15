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
		year := helpers.GetViperValueEnsureSet("year")
		day := helpers.GetViperValueEnsureSet("day")

		sessionToken := helpers.GetViperValueEnsureSet("session-token")
		rootDir := helpers.GetViperValueEnsureSet("root-dir")
		aocClient, err := client.NewAOCClient(sessionToken)

		helpers.HandleErr(err)

		dayInput, err := aocClient.GetDayInput(year, day)
		helpers.HandleErr(err)
		dayTestInput, err := aocClient.GetDayTestInput(year, day)
		helpers.HandleErr(err)

		scaffolder, err := client.NewFileClient(rootDir, year, day)
		helpers.HandleErr(err)
		helpers.HandleErr(scaffolder.ScaffoldDay(year, day))
		helpers.HandleErr(scaffolder.WriteInput([]byte(dayInput)))
		helpers.HandleErr(scaffolder.WriteTestInput([]byte(dayTestInput)))

		fmt.Printf("Successfully scaffolded year %s, day %s\n", year, day)
	},
}

func init() {
	rootCmd.AddCommand(scaffoldCmd)
}
