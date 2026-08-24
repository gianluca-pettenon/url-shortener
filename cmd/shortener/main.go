package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	defer stop()

	root := &cobra.Command{
		Use:           "shortener",
		Short:         "Shorten, redirect, load-test, and list URLs",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(shortenCmd(), loadTestCmd(), listCmd(), serveCmd())

	if err := root.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
