package main

import (
	"fmt"
	"forgeflow-api/command"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cli",
		Short: "ForgeFlow CLI - internal commands",
		Long:  "ForgeFlow CLI provides commands for internal tasks",
	}

	rootCmd.AddCommand(
		command.NewStartWorkflowCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
