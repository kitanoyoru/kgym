package run

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	return &cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
}
