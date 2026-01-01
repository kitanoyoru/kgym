package run

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/kitanoyoru/kgym/internal/apps/file/internal"
	"github.com/kitanoyoru/kgym/internal/apps/file/pkg/env"
	"github.com/spf13/cobra"
)

const (
	Version = "1.0.0"
	Short   = "Run the File application"
	Long    = "Run the File application"
	Use     = "run"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:     Use,
		Version: Version,
		Short:   Short,
		Long:    Long,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			var cfg internal.Config
			if err := env.ParseAndValidate(ctx, &cfg); err != nil {
				return err
			}

			ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
			defer cancel()

			app, err := internal.New(ctx, cfg)
			if err != nil {
				return err
			}

			errChan := make(chan error, 1)
			go func() {
				errChan <- app.Run(ctx)
			}()

			select {
			case <-ctx.Done():
			case err := <-errChan:
				if err != nil {
					return err
				}
			}

			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
			defer shutdownCancel()

			errChan = make(chan error, 1)

			go func() {
				errChan <- app.Shutdown(shutdownCtx)
			}()

			select {
			case <-shutdownCtx.Done():
				return shutdownCtx.Err()
			case err := <-errChan:
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}
