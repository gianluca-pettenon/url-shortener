package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gianluca-pettenon/url-shortener/internal/urls"
	"github.com/spf13/cobra"
)

func shortenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "shorten",
		Short: "Prompt for a URL and print the short code",
		RunE: func(cmd *cobra.Command, args []string) error {
			in := bufio.NewReader(cmd.InOrStdin())

			raw, err := prompt(cmd, in, "URL: ")

			if err != nil {
				return err
			}

			return withService(cmd, func(svc *urls.Service) error {
				code, err := svc.Create(cmd.Context(), raw)

				if err != nil {
					return err
				}

				fmt.Fprintf(cmd.OutOrStdout(), "code: %s\n", code)

				return nil
			})
		},
	}
}

func loadTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "load-test",
		Short: "Prompt for a URL and how many times to shorten it",
		RunE: func(cmd *cobra.Command, args []string) error {
			in := bufio.NewReader(cmd.InOrStdin())

			raw, err := prompt(cmd, in, "URL: ")

			if err != nil {
				return err
			}

			n, err := promptTimes(cmd, in)

			if err != nil {
				return err
			}

			return withService(cmd, func(svc *urls.Service) error {
				start := time.Now()
				first, last, err := svc.CreateMany(cmd.Context(), raw, n)

				if err != nil {
					return err
				}

				elapsed := time.Since(start)
				secs := elapsed.Seconds()

				if secs <= 0 {
					secs = 0.001
				}

				firstCode, err := svc.Code(first)

				if err != nil {
					return err
				}

				lastCode, err := svc.Code(last)

				if err != nil {
					return err
				}

				fmt.Fprintf(cmd.OutOrStdout(), "First: %s\nLast: %s\n", firstCode, lastCode)
				fmt.Fprintf(cmd.ErrOrStderr(), "%d ok in %s (%.0f/s)\n",
					n,
					elapsed.Round(time.Millisecond),
					float64(n)/secs,
				)

				return nil
			})
		},
	}
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List shortened URLs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return withService(cmd, func(svc *urls.Service) error {
				items, err := svc.List(cmd.Context())

				if err != nil {
					return err
				}

				if len(items) == 0 {
					fmt.Fprintln(cmd.OutOrStdout(), "No shortened URLs")
					return nil
				}

				for _, u := range items {
					code, err := svc.Code(u.ID)

					if err != nil {
						return err
					}

					fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\n",
						code,
						u.CreatedAt.Format("2006-01-02 15:04"),
						u.OriginalURL,
					)
				}

				return nil
			})
		},
	}
}

func withService(cmd *cobra.Command, fn func(*urls.Service) error) error {
	svc, cleanup, err := open(cmd.Context())

	if err != nil {
		return err
	}

	defer cleanup()

	return fn(svc)
}

func prompt(cmd *cobra.Command, in *bufio.Reader, label string) (string, error) {
	fmt.Fprint(cmd.OutOrStdout(), label)
	raw, err := in.ReadString('\n')

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(raw), nil
}

func promptTimes(cmd *cobra.Command, in *bufio.Reader) (int, error) {
	raw, err := prompt(cmd, in, "Times: ")

	if err != nil {
		return 0, err
	}

	n, err := strconv.Atoi(raw)

	if err != nil || n < 1 {
		return 0, fmt.Errorf("times must be an integer >= 1")
	}

	return n, nil
}
