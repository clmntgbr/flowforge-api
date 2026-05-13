package command

import (
	"flowforge-api/domain/entity"
	"flowforge-api/infrastructure/config"
	"fmt"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func NewMigrateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Migrate the database",
		Long:  "Migrate the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			env := config.Load()
			db := config.ConnectDatabase(env)

			err := db.Transaction(func(tx *gorm.DB) error {
				return tx.AutoMigrate(
					&entity.User{},
					&entity.Organization{},
					&entity.Workflow{},
					&entity.Step{},
					&entity.Endpoint{},
					&entity.Connexion{},
					&entity.WorkflowRun{},
					&entity.StepRun{},
					&entity.Insight{},
				)
			})

			if err != nil {
				return fmt.Errorf(
					"🚨 failed to migrate database: %w",
					err,
				)
			}

			fmt.Println("🎉 Database migrations completed successfully")
			return nil
		},
	}
}
