package cmd

import (
	"fmt"
	"strings"

	"github.com/aoccli/client"
	"github.com/aoccli/helpers"
	"github.com/spf13/cobra"
)

type submitter struct {
	aocClient client.AOCClient
	fileClient client.FileClient
	year string
	day string
}

func (s submitter) submitSubProblem(problemIndex int) (bool, error) {
	problemSolved, err := s.fileClient.ProblemSolved(problemIndex)
	if err != nil {
		return false, fmt.Errorf("failed to check if problem was already solved: %s", err)
	}

	if problemSolved {
		fmt.Printf("Problem %d is already solved and submitted, skipping..\n", problemIndex)
		return true, nil
	}

	exists, err := s.fileClient.SolutionFileExists(problemIndex)
	if err != nil {
		return false, fmt.Errorf("failed to check if solution file exists: %s", err)
	}

	if !exists{
		fmt.Printf("Problem %d not solved yet..\n", problemIndex)
		return false, fmt.Errorf("solve problem %d by running 'aoc solve %d' before submitting..", problemIndex, problemIndex)
	}

	fmt.Printf("Problem %d solved but not submitted, submitting..\n", problemIndex)
	solutionString, err := s.fileClient.ReadSolutionFile(problemIndex)
	if err != nil {
		return false, fmt.Errorf("failed to read solution file: %s", err)
	}
	response, err := s.aocClient.SubmitProblem(s.year, s.day, problemIndex, solutionString)
	if err != nil {
		return false, fmt.Errorf("failed to submit solution: %s", err)
	}

	fmt.Printf("== Got response: == \n\n%s\n\n==============\n", response)

	if strings.Contains(response, "That's not the right answer") {
		fmt.Printf("Solution \n'%s' to problem %d was not correct..\n", strings.TrimRight(solutionString, "\n"), problemIndex)
		return false, nil
	}

	fmt.Printf("Solution '%s' to problem %d is correct!!\n", solutionString, problemIndex)
	err = s.fileClient.SetProblemSolved(problemIndex)
	if err != nil {
		return false, fmt.Errorf("failed to set problem as solved locally: %s", err)
	}

	return true, nil
}


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

		submitter := submitter{
			aocClient: aocClient,
			fileClient: fileClient,
			year: year,
			day: day,
		}

		solved, err := submitter.submitSubProblem(1)
		if err != nil {
			helpers.HandleErr(err)
		}

		if !solved {
			return
		}

		fileAlreadyExisted, err := fileClient.CreateSecondSolutionFile()
		if err != nil {
			helpers.HandleErr(fmt.Errorf("failed to set problem as solved locally: %s", err))
		}

		if !fileAlreadyExisted {
			return
		}

		solved, err = submitter.submitSubProblem(2)
		if err != nil {
			helpers.HandleErr(err)
		}

		if !solved {
			return
		}

		fmt.Printf("All problems solved for %s/%s! 🎉\n", year, day)
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
}
