package main

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
			cfg := config.Load()
			db := config.ConnectDatabase(cfg)

			err := db.Transaction(func(tx *gorm.DB) error {
				return tx.AutoMigrate(
					&entity.User{},
					&entity.Organization{},
					&entity.Workflow{},
					&entity.Step{},
					&entity.Endpoint{},
					&entity.Connexion{},
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
