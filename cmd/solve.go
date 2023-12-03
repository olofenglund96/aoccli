package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aoccli/helpers"
	"github.com/spf13/cobra"
)

var solveCmd = &cobra.Command{
	Use:   "solve",
	Short: "Solve the specified day problem",
	Long:  "Solve the specified day problem",
	Args: func(cmd *cobra.Command, args []string) error {
		return helpers.ValidateDayArg(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		problem, err := strconv.Atoi(args[0])
		useTest, err := cmd.LocalFlags().GetBool("test")
		cobra.CheckErr(err)

		printCommand, err := cmd.LocalFlags().GetBool("print")
		cobra.CheckErr(err)

		year := helpers.GetViperValueEnsureSet[int]("year")
		day := helpers.GetViperValueEnsureSet[int]("day")
		venv := helpers.GetViperValueEnsureSet[string]("python-venv")
		rootDir := helpers.GetViperValueEnsureSet[string]("root-dir")
		pythonExecutable := filepath.Join(venv, "bin", "python")

		inputFile := "input"

		if useTest {
			inputFile = "test"
		}

		problemPath := filepath.Join(rootDir, strconv.Itoa(year), strconv.Itoa(day), fmt.Sprintf("s%d.py", problem))
		solutionPath := filepath.Join(rootDir, strconv.Itoa(year), strconv.Itoa(day), fmt.Sprintf("%d.sol", problem))

		if printCommand {
			fmt.Printf("Not running command, only printing..\n")
			fmt.Printf("%s %s %s\n", pythonExecutable, problemPath, inputFile)
			return
		}

		command := exec.Command(pythonExecutable, problemPath, inputFile)
		command.Stdout = os.Stdout
		var buffer bytes.Buffer
		command.Stderr = io.MultiWriter(os.Stderr, &buffer)

		fmt.Printf("== Solving problem %d.. ==\n", problem)
		err = command.Run()
		cobra.CheckErr(err)

		stderrStr := buffer.String()
		fmt.Printf("== Solved problem, got output: '%s' ==\n", strings.Replace(stderrStr, "\n", "", 1))
		if !useTest {
			fmt.Printf("Saving to file %s\n", solutionPath)
			err = os.WriteFile(solutionPath, []byte(stderrStr), 0755)
			cobra.CheckErr(err)
		} else {
			fmt.Println("Did not run on real input, not saving solution..")
		}
	},
}

func init() {
	rootCmd.AddCommand(solveCmd)
	solveCmd.Flags().BoolP("test", "t", false, "Use test file as input instead of the real input")
	solveCmd.Flags().BoolP("print", "p", false, "Just print run command instead of running")
}
