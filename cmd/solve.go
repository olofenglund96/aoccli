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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var solveCmd = &cobra.Command{
	Use:   "solve",
	Short: "Solve the specified day problem",
	Long:  "Solve the specified day problem",
	Run: func(cmd *cobra.Command, args []string) {
		problems := []int{1, 2}
		if err := cobra.ExactArgs(1)(cmd, args); err == nil {
			problem, err := strconv.Atoi(args[0])
			problems = []int{problem}
			cobra.CheckErr(err)
		}

		useTest, err := cmd.LocalFlags().GetBool("test")
		cobra.CheckErr(err)

		printCommand, err := cmd.LocalFlags().GetBool("print")
		cobra.CheckErr(err)

		year := viper.GetInt("year")
		day := viper.GetInt("day")
		venv := viper.GetString("python-venv")
		rootDir := viper.GetString("root-dir")
		pythonExecutable := filepath.Join(venv, "bin", "python")

		inputFile := "input"

		if useTest {
			inputFile = "test"
		}

		for _, problem := range problems {
			problemPath := filepath.Join(rootDir, strconv.Itoa(year), strconv.Itoa(day), fmt.Sprintf("s%d.py", problem))
			solutionPath := filepath.Join(rootDir, strconv.Itoa(year), strconv.Itoa(day), fmt.Sprintf("%d.sol", problem))

			if printCommand {
				fmt.Printf("Not running command, only printing..\n")
				fmt.Printf("%s %s %s\n", pythonExecutable, problemPath, inputFile)
				continue
			}

			command := exec.Command(pythonExecutable, problemPath, inputFile)
			command.Stdout = os.Stdout
			var buffer bytes.Buffer
			command.Stderr = io.MultiWriter(os.Stderr, &buffer)

			fmt.Printf("== Solving problem %d.. ==\n", problem)
			err := command.Run()
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
		}
	},
}

func init() {
	rootCmd.AddCommand(solveCmd)
	solveCmd.Flags().BoolP("test", "t", false, "Use test file as input instead of the real input")
	solveCmd.Flags().BoolP("print", "p", false, "Just print run command instead of running")
}
