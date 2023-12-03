package cmd

import (
	"fmt"

	"github.com/aoccli/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var scaffoldCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Scaffold a day",
	Long:  `Scaffold all the files for a specific day`,
	Run: func(cmd *cobra.Command, args []string) {
		year := viper.GetInt("year")
		day := viper.GetInt("day")

		domain := viper.GetString("domain")
		sessionToken := viper.GetString("session-token")
		rootDir := viper.GetString("root-dir")
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
