package command

import (
	"context"
	mainWire "flowforge-api/cmd/wire"
	"flowforge-api/infrastructure/config"
	"fmt"

	"github.com/spf13/cobra"
)

func NewExecuteWorkflowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "execute-workflow",
		Short: "Execute a workflow",
		Long:  "Execute a workflow",
		RunE: func(cmd *cobra.Command, args []string) error {
			env := config.Load()
			db := config.ConnectDatabase(env)

			container := mainWire.NewContainer(db, env)

			err := container.ExecuteWorkflowUseCase.Execute(context.Background())
			if err != nil {
				return fmt.Errorf("🚨 failed to execute workflow command: %w", err)
			}

			return nil
		},
	}
}
