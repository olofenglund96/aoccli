/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/aoccli/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download input for the specific day",
	Long:  `Downloads the input for a day and returns it`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		if _, err := strconv.Atoi(args[0]); err != nil {
			return err
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		year := viper.GetInt("year")
		day, _ := strconv.Atoi(args[0])

		domain := viper.GetString("domain")
		sessionToken := viper.GetString("session-token")
		aocClient, err := client.NewAOCClient(domain, sessionToken)
		cobra.CheckErr(err)

		dayInput, err := aocClient.GetDayInput(year, day)
		cobra.CheckErr(err)

		fmt.Printf("download called, got this: \n %s \n", dayInput)
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
