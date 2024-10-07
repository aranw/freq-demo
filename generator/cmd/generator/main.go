package main

import (
	"context"
	"fmt"
	"os"
)

var (
	Name string
)

func main() {
	if err := Run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
