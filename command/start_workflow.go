package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewStartWorkflowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "workflow:start",
		Short: "Start a workflow",
		Long:  "Start a workflow",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Starting workflows from API...")
			return nil
		},
	}
}
