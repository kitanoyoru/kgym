package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/kitanoyoru/kgym/internal/apps/sso/cmd/run"
	pkgLogger "github.com/kitanoyoru/kgym/pkg/logger"
)

var rootCmd = &cobra.Command{
	Use: "kgym-sso-service",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger, err := pkgLogger.New(pkgLogger.WithDev())
		if err != nil {
			log.Fatal(err)
		}

		cmd.SetContext(pkgLogger.Inject(cmd.Context(), logger))
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger, err := pkgLogger.FromContext(cmd.Context())
		if err != nil {
			log.Fatal(err)
		}

		if err := logger.Sync(); err != nil {
			log.Fatal(err)
		}
	},
}

func main() {
	rootCmd.AddCommand(run.Command())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
