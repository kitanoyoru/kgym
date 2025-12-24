package run

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kitanoyoru/kgym/internal/apps/user/internal"
	"github.com/kitanoyoru/kgym/internal/apps/user/pkg/env"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			var cfg internal.Config
			if err := env.ParseAndValidate(ctx, &cfg); err != nil {
				return err
			}

			ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
			defer cancel()

			app, err := internal.New(cfg)
			if err != nil {
				return err
			}

			errChan := make(chan error, 1)
			go func() {
				if err := app.Run(ctx); err != nil {
					errChan <- err
				}
			}()

			select {
			case <-ctx.Done():
			case err := <-errChan:
				if err != nil {
					return err
				}
			}

			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer shutdownCancel()
			if err := app.Shutdown(shutdownCtx); err != nil {
				return err
			}

			return ctx.Err()
		},
	}
}
