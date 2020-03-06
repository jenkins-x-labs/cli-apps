// +build !windows

package app

import (
	"github.com/jenkins-x-labs/cli-apps/cmd/apps"
)

// Run runs the command, if args are not nil they will be set on the command
func Run(args []string) error {
	cmd := apps.NewCmdHelloWorld()
	if args != nil {
		args = args[1:]
		cmd.SetArgs(args)
	}
	return cmd.Execute()
}
