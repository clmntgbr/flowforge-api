package handler_test

import (
	"context"
	"errors"
	"forgeflow-api/dto"
	"forgeflow-api/handler"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWorkflowHandler_GetWorkflows_Success(t *testing.T) {
	app := newTestApp()
	mockService := NewMockWorkflowService()
	workflowHandler := handler.NewWorkflowHandler(mockService, NewMockStepService())

	orgID := uuid.New()
	expectedResponse := dto.PaginateResponse{
		Total:      2,
		Page:       1,
		Limit:      20,
		TotalPages: 1,
		Members: []dto.MinimalWorkflowOutput{
			{ID: uuid.New().String(), Name: "Workflow 1"},
			{ID: uuid.New().String(), Name: "Workflow 2"},
		},
	}

	mockService.GetWorkflowsFunc = func(c fiber.Ctx, organizationID uuid.UUID, query dto.PaginateQuery) (dto.PaginateResponse, error) {
		assert.Equal(t, orgID, organizationID)
		return expectedResponse, nil
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/workflows", workflowHandler.GetWorkflows)

	req, err := makeJSONRequest("GET", "/workflows?page=1&limit=20", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestWorkflowHandler_GetWorkflows_Unauthorized(t *testing.T) {
	app := newTestApp()
	mockService := NewMockWorkflowService()
	workflowHandler := handler.NewWorkflowHandler(mockService, NewMockStepService())

	app.Get("/workflows", workflowHandler.GetWorkflows)

	req, err := makeJSONRequest("GET", "/workflows", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestWorkflowHandler_CreateWorkflow_Success(t *testing.T) {
	app := newTestApp()
	mockService := NewMockWorkflowService()
	workflowHandler := handler.NewWorkflowHandler(mockService, NewMockStepService())

	orgID := uuid.New()
	validInput := dto.CreateWorkflowInput{
		Name:        "New Workflow",
		Description: "A test workflow",
	}

	mockService.CreateWorkflowFunc = func(c fiber.Ctx, organizationID uuid.UUID, req dto.CreateWorkflowInput) (dto.WorkflowOutput, error) {
		assert.Equal(t, orgID, organizationID)
		assert.Equal(t, validInput.Name, req.Name)
		return dto.WorkflowOutput{}, nil
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Post("/workflows", workflowHandler.CreateWorkflow)

	req, err := makeJSONRequest("POST", "/workflows", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestWorkflowHandler_CreateWorkflow_InvalidData(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
	}{
		{
			name:  "Empty name",
			input: map[string]interface{}{"name": "", "description": "Test"},
		},
		{
			name:  "Name too short",
			input: map[string]interface{}{"name": "A", "description": "Test"},
		},
		{
			name:  "Name too long",
			input: map[string]interface{}{"name": string(make([]byte, 300)), "description": "Test"},
		},
		{
			name:  "Description too short",
			input: map[string]interface{}{"name": "Valid Name", "description": "A"},
		},
		{
			name:  "Missing name",
			input: map[string]interface{}{"description": "Test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApp()
			mockService := NewMockWorkflowService()
			workflowHandler := handler.NewWorkflowHandler(mockService, NewMockStepService())

			orgID := uuid.New()

			app.Use(setOrganizationIDInContext(app, orgID))
			app.Post("/workflows", workflowHandler.CreateWorkflow)

			req, err := makeJSONRequest("POST", "/workflows", tt.input)
			assert.NoError(t, err)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestWorkflowHandler_GetWorkflowByID_Success(t *testing.T) {
	app := newTestApp()
	mockService := NewMockWorkflowService()
	workflowHandler := handler.NewWorkflowHandler(mockService, NewMockStepService())

	orgID := uuid.New()
	workflowID := uuid.New()
	expectedOutput := dto.WorkflowOutput{
		MinimalWorkflowOutput: dto.MinimalWorkflowOutput{
			ID:   workflowID.String(),
			Name: "Test Workflow",
		},
		Description: "Test Description",
	}

	mockService.GetWorkflowByIDFunc = func(c fiber.Ctx, organizationID uuid.UUID, wfID uuid.UUID) (dto.WorkflowOutput, error) {
		assert.Equal(t, orgID, organizationID)
		assert.Equal(t, workflowID, wfID)
		return expectedOutput, nil
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/workflows/:id", workflowHandler.GetWorkflowByID)

	req, err := makeJSONRequest("GET", "/workflows/"+workflowID.String(), nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestWorkflowHandler_GetWorkflowByID_InvalidUUID(t *testing.T) {
	app := newTestApp()
	mockService := NewMockWorkflowService()
	workflowHandler := handler.NewWorkflowHandler(mockService, NewMockStepService())

	orgID := uuid.New()

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/workflows/:id", workflowHandler.GetWorkflowByID)

	req, err := makeJSONRequest("GET", "/workflows/invalid-uuid", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusInternalServerError)
}

func TestWorkflowHandler_UpdateWorkflow_Success(t *testing.T) {
	app := newTestApp()
	mockService := NewMockWorkflowService()
	workflowHandler := handler.NewWorkflowHandler(mockService, NewMockStepService())

	orgID := uuid.New()
	workflowID := uuid.New()
	validInput := dto.UpdateWorkflowInput{
		Name:        "Updated Workflow",
		Description: "Updated Description",
	}

	mockService.UpdateWorkflowFunc = func(c fiber.Ctx, organizationID uuid.UUID, wfID uuid.UUID, req dto.UpdateWorkflowInput) (dto.WorkflowOutput, error) {
		assert.Equal(t, orgID, organizationID)
		assert.Equal(t, workflowID, wfID)
		assert.Equal(t, validInput.Name, req.Name)
		return dto.WorkflowOutput{}, nil
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Put("/workflows/:id", workflowHandler.UpdateWorkflow)

	req, err := makeJSONRequest("PUT", "/workflows/"+workflowID.String(), validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestWorkflowHandler_UpdateWorkflow_InvalidData(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
	}{
		{
			name:  "Empty name",
			input: map[string]interface{}{"name": "", "description": "Test"},
		},
		{
			name:  "Name too short",
			input: map[string]interface{}{"name": "A", "description": "Test"},
		},
		{
			name:  "Missing name",
			input: map[string]interface{}{"description": "Test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApp()
			mockService := NewMockWorkflowService()
			workflowHandler := handler.NewWorkflowHandler(mockService, NewMockStepService())

			orgID := uuid.New()
			workflowID := uuid.New()

			app.Use(setOrganizationIDInContext(app, orgID))
			app.Put("/workflows/:id", workflowHandler.UpdateWorkflow)

			req, err := makeJSONRequest("PUT", "/workflows/"+workflowID.String(), tt.input)
			assert.NoError(t, err)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestWorkflowHandler_UpdateWorkflowSteps_Success(t *testing.T) {
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
		assert.Equal(t, workflowID, wfID)
		assert.Len(t, steps, 1)
		return nil
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Put("/workflows/:id/steps", workflowHandler.UpdateWorkflowSteps)

	req, err := makeJSONRequest("PUT", "/workflows/"+workflowID.String()+"/steps", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestWorkflowHandler_ServiceError(t *testing.T) {
	app := newTestApp()
	mockService := NewMockWorkflowService()
	workflowHandler := handler.NewWorkflowHandler(mockService, NewMockStepService())

	orgID := uuid.New()
	serviceError := errors.New("database error")

	mockService.GetWorkflowsFunc = func(c fiber.Ctx, organizationID uuid.UUID, query dto.PaginateQuery) (dto.PaginateResponse, error) {
		return dto.PaginateResponse{}, serviceError
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/workflows", workflowHandler.GetWorkflows)

	req, err := makeJSONRequest("GET", "/workflows", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}
