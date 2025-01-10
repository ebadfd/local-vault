package main

import (
	"os"

	"github.com/ebadfd/local-vault/cmd"
)

var version = "0.0.1"

func main() {
	if err := cmd.Execute(version); err != nil {
		os.Exit(1)
	}
}
