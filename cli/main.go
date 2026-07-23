package main

import (
	"fmt"
	"os"

	"ticket-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
	_ = fmt.Sprintf
}
