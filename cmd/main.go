package main

import (
	"os"

	"github.com/jenkins-x-labs/cli-apps/cmd/app"
)

// Entrypoint for the command
func main() {
	err := app.Run(nil)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
