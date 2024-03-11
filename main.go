package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/wiz-sec/cibus/internal"
	"golang.org/x/term"
	"log"
	"os"
	"strings"
	"syscall"
)

func main() {
	username, password, err := GetCredentials()
	if err != nil {
		log.Fatalf("Unable to retrieve credentials: %s", err)
	}

	err = internal.AddAllFriends(username, password)
	if err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			log.Fatalf("%s", err.Error())
		}
	}
}

func GetCredentials() (string, string, error) {
	username := os.Getenv("CIBUS_USERNAME")
	password := os.Getenv("CIBUS_PASSWORD")
	var err error

	reader := bufio.NewReader(os.Stdin)

	if username == "" {
		fmt.Print("Enter Username: ")
		username, err = reader.ReadString('\n')
		if err != nil {
			return "", "", err
		}
		username = strings.TrimSpace(username)
	}

	if password == "" {
		fmt.Print("Enter Password: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", "", err
		}
		password = string(bytePassword)
	}

	return username, password, nil
}
