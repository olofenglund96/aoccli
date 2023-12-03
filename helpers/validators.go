package helpers

import (
	"strconv"

	"github.com/spf13/cobra"
)

func ValidateDayArg(cmd *cobra.Command, args []string) error {
	if err := cobra.ExactArgs(1)(cmd, args); err != nil {
		return err
	}

	if _, err := strconv.Atoi(args[0]); err != nil {
		return err
	}

	return nil
}
