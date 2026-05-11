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
)

func main() {
	env := config.Load()
	config.ConnectDatabase(env)

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

	setupHealthChecks(app)

	log.Println("🚀 Server is running on port", env.Port)
	log.Fatal(app.Listen(":" + env.Port))
}

func setupHealthChecks(app *fiber.App) {
	app.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	app.Get(healthcheck.ReadinessEndpoint, healthcheck.New())
	app.Get(healthcheck.StartupEndpoint, healthcheck.New())
}
