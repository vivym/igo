package main

import (
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
)

// Provisioned by ldflags
// nolint: gochecknoglobals
var (
	version    string
	commitHash string
	buildDate  string
)

var verbose *bool

func main() {
	app := cli.App("igo", "iCloud Client in go")

	app.Spec = "[-v]"

	verbose = app.BoolOpt("v verbose", false, "Verbose debug mode")

	app.Before = func() {
		if *verbose {
			// TODO: enable verbose
			fmt.Println("Verbose mode enabled")
		}
		configure()
	}

	app.Command("login", "iCloud login", loginCmd)

	app.Command("drive", "iCloud Drive", driveCmd)

	app.Command("version", "version", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			fmt.Printf("%s version %s (%s) built on %s\n", friendlyAppName, version, commitHash, buildDate)
		}
	})

	app.Run(os.Args)
}
