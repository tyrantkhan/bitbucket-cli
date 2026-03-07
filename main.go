package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tyrantkhan/bb/cmd"
)

func main() {
	if err := cmd.NewRootCommand().Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
