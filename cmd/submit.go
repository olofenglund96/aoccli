package cmd

import (
	"fmt"
	"strings"

	"github.com/aoccli/client"
	"github.com/aoccli/helpers"
	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit the solution",
	Long:  "Submit the solution to a problem. Will try to figure out what to submit.",
	Run: func(cmd *cobra.Command, args []string) {
		year := helpers.GetViperValueEnsureSet("year")
		day := helpers.GetViperValueEnsureSet("day")

		sessionToken := helpers.GetViperValueEnsureSet("session-token")
		rootDir := helpers.GetViperValueEnsureSet("root-dir")

		aocClient, err := client.NewAOCClient(sessionToken)
		helpers.HandleErr(err)

		fileClient, err := client.NewFileClient(rootDir, year, day)
		helpers.HandleErr(err)

		for i := 1; i <= 2; i++ {
			if problemSolved, err := fileClient.ProblemSolved(i); problemSolved && err == nil {
				fmt.Printf("Problem %d is already solved and submitted, skipping..\n", i)
				continue
			}

			if exists, err := fileClient.SolutionFileExists(i); !exists && err == nil {
				fmt.Printf("Problem %d not solved yet..\n", i)
				break
			}

			fmt.Printf("Problem %d solved but not submitted, submitting..\n", i)
			solutionString, err := fileClient.ReadSolutionFile(i)
			helpers.HandleErr(err)
			response, err := aocClient.SubmitProblem(year, day, i, solutionString)
			helpers.HandleErr(err)

			fmt.Printf("== Got response: == \n\n%s\n\n==============\n", response)

			if strings.Contains(response, "That's not the right answer") {
				fmt.Printf("Solution \n'%s' to problem %d was not correct..\n", strings.TrimRight(solutionString, "\n"), i)
				break
			}
			fmt.Printf("Solution '%s' to problem %d is correct!!\n", solutionString, i)
			err = fileClient.SetProblemSolved(i)
			helpers.HandleErr(err)
			err = fileClient.CreateSecondSolutionFile()
			helpers.HandleErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
}
