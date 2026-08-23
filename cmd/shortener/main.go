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
		Short:         "Shorten, load-test, and list URLs",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(shortenCmd(), loadTestCmd(), listCmd())

	if err := root.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
