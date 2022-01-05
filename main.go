package main

import (
	"context"
	"errors"
	"github.com/wiz-sec/cibus/internal"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: cibus [username] [password]")
	}

	err := internal.AddAllFriends(os.Args[1], os.Args[2])
	if err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			log.Fatalf("%s", err.Error())
		}
	}
}
