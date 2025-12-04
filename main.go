package main

import (
	"os"

	"github.com/Myrtilli/link-shortener-svc/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
