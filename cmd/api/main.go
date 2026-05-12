package main

import (
	"flowforge-api/infrastructure/config"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
	env := config.Load()
	db := config.ConnectDatabase(env)

	app := fiber.New(fiber.Config{
		AppName:       "Flowforge API",
		ServerHeader:  "Flowforge API",
		CaseSensitive: true,
		StrictRouting: true,
		UnescapePath:  true,
	})

	app.Use(helmet.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     env.CORSAllowedOrigins,
		AllowMethods:     env.CORSAllowMethods,
		AllowHeaders:     env.CORSAllowHeaders,
		AllowCredentials: env.CORSAllowCredentials,
		MaxAge:           env.CORSMaxAge,
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        env.RateLimitMax,
		Expiration: 1 * time.Minute,
		LimitReached: func(c fiber.Ctx) error {
			log.Println("rate limit exceeded: ", c.IP(), c.Path(), c.Method())
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "too many requests, please try again later",
			})
		},
	}))

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Use(func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		c.Append("Server-Timing", "app;dur="+duration.String())
		return err
	})

	deps := NewContainer(db, env)

	setupHealthChecks(app)
	setupWebhooks(app, deps)
	setupAPIRoutes(app, deps)

	log.Println("🚀 Server is running on port", env.Port)
	log.Fatal(app.Listen(":" + env.Port))
}

func setupWebhooks(app *fiber.App, container *Container) {
	webhooks := app.Group("/webhook")

	webhooks.Post("/clerk", container.ClerkMiddleware.Protected(), container.ClerkHandler.Execute)
}

func setupHealthChecks(app *fiber.App) {
	app.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	app.Get(healthcheck.ReadinessEndpoint, healthcheck.New())
	app.Get(healthcheck.StartupEndpoint, healthcheck.New())
}

func setupAPIRoutes(app *fiber.App, container *Container) {
	api := app.Group("/api")

	api.Use(container.AuthenticateMiddleware.Protected())
	setupUsersRoutes(api, container)
	setupOrganizationsRoutes(api, container)
}

func setupUsersRoutes(api fiber.Router, container *Container) {
	api.Get("/users/me", container.UserHandler.GetUser)
}

func setupOrganizationsRoutes(api fiber.Router, container *Container) {
	api.Get("/organizations", container.OrganizationHandler.GetOrganizations)
	api.Post("/organizations", container.OrganizationHandler.CreateOrganization)
	api.Get("/organizations/:id", container.OrganizationHandler.GetOrganizationByID)
}
