package cmd

import (
	"log"

	"github.com/kitanoyoru/kgym/internal/apps/user/cmd/run"
	"github.com/kitanoyoru/kgym/internal/apps/user/cmd/shutdown"
	pkgLogger "github.com/kitanoyoru/kgym/pkg/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "kgym-user-service",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger, err := pkgLogger.New(pkgLogger.WithDev())
		if err != nil {
			log.Fatal(err)
		}

		/*
			TODO: What I want to see

			if err := dependecy.InjectToContext(cmd.Context(), []dependency.Dependecy{
				logger.AsDependency(),
			}); err != nil {
				log.Fatal(err)
			}
		*/

		cmd.SetContext(pkgLogger.Inject(cmd.Context(), logger.Logger))
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger, err := pkgLogger.FromContext(cmd.Context())
		if err != nil {
			log.Fatal(err)
		}

		/*
			TODO: What I want to see

			if err := dependenncy.CloseAll(cmd.Context()); err != nil {
				log.Fatal(err)
			}

		*/

		if err := logger.Sync(); err != nil {
			log.Fatal(err)
		}
	},
}

func main() {
	rootCmd.AddCommand(run.Command())
	rootCmd.AddCommand(shutdown.Command())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
