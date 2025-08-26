package main

import (
	"fmt"
	"os"

	"challenge/cli"
)

func main() {
	cli := cli.NewCLI()

	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
