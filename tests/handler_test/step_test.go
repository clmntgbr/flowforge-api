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
	"gorm.io/datatypes"
)

func TestStepHandler_GetStepByID_Success(t *testing.T) {
	app := newTestApp()
	mockService := NewMockStepService()
	stepHandler := handler.NewStepHandler(mockService)

	orgID := uuid.New()
	stepID := uuid.New()
	expectedOutput := dto.StepOutput{
		MinimalStepOutput: dto.MinimalStepOutput{
			ID:   stepID.String(),
			Name: "Test Step",
		},
		Description: "Test Description",
		Timeout:     5000,
	}

	mockService.GetStepByIDFunc = func(ctx context.Context, organizationID uuid.UUID, sID uuid.UUID) (dto.StepOutput, error) {
		assert.Equal(t, orgID, organizationID)
		assert.Equal(t, stepID, sID)
		return expectedOutput, nil
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/steps/:id", stepHandler.GetStepByID)

	req, err := makeJSONRequest("GET", "/steps/"+stepID.String(), nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestStepHandler_GetStepByID_Unauthorized(t *testing.T) {
	app := newTestApp()
	mockService := NewMockStepService()
	stepHandler := handler.NewStepHandler(mockService)

	stepID := uuid.New()

	app.Get("/steps/:id", stepHandler.GetStepByID)

	req, err := makeJSONRequest("GET", "/steps/"+stepID.String(), nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestStepHandler_GetStepByID_InvalidUUID(t *testing.T) {
	app := newTestApp()
	mockService := NewMockStepService()
	stepHandler := handler.NewStepHandler(mockService)

	orgID := uuid.New()

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/steps/:id", stepHandler.GetStepByID)

	req, err := makeJSONRequest("GET", "/steps/invalid-uuid", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusInternalServerError)
}

func TestStepHandler_UpdateStep_Success(t *testing.T) {
	app := newTestApp()
	mockService := NewMockStepService()
	stepHandler := handler.NewStepHandler(mockService)

	orgID := uuid.New()
	stepID := uuid.New()
	validInput := dto.UpdateStepInput{
		Name:           "Updated Step",
		Description:    "Updated Description",
		Timeout:        10000,
		Query:          datatypes.JSON([]byte(`{"key": "value"}`)),
		RetryOnFailure: true,
		RetryCount:     5,
		RetryDelay:     2000,
	}

	mockService.UpdateStepFunc = func(ctx context.Context, organizationID uuid.UUID, sID uuid.UUID, req dto.UpdateStepInput) (dto.StepOutput, error) {
		assert.Equal(t, orgID, organizationID)
		assert.Equal(t, stepID, sID)
		assert.Equal(t, validInput.Name, req.Name)
		return dto.StepOutput{}, nil
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Put("/steps/:id", stepHandler.UpdateStep)

	req, err := makeJSONRequest("PUT", "/steps/"+stepID.String(), validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestStepHandler_UpdateStep_InvalidData(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
	}{
		{
			name:  "Empty name",
			input: map[string]interface{}{"name": "", "description": "Test", "timeout": 5000},
		},
		{
			name:  "Name too short",
			input: map[string]interface{}{"name": "A", "description": "Test", "timeout": 5000},
		},
		{
			name:  "Name too long",
			input: map[string]interface{}{"name": string(make([]byte, 150)), "description": "Test", "timeout": 5000},
		},
		{
			name:  "Missing name",
			input: map[string]interface{}{"description": "Test", "timeout": 5000},
		},
		{
			name:  "Negative timeout",
			input: map[string]interface{}{"name": "Valid Name", "description": "Test", "timeout": -100},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApp()
			mockService := NewMockStepService()
			stepHandler := handler.NewStepHandler(mockService)

			orgID := uuid.New()
			stepID := uuid.New()

			app.Use(setOrganizationIDInContext(app, orgID))
			app.Put("/steps/:id", stepHandler.UpdateStep)

			req, err := makeJSONRequest("PUT", "/steps/"+stepID.String(), tt.input)
			assert.NoError(t, err)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestStepHandler_UpdateStep_InvalidUUID(t *testing.T) {
	app := newTestApp()
	mockService := NewMockStepService()
	stepHandler := handler.NewStepHandler(mockService)

	orgID := uuid.New()
	validInput := dto.UpdateStepInput{
		Name:           "Updated",
		Description:    "Test",
		Timeout:        5000,
		RetryOnFailure: false,
		RetryCount:     3,
		RetryDelay:     1000,
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Put("/steps/:id", stepHandler.UpdateStep)

	req, err := makeJSONRequest("PUT", "/steps/invalid-uuid", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusInternalServerError)
}

func TestStepHandler_ServiceError(t *testing.T) {
	app := newTestApp()
	mockService := NewMockStepService()
	stepHandler := handler.NewStepHandler(mockService)

	orgID := uuid.New()
	stepID := uuid.New()
	serviceError := errors.New("database error")

	mockService.GetStepByIDFunc = func(ctx context.Context, organizationID uuid.UUID, sID uuid.UUID) (dto.StepOutput, error) {
		return dto.StepOutput{}, serviceError
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/steps/:id", stepHandler.GetStepByID)

	req, err := makeJSONRequest("GET", "/steps/"+stepID.String(), nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}
