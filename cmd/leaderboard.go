package cmd

import (
	"fmt"

	"github.com/aoccli/helpers"

	"github.com/aoccli/client"
	"github.com/spf13/cobra"
)

var leaderboardCmd = &cobra.Command{
	Use:   "leaderboard",
	Short: "Show current leaderboard standings",
	Long:  "Show current leaderboard standings",
	Run: func(cmd *cobra.Command, args []string) {
		sessionToken := helpers.GetViperValueEnsureSet("session-token")
		leaderboardId := helpers.GetViperValueEnsureSet("leaderboard")
		aocClient, err := client.NewAOCClient(sessionToken)
		helpers.HandleErr(err)

		leaderboard, err := aocClient.GetLeaderboard(leaderboardId)
		helpers.HandleErr(err)

		fmt.Println("Current leaderboard standings:")
		fmt.Print(leaderboard)
	},
}

func init() {
	rootCmd.AddCommand(leaderboardCmd)
}
