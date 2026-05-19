package main

import (
	cliCommand "flowforge-api/cmd/cli/command"
	cliWire "flowforge-api/cmd/cli/wire"
	"flowforge-api/infrastructure/config"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	env := config.Load()
	db := config.ConnectDatabase(env)

	container := cliWire.NewContainer(db, env)

	rootCmd := &cobra.Command{
		Use:   "cli",
		Short: "Flowforge CLI - cmd commands",
		Long:  "Flowforge CLI provides commands for cmd tasks",
	}

	rootCmd.AddCommand(
		cliCommand.NewMigrateCommand(),
		cliCommand.NewExecuteWorkflowCommand(container),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
