package main

import (
	"context"
	"fmt"
	"os"
	"rs/pkg/server"
)

func main() {
	ctx := context.Background()
	if err := server.Run(ctx, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
