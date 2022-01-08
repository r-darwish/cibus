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

func GetCredentials() (username, password string, err error) {
	if len(os.Args) != 3 {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter Username: ")
		username, err := reader.ReadString('\n')
		if err != nil {
			return "", "", err
		}

		fmt.Print("Enter Password: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", "", err
		}

		password := string(bytePassword)
		return strings.TrimSpace(username), password, nil

	} else {
		return os.Args[1], os.Args[2], nil
	}
}
