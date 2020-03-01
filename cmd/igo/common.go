package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

func promptString(name string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(name + ": ")
	str, _ := reader.ReadString('\n')

	return strings.TrimSpace(str)
}

func promptPassword(name string) string {
	fmt.Print(name + ": ")
	bytePassword, _ := terminal.ReadPassword(0)
	fmt.Print("\n")
	password := string(bytePassword)
	return strings.TrimSpace(password)
}

func verifyEmailFormat(str string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(str)
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func pathExists(p string) (bool, error) {
	_, err := os.Stat(p)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, nil
}
