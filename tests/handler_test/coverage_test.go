package handler_test

import (
	"context"
	"errors"
	"forgeflow-api/domain"
	"forgeflow-api/dto"
	"forgeflow-api/handler"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBaseHandler_BindAndValidate_InvalidJSON(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	orgID := uuid.New()

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Post("/endpoints", endpointHandler.CreateEndpoint)

	// Send invalid JSON
	req, err := makeJSONRequest("POST", "/endpoints", nil)
	assert.NoError(t, err)
	req.Body = http.NoBody
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestBaseHandler_ParseUUIDParam_EmptyParam(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	orgID := uuid.New()

	app.Use(setOrganizationIDInContext(app, orgID))
	// Route without :id parameter
	app.Get("/endpoints/", endpointHandler.GetEndpointByID)

	req, err := makeJSONRequest("GET", "/endpoints/", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

func TestBaseHandler_BindPaginateQuery_InvalidQuery(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	orgID := uuid.New()

	mockService.GetEndpointsFunc = func(c fiber.Ctx, organizationID uuid.UUID, query dto.PaginateQuery) (dto.PaginateResponse, error) {
		return dto.PaginateResponse{
			Total:      0,
			Page:       query.Page,
			Limit:      query.Limit,
			TotalPages: 0,
			Members:    []dto.MinimalEndpointOutput{},
		}, nil
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/endpoints", endpointHandler.GetEndpoints)

	// Test with invalid page parameter (will normalize)
	req, err := makeJSONRequest("GET", "/endpoints?page=-1&limit=1000", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestOrganizationHandler_GetOrganizationByID_ServiceError(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	testUser := makeTestUser()
	orgID := uuid.New()

	mockService.GetOrganizationByIDFunc = func(c fiber.Ctx, user *domain.User, organizationID uuid.UUID) (dto.OrganizationOutput, error) {
		return dto.OrganizationOutput{}, errors.New("database error")
	}

	app.Use(setUserInContext(app, testUser))
	app.Get("/organizations/:id", orgHandler.GetOrganizationByID)

	req, err := makeJSONRequest("GET", "/organizations/"+orgID.String(), nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

func TestOrganizationHandler_CreateOrganization_ServiceError(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	testUser := makeTestUser()
	validInput := dto.CreateOrganizationInput{
		Name: "Test Organization",
	}

	mockService.CreateOrganizationFunc = func(c fiber.Ctx, user *domain.User, name string) (dto.OrganizationOutput, error) {
		return dto.OrganizationOutput{}, errors.New("max organizations reached")
	}

	app.Use(setUserInContext(app, testUser))
	app.Post("/organizations", orgHandler.CreateOrganization)

	req, err := makeJSONRequest("POST", "/organizations", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

func TestOrganizationHandler_UpdateOrganization_ServiceError(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	testUser := makeTestUser()
	orgID := uuid.New()
	validInput := dto.UpdateOrganizationInput{
		Name: "Updated Name",
	}

	mockService.UpdateOrganizationFunc = func(c fiber.Ctx, user *domain.User, organizationID uuid.UUID, req dto.UpdateOrganizationInput) (dto.OrganizationOutput, error) {
		return dto.OrganizationOutput{}, errors.New("organization not found")
	}

	app.Use(setUserInContext(app, testUser))
	app.Put("/organizations/:id", orgHandler.UpdateOrganization)

	req, err := makeJSONRequest("PUT", "/organizations/"+orgID.String(), validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

func TestOrganizationHandler_ActivateOrganization_ServiceError(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	testUser := makeTestUser()
	orgID := uuid.New()

	mockService.ActivateOrganizationFunc = func(ctx context.Context, userID uuid.UUID, organizationID uuid.UUID) (dto.OrganizationOutput, error) {
		return dto.OrganizationOutput{}, errors.New("organization not found")
	}

	app.Use(setUserInContext(app, testUser))
	app.Patch("/organizations/:id/activate", orgHandler.ActivateOrganization)

	req, err := makeJSONRequest("PATCH", "/organizations/"+orgID.String()+"/activate", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

func TestWorkflowHandler_CreateWorkflow_ServiceError(t *testing.T) {
	app := newTestApp()
	mockService := NewMockWorkflowService()
	workflowHandler := handler.NewWorkflowHandler(mockService, NewMockStepService())

	orgID := uuid.New()
	validInput := dto.CreateWorkflowInput{
		Name:        "Test Workflow",
		Description: "Test",
	}

	mockService.CreateWorkflowFunc = func(c fiber.Ctx, organizationID uuid.UUID, req dto.CreateWorkflowInput) (dto.WorkflowOutput, error) {
		return dto.WorkflowOutput{}, errors.New("failed to create workflow")
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Post("/workflows", workflowHandler.CreateWorkflow)

	req, err := makeJSONRequest("POST", "/workflows", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

func TestWorkflowHandler_GetWorkflowByID_ServiceError(t *testing.T) {
	app := newTestApp()
	mockService := NewMockWorkflowService()
	workflowHandler := handler.NewWorkflowHandler(mockService, NewMockStepService())

	orgID := uuid.New()
	workflowID := uuid.New()

	mockService.GetWorkflowByIDFunc = func(c fiber.Ctx, organizationID uuid.UUID, wfID uuid.UUID) (dto.WorkflowOutput, error) {
		return dto.WorkflowOutput{}, errors.New("workflow not found")
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/workflows/:id", workflowHandler.GetWorkflowByID)

	req, err := makeJSONRequest("GET", "/workflows/"+workflowID.String(), nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

func TestWorkflowHandler_UpdateWorkflow_ServiceError(t *testing.T) {
	app := newTestApp()
	mockService := NewMockWorkflowService()
	workflowHandler := handler.NewWorkflowHandler(mockService, NewMockStepService())

	orgID := uuid.New()
	workflowID := uuid.New()
	validInput := dto.UpdateWorkflowInput{
		Name:        "Updated",
		Description: "Test",
	}

	mockService.UpdateWorkflowFunc = func(c fiber.Ctx, organizationID uuid.UUID, wfID uuid.UUID, req dto.UpdateWorkflowInput) (dto.WorkflowOutput, error) {
		return dto.WorkflowOutput{}, errors.New("failed to update workflow")
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Put("/workflows/:id", workflowHandler.UpdateWorkflow)

	req, err := makeJSONRequest("PUT", "/workflows/"+workflowID.String(), validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

func TestWorkflowHandler_UpdateWorkflowSteps_WorkflowNotFound(t *testing.T) {
	app := newTestApp()
	mockWorkflowService := NewMockWorkflowService()
	mockStepService := NewMockStepService()
	workflowHandler := handler.NewWorkflowHandler(mockWorkflowService, mockStepService)

	orgID := uuid.New()
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

	mockWorkflowService.GetWorkflowByIDFunc = func(c fiber.Ctx, organizationID uuid.UUID, wfID uuid.UUID) (dto.WorkflowOutput, error) {
		return dto.WorkflowOutput{}, errors.New("workflow not found")
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Put("/workflows/:id/steps", workflowHandler.UpdateWorkflowSteps)

	req, err := makeJSONRequest("PUT", "/workflows/"+workflowID.String()+"/steps", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

func TestWorkflowHandler_UpdateWorkflowSteps_StepServiceError(t *testing.T) {
	app := newTestApp()
	mockWorkflowService := NewMockWorkflowService()
	mockStepService := NewMockStepService()
	workflowHandler := handler.NewWorkflowHandler(mockWorkflowService, mockStepService)

	orgID := uuid.New()
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

	mockWorkflowService.GetWorkflowByIDFunc = func(c fiber.Ctx, organizationID uuid.UUID, wfID uuid.UUID) (dto.WorkflowOutput, error) {
		return dto.WorkflowOutput{}, nil
	}

	mockStepService.UpsertStepsFunc = func(ctx context.Context, wfID uuid.UUID, steps []dto.UpdateWorkflowStepInput) error {
		return errors.New("failed to upsert steps")
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Put("/workflows/:id/steps", workflowHandler.UpdateWorkflowSteps)

	req, err := makeJSONRequest("PUT", "/workflows/"+workflowID.String()+"/steps", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

func TestStepHandler_UpdateStep_ServiceError(t *testing.T) {
	app := newTestApp()
	mockService := NewMockStepService()
	stepHandler := handler.NewStepHandler(mockService)

	orgID := uuid.New()
	stepID := uuid.New()
	validInput := dto.UpdateStepInput{
		Name:        "Updated Step",
		Description: "Test",
		Timeout:     5000,
	}

	mockService.UpdateStepFunc = func(ctx context.Context, organizationID uuid.UUID, sID uuid.UUID, req dto.UpdateStepInput) (dto.StepOutput, error) {
		return dto.StepOutput{}, errors.New("failed to update step")
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Put("/steps/:id", stepHandler.UpdateStep)

	req, err := makeJSONRequest("PUT", "/steps/"+stepID.String(), validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

func TestConnexionHandler_CreateConnexion_ServiceError(t *testing.T) {
	app := newTestApp()
	mockConnexionService := NewMockConnexionService()
	mockWorkflowService := NewMockWorkflowService()
	connexionHandler := handler.NewConnexionHandler(mockConnexionService, mockWorkflowService)

	orgID := uuid.New()
	workflowID := uuid.New()

	validInput := dto.CreateConnexionInput{
		WorkflowID: workflowID,
		From:       uuid.New().String(),
		To:         uuid.New().String(),
	}

	mockWorkflowService.GetWorkflowByIDFunc = func(c fiber.Ctx, organizationID uuid.UUID, wfID uuid.UUID) (dto.WorkflowOutput, error) {
		return dto.WorkflowOutput{}, nil
	}

	mockConnexionService.CreateConnexionFunc = func(c fiber.Ctx, wfID uuid.UUID, req dto.CreateConnexionInput) (dto.ConnexionOutput, error) {
		return dto.ConnexionOutput{}, errors.New("failed to create connexion")
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Post("/connexions", connexionHandler.CreateConnexion)

	req, err := makeJSONRequest("POST", "/connexions", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

func TestEndpointHandler_GetEndpoints_BadPaginationQuery(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	orgID := uuid.New()

	mockService.GetEndpointsFunc = func(c fiber.Ctx, organizationID uuid.UUID, query dto.PaginateQuery) (dto.PaginateResponse, error) {
		return dto.PaginateResponse{
			Total:      0,
			Page:       query.Page,
			Limit:      query.Limit,
			TotalPages: 0,
			Members:    []dto.MinimalEndpointOutput{},
		}, nil
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/endpoints", endpointHandler.GetEndpoints)

	// Test with various edge cases
	req, err := makeJSONRequest("GET", "/endpoints?page=0&limit=0", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestWorkflowHandler_GetWorkflows_BadPaginationQuery(t *testing.T) {
	app := newTestApp()
	mockService := NewMockWorkflowService()
	workflowHandler := handler.NewWorkflowHandler(mockService, NewMockStepService())

	orgID := uuid.New()

	mockService.GetWorkflowsFunc = func(c fiber.Ctx, organizationID uuid.UUID, query dto.PaginateQuery) (dto.PaginateResponse, error) {
		return dto.PaginateResponse{
			Total:      0,
			Page:       query.Page,
			Limit:      query.Limit,
			TotalPages: 0,
			Members:    []dto.MinimalWorkflowOutput{},
		}, nil
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/workflows", workflowHandler.GetWorkflows)

	// Test with edge cases that get normalized
	req, err := makeJSONRequest("GET", "/workflows?page=0&limit=0", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
