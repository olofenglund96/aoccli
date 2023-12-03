/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/aoccli/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// scaffoldCmd represents the scaffold command
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scaffoldCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scaffoldCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
