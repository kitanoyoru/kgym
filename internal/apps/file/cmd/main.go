package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/kitanoyoru/kgym/internal/apps/file/cmd/run"
	"github.com/kitanoyoru/kgym/pkg/tracing"
)

var tracingFlushFunc tracing.FlushFunc

var rootCmd = &cobra.Command{
	Use: "kgym-file-service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		cfg, err := tracing.ConfigFromEnv(ctx)
		if err != nil {
			return err
		}

		flush, err := tracing.Init(ctx, cfg)
		if err != nil {
			return err
		}

		tracingFlushFunc = flush

		return err
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		if tracingFlushFunc != nil {
			if err := tracingFlushFunc(ctx); err != nil {
				return err
			}
		}

		return nil
	},
}

func main() {
	rootCmd.AddCommand(run.Command())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
