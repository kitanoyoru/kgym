package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/kitanoyoru/kgym/internal/apps/sso/cmd/run"
)

var rootCmd = &cobra.Command{
	Use: "kgym-sso-service",
}

func main() {
	rootCmd.AddCommand(run.Command())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
