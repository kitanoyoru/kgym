package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/kitanoyoru/kgym/internal/apps/user/cmd/run"
	"github.com/kitanoyoru/kgym/internal/apps/user/cmd/shutdown"
)

var rootCmd = &cobra.Command{
	Use: "kgym-user-service",
}

func main() {
	rootCmd.AddCommand(run.Command())
	rootCmd.AddCommand(shutdown.Command())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
