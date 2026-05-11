package config

import (
	"flowforge-api/domain/entity"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase(cfg *Config) *gorm.DB {
	logLevel := logger.Warn
	if cfg.Environment == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get database instance: %v", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	log.Println("database connection pool configured")

	return db
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&entity.User{},
		&entity.Organization{},
		&entity.Endpoint{},
		&entity.Workflow{},
		&entity.Step{},
		&entity.Connexion{},
	)

	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("database migrations completed")
}
