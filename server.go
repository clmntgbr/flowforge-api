package main

import (
	"fmt"
	"forgeflow-api/config"
	"forgeflow-api/deps"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
	cfg := config.Load()

	db := config.ConnectDatabase(cfg)
	config.AutoMigrate(db)

	app := fiber.New(fiber.Config{
		AppName:       "Flowforge API",
		ServerHeader:  "Flowforge API",
		CaseSensitive: true,
		StrictRouting: true,
		UnescapePath:  true,
	})

	app.Use(helmet.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	deps := deps.New(db, cfg)

	setupHealthChecks(app)
	setupWebhooks(app, deps)
	setupAPIRoutes(app, deps)

	fmt.Println("🚀 Server is running on port", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}

func setupHealthChecks(app *fiber.App) {
	app.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	app.Get(healthcheck.ReadinessEndpoint, healthcheck.New())
	app.Get(healthcheck.StartupEndpoint, healthcheck.New())
}

func setupWebhooks(app *fiber.App, deps *deps.Dependencies) {
	webhooks := app.Group("/webhook")

	webhooks.Post("/clerk", deps.ClerkWebhookMiddleware.Protected(), deps.WebhookClerkHandler.Handle)
}

func setupAPIRoutes(app *fiber.App, deps *deps.Dependencies) {
	api := app.Group("/api")

	api.Use(deps.AuthenticateMiddleware.Protected())
	setupUsersRoutes(api, deps)
	setupProjectsRoutes(api, deps)
}

func setupProjectsRoutes(api fiber.Router, deps *deps.Dependencies) {
	api.Get("/projects", deps.ProjectHandler.GetProjects)
}

func setupUsersRoutes(api fiber.Router, deps *deps.Dependencies) {
	api.Get("/users/me", deps.UserHandler.GetUser)
}
