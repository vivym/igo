package main

import (
	cli "github.com/jawher/mow.cli"
	"github.com/vivym/igo"
)

// Cmd defines commands to access iCloud Drive
func driveCmd(cmd *cli.Cmd) {
	client := igo.New()
	loadSession(client)

	cmd.Command("ls", "list directory", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			client.EnableDrive().Test("FILE::com.apple.CloudDocs::1D4D6C14-EA1A-4E45-8B5E-47C16A02ADF5")
		}
	})
}
