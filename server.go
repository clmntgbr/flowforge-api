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
	setupOrganizationsRoutes(api, deps)
	setupEndpointsRoutes(api, deps)
}

func setupOrganizationsRoutes(api fiber.Router, deps *deps.Dependencies) {
	api.Get("/organizations", deps.OrganizationHandler.GetOrganizations)
	api.Get("/organizations/:id", deps.OrganizationHandler.GetOrganizationByID)
	api.Post("/organizations", deps.OrganizationHandler.CreateOrganization)
	api.Put("/organizations/:id", deps.OrganizationHandler.UpdateOrganization)
	api.Put("/organizations/:id/activate", deps.OrganizationHandler.ActivateOrganization)
}

func setupEndpointsRoutes(api fiber.Router, deps *deps.Dependencies) {
	api.Get("/endpoints", deps.EndpointHandler.GetEndpoints)
	// api.Get("/endpoints/:id", deps.EndpointHandler.GetEndpointByID)
	api.Post("/endpoints", deps.EndpointHandler.CreateEndpoint)
	api.Put("/endpoints/:id", deps.EndpointHandler.UpdateEndpoint)
}

func setupUsersRoutes(api fiber.Router, deps *deps.Dependencies) {
	api.Get("/users/me", deps.UserHandler.GetUser)
}
