package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/tyrantkhan/bb/cmd"
	"github.com/tyrantkhan/bb/internal/cmdutil"
)

func main() {
	if err := cmd.NewRootCommand().Run(context.Background(), os.Args); err != nil {
		if errors.Is(err, cmdutil.ErrShowedUsage) {
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
