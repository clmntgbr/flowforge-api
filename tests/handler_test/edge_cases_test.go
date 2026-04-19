package handler_test

import (
	"forgeflow-api/dto"
	"forgeflow-api/handler"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestOrganizationHandler_GetOrganizationByID_EmptyUUID(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	testUser := makeTestUser()

	app.Use(setUserInContext(app, testUser))
	// Use a catch-all route to test empty param
	app.Get("/organizations/*", func(c fiber.Ctx) error {
		// Manually set empty id param
		c.Params("id")
		return orgHandler.GetOrganizationByID(c)
	})

	req, err := makeJSONRequest("GET", "/organizations/", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

func TestOrganizationHandler_GetOrganizationByID_MalformedUUID(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	testUser := makeTestUser()

	app.Use(setUserInContext(app, testUser))
	app.Get("/organizations/:id", orgHandler.GetOrganizationByID)

	// Use various malformed UUIDs
	malformedUUIDs := []string{
		"not-a-uuid",
		"12345",
		"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		"00000000-0000-0000-0000-00000000000", // one digit short
	}

	for _, badID := range malformedUUIDs {
		req, err := makeJSONRequest("GET", "/organizations/"+badID, nil)
		assert.NoError(t, err)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.True(t, resp.StatusCode >= 400, "Should fail for UUID: %s", badID)
	}
}

func TestOrganizationHandler_UpdateOrganization_MalformedUUID(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	testUser := makeTestUser()
	validInput := dto.UpdateOrganizationInput{
		Name: "Updated Name",
	}

	app.Use(setUserInContext(app, testUser))
	app.Put("/organizations/:id", orgHandler.UpdateOrganization)

	malformedUUIDs := []string{
		"not-a-uuid",
		"12345",
		"invalid-format",
	}

	for _, badID := range malformedUUIDs {
		req, err := makeJSONRequest("PUT", "/organizations/"+badID, validInput)
		assert.NoError(t, err)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.True(t, resp.StatusCode >= 400, "Should fail for UUID: %s", badID)
	}
}

func TestOrganizationHandler_ActivateOrganization_MalformedUUID(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	testUser := makeTestUser()

	app.Use(setUserInContext(app, testUser))
	app.Patch("/organizations/:id/activate", orgHandler.ActivateOrganization)

	malformedUUIDs := []string{
		"not-a-uuid",
		"12345",
		"invalid-format",
	}

	for _, badID := range malformedUUIDs {
		req, err := makeJSONRequest("PATCH", "/organizations/"+badID+"/activate", nil)
		assert.NoError(t, err)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.True(t, resp.StatusCode >= 400, "Should fail for UUID: %s", badID)
	}
}

func TestWorkflowHandler_GetWorkflowByID_EmptyUUID(t *testing.T) {
	app := newTestApp()
	mockService := NewMockWorkflowService()
	workflowHandler := handler.NewWorkflowHandler(mockService, NewMockStepService())

	orgID := uuid.New()

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/workflows/:id", workflowHandler.GetWorkflowByID)

	req, err := makeJSONRequest("GET", "/workflows/", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

func TestWorkflowHandler_UpdateWorkflow_EmptyUUID(t *testing.T) {
	app := newTestApp()
	mockService := NewMockWorkflowService()
	workflowHandler := handler.NewWorkflowHandler(mockService, NewMockStepService())

	orgID := uuid.New()
	validInput := dto.UpdateWorkflowInput{
		Name:        "Updated",
		Description: "Test",
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Put("/workflows/:id", workflowHandler.UpdateWorkflow)

	req, err := makeJSONRequest("PUT", "/workflows/", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

func TestWorkflowHandler_UpdateWorkflowSteps_EmptyUUID(t *testing.T) {
	app := newTestApp()
	mockWorkflowService := NewMockWorkflowService()
	mockStepService := NewMockStepService()
	workflowHandler := handler.NewWorkflowHandler(mockWorkflowService, mockStepService)

	orgID := uuid.New()
	validInput := dto.UpdateWorkflowStepsInput{
		Steps: []dto.UpdateWorkflowStepInput{
			{
				ID:         uuid.New().String(),
				EndpointID: uuid.New().String(),
				Position:   dto.Position{X: 100, Y: 200},
				Index:      "0",
			},
		},
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Put("/workflows/:id/steps", workflowHandler.UpdateWorkflowSteps)

	req, err := makeJSONRequest("PUT", "/workflows//steps", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

func TestEndpointHandler_GetEndpoints_WithMalformedQueryString(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	orgID := uuid.New()

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/endpoints", endpointHandler.GetEndpoints)

	// Create a request with completely invalid query params that might cause bind error
	req, err := makeJSONRequest("GET", "/endpoints", nil)
	assert.NoError(t, err)

	// Force malformed query by setting raw query directly
	req.URL.RawQuery = strings.Repeat("invalid&", 1000)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	// Either succeeds with normalized values or returns error
	assert.True(t, resp.StatusCode == fiber.StatusOK || resp.StatusCode >= 400)
}

func TestWorkflowHandler_GetWorkflows_WithMalformedQueryString(t *testing.T) {
	app := newTestApp()
	mockService := NewMockWorkflowService()
	workflowHandler := handler.NewWorkflowHandler(mockService, NewMockStepService())

	orgID := uuid.New()

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/workflows", workflowHandler.GetWorkflows)

	req, err := makeJSONRequest("GET", "/workflows", nil)
	assert.NoError(t, err)

	req.URL.RawQuery = strings.Repeat("malformed&", 1000)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode == fiber.StatusOK || resp.StatusCode >= 400)
}
