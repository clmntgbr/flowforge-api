package main

import (
	"flowforge-api/cmd/cli/command"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cli",
		Short: "Flowforge CLI - cmd commands",
		Long:  "Flowforge CLI provides commands for cmd tasks",
	}

	rootCmd.AddCommand(
		command.NewMigrateCommand(),
		command.NewExecuteWorkflowCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
