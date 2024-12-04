package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/aoccli/helpers"
	"github.com/spf13/cobra"
)

var aoccliFolderPath string

var installLibsCmd = &cobra.Command{
	Use:   "install-libs",
	Short: "Install python helper libs",
	Long:  `Installs a python lib called aoclib that contains useful functions for solving aoc problems`,
	Run: func(cmd *cobra.Command, args []string) {
		pythonExecutable := helpers.GetViperValueEnsureSet("python-exec")

		pathToLibs := path.Join(aoccliFolderPath, "libs", "python")

		command := exec.Command(pythonExecutable, "-m", "pip", "install", "-e", pathToLibs)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		fmt.Print("== Installing python libs.. ==\n")
		err := command.Run()
		helpers.HandleErr(err)

		fmt.Println("Successfully installed python libs called 'aoclib'. Use by adding 'import aoclib' to .py files")
	},
}

func init() {
	installLibsCmd.Flags().StringVarP(&aoccliFolderPath, "aoccli-directory", "d", "", "Absolute path to the source directory of aoccli")
	installLibsCmd.MarkFlagRequired("aoccli-directory")
	rootCmd.AddCommand(installLibsCmd)
}
