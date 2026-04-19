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

func TestConnexionHandler_CreateConnexion_Success(t *testing.T) {
	app := newTestApp()
	mockConnexionService := NewMockConnexionService()
	mockWorkflowService := NewMockWorkflowService()
	connexionHandler := handler.NewConnexionHandler(mockConnexionService, mockWorkflowService)

	orgID := uuid.New()
	workflowID := uuid.New()
	fromStepID := uuid.New()
	toStepID := uuid.New()
	
	validInput := dto.CreateConnexionInput{
		WorkflowID: workflowID,
		From:       fromStepID.String(),
		To:         toStepID.String(),
	}

	mockWorkflowService.GetWorkflowByIDFunc = func(c fiber.Ctx, organizationID uuid.UUID, wfID uuid.UUID) (dto.WorkflowOutput, error) {
		assert.Equal(t, orgID, organizationID)
		assert.Equal(t, workflowID, wfID)
		return dto.WorkflowOutput{}, nil
	}

	mockConnexionService.CreateConnexionFunc = func(c fiber.Ctx, wfID uuid.UUID, req dto.CreateConnexionInput) (dto.ConnexionOutput, error) {
		assert.Equal(t, workflowID, wfID)
		assert.Equal(t, validInput.From, req.From)
		assert.Equal(t, validInput.To, req.To)
		return dto.ConnexionOutput{}, nil
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Post("/connexions", connexionHandler.CreateConnexion)

	req, err := makeJSONRequest("POST", "/connexions", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestConnexionHandler_CreateConnexion_Unauthorized(t *testing.T) {
	app := newTestApp()
	mockConnexionService := NewMockConnexionService()
	mockWorkflowService := NewMockWorkflowService()
	connexionHandler := handler.NewConnexionHandler(mockConnexionService, mockWorkflowService)

	workflowID := uuid.New()
	validInput := dto.CreateConnexionInput{
		WorkflowID: workflowID,
		From:       uuid.New().String(),
		To:         uuid.New().String(),
	}

	app.Post("/connexions", connexionHandler.CreateConnexion)

	req, err := makeJSONRequest("POST", "/connexions", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestConnexionHandler_CreateConnexion_InvalidData(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
	}{
		{
			name:  "Missing workflowId",
			input: map[string]interface{}{"from": uuid.New().String(), "to": uuid.New().String()},
		},
		{
			name:  "Missing from",
			input: map[string]interface{}{"workflowId": uuid.New().String(), "to": uuid.New().String()},
		},
		{
			name:  "Missing to",
			input: map[string]interface{}{"workflowId": uuid.New().String(), "from": uuid.New().String()},
		},
		{
			name:  "Invalid from UUID",
			input: map[string]interface{}{"workflowId": uuid.New().String(), "from": "invalid", "to": uuid.New().String()},
		},
		{
			name:  "Invalid to UUID",
			input: map[string]interface{}{"workflowId": uuid.New().String(), "from": uuid.New().String(), "to": "invalid"},
		},
		{
			name:  "All fields missing",
			input: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApp()
			mockConnexionService := NewMockConnexionService()
			mockWorkflowService := NewMockWorkflowService()
			connexionHandler := handler.NewConnexionHandler(mockConnexionService, mockWorkflowService)

			orgID := uuid.New()

			app.Use(setOrganizationIDInContext(app, orgID))
			app.Post("/connexions", connexionHandler.CreateConnexion)

			req, err := makeJSONRequest("POST", "/connexions", tt.input)
			assert.NoError(t, err)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestConnexionHandler_CreateConnexion_WorkflowNotFound(t *testing.T) {
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
		return dto.WorkflowOutput{}, errors.New("workflow not found")
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Post("/connexions", connexionHandler.CreateConnexion)

	req, err := makeJSONRequest("POST", "/connexions", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

func TestConnexionHandler_DeleteConnexion_Success(t *testing.T) {
	app := newTestApp()
	mockConnexionService := NewMockConnexionService()
	mockWorkflowService := NewMockWorkflowService()
	connexionHandler := handler.NewConnexionHandler(mockConnexionService, mockWorkflowService)

	connexionID := uuid.New()

	mockConnexionService.DeleteConnexionFunc = func(ctx context.Context, cID uuid.UUID) (dto.ConnexionOutput, error) {
		assert.Equal(t, connexionID, cID)
		return dto.ConnexionOutput{}, nil
	}

	app.Delete("/connexions/:id", connexionHandler.DeleteConnexion)

	req, err := makeJSONRequest("DELETE", "/connexions/"+connexionID.String(), nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestConnexionHandler_DeleteConnexion_InvalidUUID(t *testing.T) {
	app := newTestApp()
	mockConnexionService := NewMockConnexionService()
	mockWorkflowService := NewMockWorkflowService()
	connexionHandler := handler.NewConnexionHandler(mockConnexionService, mockWorkflowService)

	app.Delete("/connexions/:id", connexionHandler.DeleteConnexion)

	req, err := makeJSONRequest("DELETE", "/connexions/invalid-uuid", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusInternalServerError)
}

func TestConnexionHandler_DeleteConnexion_ServiceError(t *testing.T) {
	app := newTestApp()
	mockConnexionService := NewMockConnexionService()
	mockWorkflowService := NewMockWorkflowService()
	connexionHandler := handler.NewConnexionHandler(mockConnexionService, mockWorkflowService)

	connexionID := uuid.New()
	serviceError := errors.New("database error")

	mockConnexionService.DeleteConnexionFunc = func(ctx context.Context, cID uuid.UUID) (dto.ConnexionOutput, error) {
		return dto.ConnexionOutput{}, serviceError
	}

	app.Delete("/connexions/:id", connexionHandler.DeleteConnexion)

	req, err := makeJSONRequest("DELETE", "/connexions/"+connexionID.String(), nil)
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
