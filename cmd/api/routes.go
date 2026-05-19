package main

import (
	"flowforge-api/cmd/api/wire"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
)

func setupRoutes(app *fiber.App, container *wire.Container) {
	setupHealthChecks(app)
	setupWebhooks(app, container)
	setupAPIRoutes(app, container)
}

func setupWebhooks(app *fiber.App, container *wire.Container) {
	webhooks := app.Group("/webhook")

	webhooks.Post("/clerk", container.ClerkMiddleware.Protected(), container.ClerkHandler.Execute)
}

func setupHealthChecks(app *fiber.App) {
	app.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	app.Get(healthcheck.ReadinessEndpoint, healthcheck.New())
	app.Get(healthcheck.StartupEndpoint, healthcheck.New())
}

func setupAPIRoutes(app *fiber.App, container *wire.Container) {
	api := app.Group("/api")

	api.Use(container.AuthenticateMiddleware.Protected())
	setupUsersRoutes(api, container)
	setupOrganizationsRoutes(api, container)
	setupEndpointsRoutes(api, container)
	setupConnexionsRoutes(api, container)
	setupStepsRoutes(api, container)
	setupWorkflowsRoutes(api, container)
}

func setupUsersRoutes(api fiber.Router, container *wire.Container) {
	api.Get("/users/me", container.UserHandler.GetUser)
}

func setupOrganizationsRoutes(api fiber.Router, container *wire.Container) {
	api.Get("/organizations", container.OrganizationHandler.GetOrganizations)
	api.Post("/organizations", container.OrganizationHandler.CreateOrganization)
	api.Get("/organizations/:id", container.OrganizationHandler.GetOrganizationByID)
	api.Put("/organizations/:id", container.OrganizationHandler.UpdateOrganization)
	api.Put("/organizations/:id/activate", container.OrganizationHandler.ActivateOrganization)
}

func setupEndpointsRoutes(api fiber.Router, container *wire.Container) {
	api.Get("/endpoints", container.EndpointHandler.GetEndpoints)
	api.Post("/endpoints", container.EndpointHandler.CreateEndpoint)
	api.Get("/endpoints/:id", container.EndpointHandler.GetEndpointByID)
	api.Put("/endpoints/:id", container.EndpointHandler.UpdateEndpoint)
}

func setupConnexionsRoutes(api fiber.Router, container *wire.Container) {
	api.Delete("/connexions/:id", container.ConnexionHandler.DeleteConnexion)
	api.Post("/connexions", container.ConnexionHandler.CreateConnexion)
}

func setupStepsRoutes(api fiber.Router, container *wire.Container) {
	api.Get("/workflows/:workflowId/steps/:id", container.StepHandler.GetStepByID)
	api.Put("/workflows/:workflowId/steps/:id", container.StepHandler.UpdateStep)
	api.Delete("/workflows/:workflowId/steps/:id", container.StepHandler.DeleteStep)
}

func setupWorkflowsRoutes(api fiber.Router, container *wire.Container) {
	api.Get("/workflows", container.WorkflowHandler.GetWorkflows)
	api.Post("/workflows", container.WorkflowHandler.CreateWorkflow)
	api.Get("/workflows/:id", container.WorkflowHandler.GetWorkflowByID)
	api.Put("/workflows/:id", container.WorkflowHandler.UpdateWorkflow)
	api.Put("/workflows/:id/activate", container.WorkflowHandler.ActivateWorkflow)
	api.Put("/workflows/:id/deactivate", container.WorkflowHandler.DeactivateWorkflow)
	api.Put("/workflows/:id/upsert", container.WorkflowHandler.UpsertWorkflow)
	api.Get("/workflows/:id/runs", container.WorkflowHandler.GetWorkflowRuns)
}
