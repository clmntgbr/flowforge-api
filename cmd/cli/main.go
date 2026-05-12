package main

import (
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
		NewMigrateCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
