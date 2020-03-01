package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
)

// Cmd defines commands to access iCloud Drive
func driveCmd(cmd *cli.Cmd) {
	cmd.Command("ls", "list directory", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			fmt.Println("ls test")
		}
	})
}
