package command

import (
	"context"
	"flowforge-api/cmd/cli/wire"
	"fmt"

	"github.com/spf13/cobra"
)

func NewExecuteWorkflowCommand(container *wire.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "execute-workflow",
		Short: "Execute a workflow",
		Long:  "Execute a workflow",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := container.ExecuteWorkflowUseCase.Execute(context.Background())
			if err != nil {
				return fmt.Errorf("🚨 failed to execute workflow command: %w", err)
			}

			return nil
		},
	}
}
