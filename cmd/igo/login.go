package main

import (
	"fmt"
	"os"
	"path"

	"emperror.dev/emperror"
	"emperror.dev/errors"
	cli "github.com/jawher/mow.cli"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/vivym/igo"
)

func loginCmd(cmd *cli.Cmd) {
	var (
		username = cmd.StringOpt("u username", "", "Username (email address)")
		password = cmd.StringOpt("p password", "", "Password")
	)

	cmd.Action = func() {
		if *username == "" {
			*username = promptString("Username")
		}
		if !verifyEmailFormat(*username) {
			fmt.Println("invalid username: ", *username)
			os.Exit(-1)
		}

		if *password == "" {
			*password = promptPassword("Password")
		}
		if *password == "" {
			fmt.Println("password is required.")
			os.Exit(-1)
		}

		client := igo.New()
		defer client.Close()
		login(client, *username, *password)
	}
}

func getSessionPath() string {
	home, _ := homedir.Dir()
	return path.Join(home, ".igo", "igo.session")
}

func loadSession(client *igo.Client) error {
	sessionPath := getSessionPath()
	fp, err := os.OpenFile(sessionPath, os.O_CREATE|os.O_RDONLY, 0666)
	defer fp.Close()
	if err != nil {
		return nil
	}

	err = client.LoadSession(fp)
	return err
}

func saveSession(client *igo.Client) {
	sessionPath := getSessionPath()
	fp, err := os.OpenFile(sessionPath, os.O_CREATE|os.O_WRONLY, 0666)
	defer fp.Close()
	emperror.Panic(errors.Wrap(err, "failed to open session file: "+sessionPath))

	err = client.SaveSession(fp)
	emperror.Panic(errors.Wrap(err, "failed to save session"))

	if *verbose {
		fmt.Println("Session is saved to /tmp/igo.session.")
	}
}

func login(client *igo.Client, username, password string) {
	if err := client.Login(username, password); err != nil {
		fmt.Println("Login failed:", err)
		return
	}
	if client.TwoFactorAuthenticationIsRequired() {
		securityCode := promptString("SecurityCode")
		if err := client.SetSecurityCode(securityCode); err != nil {
			fmt.Println("err:", err)
			return
		}
	}
	fmt.Println("Login successfualy.")

	saveSession(client)
}
