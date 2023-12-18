package helpers

import (
	"fmt"
	"os"
	"runtime/debug"
)

func HandleErr(err error) {
	if err != nil {
		fmt.Printf("Error: %s", err)
		debug.PrintStack()
		os.Exit(1)
	}
}
