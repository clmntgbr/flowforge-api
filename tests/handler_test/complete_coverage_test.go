package handler_test

import (
	"forgeflow-api/dto"
	"forgeflow-api/handler"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Tests pour atteindre 100% de coverage

func TestOrganizationHandler_GetOrganizationByID_Unauthorized(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	orgID := uuid.New()

	// Pas de user dans le contexte
	app.Get("/organizations/:id", orgHandler.GetOrganizationByID)

	req, err := makeJSONRequest("GET", "/organizations/"+orgID.String(), nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestOrganizationHandler_CreateOrganization_Unauthorized(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	validInput := dto.CreateOrganizationInput{
		Name: "Test Organization",
	}

	// Pas de user dans le contexte
	app.Post("/organizations", orgHandler.CreateOrganization)

	req, err := makeJSONRequest("POST", "/organizations", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestOrganizationHandler_UpdateOrganization_Unauthorized(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	orgID := uuid.New()
	validInput := dto.UpdateOrganizationInput{
		Name: "Updated Name",
	}

	// Pas de user dans le contexte
	app.Put("/organizations/:id", orgHandler.UpdateOrganization)

	req, err := makeJSONRequest("PUT", "/organizations/"+orgID.String(), validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestOrganizationHandler_ActivateOrganization_Unauthorized(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	orgID := uuid.New()

	// Pas de user dans le contexte
	app.Patch("/organizations/:id/activate", orgHandler.ActivateOrganization)

	req, err := makeJSONRequest("PATCH", "/organizations/"+orgID.String()+"/activate", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestStepHandler_UpdateStep_Unauthorized(t *testing.T) {
	app := newTestApp()
	mockService := NewMockStepService()
	stepHandler := handler.NewStepHandler(mockService)

	stepID := uuid.New()
	validInput := dto.UpdateStepInput{
		Name:        "Updated Step",
		Description: "Test",
		Timeout:     5000,
	}

	// Pas d'organizationID dans le contexte
	app.Put("/steps/:id", stepHandler.UpdateStep)

	req, err := makeJSONRequest("PUT", "/steps/"+stepID.String(), validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestWorkflowHandler_CreateWorkflow_Unauthorized(t *testing.T) {
	app := newTestApp()
	mockService := NewMockWorkflowService()
	workflowHandler := handler.NewWorkflowHandler(mockService, NewMockStepService())

	validInput := dto.CreateWorkflowInput{
		Name:        "Test Workflow",
		Description: "Test",
	}

	// Pas d'organizationID dans le contexte
	app.Post("/workflows", workflowHandler.CreateWorkflow)

	req, err := makeJSONRequest("POST", "/workflows", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestWorkflowHandler_GetWorkflowByID_Unauthorized(t *testing.T) {
	app := newTestApp()
	mockService := NewMockWorkflowService()
	workflowHandler := handler.NewWorkflowHandler(mockService, NewMockStepService())

	workflowID := uuid.New()

	// Pas d'organizationID dans le contexte
	app.Get("/workflows/:id", workflowHandler.GetWorkflowByID)

	req, err := makeJSONRequest("GET", "/workflows/"+workflowID.String(), nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestWorkflowHandler_UpdateWorkflow_Unauthorized(t *testing.T) {
	app := newTestApp()
	mockService := NewMockWorkflowService()
	workflowHandler := handler.NewWorkflowHandler(mockService, NewMockStepService())

	workflowID := uuid.New()
	validInput := dto.UpdateWorkflowInput{
		Name:        "Updated",
		Description: "Test",
	}

	// Pas d'organizationID dans le contexte
	app.Put("/workflows/:id", workflowHandler.UpdateWorkflow)

	req, err := makeJSONRequest("PUT", "/workflows/"+workflowID.String(), validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestWorkflowHandler_UpdateWorkflowSteps_Unauthorized(t *testing.T) {
	app := newTestApp()
	mockWorkflowService := NewMockWorkflowService()
	mockStepService := NewMockStepService()
	workflowHandler := handler.NewWorkflowHandler(mockWorkflowService, mockStepService)

	workflowID := uuid.New()
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

	// Pas d'organizationID dans le contexte
	app.Put("/workflows/:id/steps", workflowHandler.UpdateWorkflowSteps)

	req, err := makeJSONRequest("PUT", "/workflows/"+workflowID.String()+"/steps", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestWorkflowHandler_UpdateWorkflowSteps_InvalidData(t *testing.T) {
	app := newTestApp()
	mockWorkflowService := NewMockWorkflowService()
	mockStepService := NewMockStepService()
	workflowHandler := handler.NewWorkflowHandler(mockWorkflowService, mockStepService)

	orgID := uuid.New()
	workflowID := uuid.New()
	invalidInput := map[string]interface{}{
		"steps": "not an array",
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Put("/workflows/:id/steps", workflowHandler.UpdateWorkflowSteps)

	req, err := makeJSONRequest("PUT", "/workflows/"+workflowID.String()+"/steps", invalidInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestWorkflowHandler_GetWorkflows_BadRequest(t *testing.T) {
	app := newTestApp()
	mockService := NewMockWorkflowService()
	workflowHandler := handler.NewWorkflowHandler(mockService, NewMockStepService())

	orgID := uuid.New()

	mockService.GetWorkflowsFunc = func(c fiber.Ctx, organizationID uuid.UUID, query dto.PaginateQuery) (dto.PaginateResponse, error) {
		return dto.PaginateResponse{}, fiber.ErrBadRequest
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/workflows", workflowHandler.GetWorkflows)

	req, err := makeJSONRequest("GET", "/workflows", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}
