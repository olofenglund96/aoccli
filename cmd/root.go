package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "aoccli",
		Short: "An Advent of Code CLI",
		Long: `A CLI tool to help you with Advent of Code. Can, for example,
open the daily problem web page, scaffold a day, download input, submit etc..
Uses the session cookie copied from the browser for authentication.

A good place to start would be to run 'configure' to set up configuration.`,
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("%+v\n", errors.WithStack(err))
		os.Exit(1)
	}
}

func init() {
	initConfig()
}

func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	configPath := filepath.Join(home, "projects", "advent-of-code")

	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".aoc")
	viper.SafeWriteConfig()

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("Config file not found: %s..\n", err)
			cobra.CheckErr(err)
		} else {
			fmt.Printf("Unknown error occurred: %s\n", err)
			cobra.CheckErr(err)
		}
	}
}
