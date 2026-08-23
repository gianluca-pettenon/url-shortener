package main

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func shortenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "shorten",
		Short: "Prompt for a URL and print the short code",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprint(cmd.OutOrStdout(), "URL: ")
			raw, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')

			if err != nil {
				return err
			}

			svc, cleanup, err := open(cmd.Context())

			if err != nil {
				return err
			}

			defer cleanup()

			code, err := svc.Create(cmd.Context(), strings.TrimSpace(raw))

			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "code: %s\n", code)

			return nil
		},
	}
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List shortened URLs",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, cleanup, err := open(cmd.Context())

			if err != nil {
				return err
			}

			defer cleanup()

			items, err := svc.List(cmd.Context())

			if err != nil {
				return err
			}

			if len(items) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no shortened URLs")
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
		},
	}
}
