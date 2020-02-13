/*
	Credential Server
	version 0.9
	author: Adrian Pareja Abarca
	email: adriancc5.5@gmail.com
*/

package main

import "os"

var (
	blockingStart = true
)

// The credential provider server main
func main() {
	if err := RunMain(os.Args); err != nil {
		os.Exit(1)
	}
}

// RunMain is the server main
func RunMain(args []string) error {
	// Save the os.Args
	saveOsArgs := os.Args
	os.Args = args

	cmdName := ""
	if len(args) > 1 {
		cmdName = args[1]
	}
	scmd := NewCommand(cmdName, blockingStart)

	// Execute the command
	err := scmd.Execute()

	// Restore original os.Args
	os.Args = saveOsArgs

	return err
}