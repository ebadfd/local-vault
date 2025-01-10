package main

import (
	"fmt"
	"os"

	"github.com/ebadfd/local-vault/cmd"
)

var version = "0.0.1"

func main() {
	if err := cmd.Execute(version); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
