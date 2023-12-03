package cmd

import (
	"fmt"
	"strings"

	"github.com/aoccli/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit the solution",
	Long:  "Submit the solution to a problem. Will try to figure out what to submit.",
	Run: func(cmd *cobra.Command, args []string) {
		year := viper.GetInt("year")
		day := viper.GetInt("day")

		domain := viper.GetString("domain")
		sessionToken := viper.GetString("session-token")
		rootDir := viper.GetString("root-dir")
		aocClient, err := client.NewAOCClient(domain, sessionToken)
		cobra.CheckErr(err)

		fileClient, err := client.NewFileClient(rootDir, year, day)
		cobra.CheckErr(err)

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
			cobra.CheckErr(err)
			response, err := aocClient.SubmitProblem(year, day, i, solutionString)
			cobra.CheckErr(err)

			if strings.Contains(response, "You gave an answer too recently") {
				leftIx := strings.Index(response, " left to wait.")
				startIx := -1
				for j := leftIx; j > leftIx-20; j-- {
					if response[j] == '.' {
						startIx = j
						break
					}
				}

				fmt.Printf("Too many requests!! %s\n", response[startIx:leftIx+len(" left to wait.")])
				break
			}

			if strings.Contains(response, "That's not the right answer") {
				fmt.Printf("Solution \n'%s' to problem %d was not correct..\n", strings.TrimRight(solutionString, "\n"), i)
				break
			}

			fmt.Printf("Solution '%s' to problem %d is correct!!\n", solutionString, i)
			err = fileClient.SetProblemSolved(i)
			cobra.CheckErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
}
