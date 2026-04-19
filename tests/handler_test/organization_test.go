package handler_test

import (
	"context"
	"errors"
	"forgeflow-api/domain"
	"forgeflow-api/dto"
	"forgeflow-api/handler"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestOrganizationHandler_GetOrganizations_Success(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	testUser := makeTestUser()
	orgID := uuid.New()
	expectedOutput := []dto.MinimalOrganizationOutput{
		{ID: uuid.New().String(), Name: "Org 1", IsActive: true},
		{ID: uuid.New().String(), Name: "Org 2", IsActive: false},
	}

	mockService.GetOrganizationsFunc = func(c fiber.Ctx, user *domain.User, activeOrganizationID uuid.UUID) ([]dto.MinimalOrganizationOutput, error) {
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, orgID, activeOrganizationID)
		return expectedOutput, nil
	}

	app.Use(setUserInContext(app, testUser))
	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/organizations", orgHandler.GetOrganizations)

	req, err := makeJSONRequest("GET", "/organizations", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestOrganizationHandler_GetOrganizations_Unauthorized_NoUser(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	orgID := uuid.New()

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/organizations", orgHandler.GetOrganizations)

	req, err := makeJSONRequest("GET", "/organizations", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestOrganizationHandler_GetOrganizations_Unauthorized_NoOrgID(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	testUser := makeTestUser()

	app.Use(setUserInContext(app, testUser))
	app.Get("/organizations", orgHandler.GetOrganizations)

	req, err := makeJSONRequest("GET", "/organizations", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestOrganizationHandler_CreateOrganization_Success(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	testUser := makeTestUser()
	validInput := dto.CreateOrganizationInput{
		Name: "My New Organization",
	}

	mockService.CreateOrganizationFunc = func(c fiber.Ctx, user *domain.User, name string) (dto.OrganizationOutput, error) {
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, validInput.Name, name)
		return dto.OrganizationOutput{}, nil
	}

	app.Use(setUserInContext(app, testUser))
	app.Post("/organizations", orgHandler.CreateOrganization)

	req, err := makeJSONRequest("POST", "/organizations", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestOrganizationHandler_CreateOrganization_InvalidData(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
	}{
		{
			name:  "Empty name",
			input: map[string]interface{}{"name": ""},
		},
		{
			name:  "Name too short",
			input: map[string]interface{}{"name": "A"},
		},
		{
			name:  "Name too long",
			input: map[string]interface{}{"name": string(make([]byte, 300))},
		},
		{
			name:  "Missing name",
			input: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApp()
			mockService := NewMockOrganizationService()
			orgHandler := handler.NewOrganizationHandler(mockService)

			testUser := makeTestUser()

			app.Use(setUserInContext(app, testUser))
			app.Post("/organizations", orgHandler.CreateOrganization)

			req, err := makeJSONRequest("POST", "/organizations", tt.input)
			assert.NoError(t, err)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestOrganizationHandler_UpdateOrganization_Success(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	testUser := makeTestUser()
	orgID := uuid.New()
	validInput := dto.UpdateOrganizationInput{
		Name: "Updated Organization Name",
	}

	mockService.UpdateOrganizationFunc = func(c fiber.Ctx, user *domain.User, organizationID uuid.UUID, req dto.UpdateOrganizationInput) (dto.OrganizationOutput, error) {
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, orgID, organizationID)
		assert.Equal(t, validInput.Name, req.Name)
		return dto.OrganizationOutput{}, nil
	}

	app.Use(setUserInContext(app, testUser))
	app.Put("/organizations/:id", orgHandler.UpdateOrganization)

	req, err := makeJSONRequest("PUT", "/organizations/"+orgID.String(), validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestOrganizationHandler_UpdateOrganization_InvalidUUID(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	testUser := makeTestUser()
	validInput := dto.UpdateOrganizationInput{
		Name: "Updated Name",
	}

	app.Use(setUserInContext(app, testUser))
	app.Put("/organizations/:id", orgHandler.UpdateOrganization)

	req, err := makeJSONRequest("PUT", "/organizations/invalid-uuid", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusInternalServerError)
}

func TestOrganizationHandler_ActivateOrganization_Success(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	testUser := makeTestUser()
	orgID := uuid.New()

	mockService.ActivateOrganizationFunc = func(ctx context.Context, userID uuid.UUID, organizationID uuid.UUID) (dto.OrganizationOutput, error) {
		assert.Equal(t, testUser.ID, userID)
		assert.Equal(t, orgID, organizationID)
		return dto.OrganizationOutput{}, nil
	}

	app.Use(setUserInContext(app, testUser))
	app.Patch("/organizations/:id/activate", orgHandler.ActivateOrganization)

	req, err := makeJSONRequest("PATCH", "/organizations/"+orgID.String()+"/activate", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestOrganizationHandler_GetOrganizationByID_Success(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	testUser := makeTestUser()
	orgID := uuid.New()
	expectedOutput := dto.OrganizationOutput{
		MinimalOrganizationOutput: dto.MinimalOrganizationOutput{
			ID:       orgID.String(),
			Name:     "Test Org",
			IsActive: true,
		},
	}

	mockService.GetOrganizationByIDFunc = func(c fiber.Ctx, user *domain.User, organizationID uuid.UUID) (dto.OrganizationOutput, error) {
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, orgID, organizationID)
		return expectedOutput, nil
	}

	app.Use(setUserInContext(app, testUser))
	app.Get("/organizations/:id", orgHandler.GetOrganizationByID)

	req, err := makeJSONRequest("GET", "/organizations/"+orgID.String(), nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestOrganizationHandler_ServiceError(t *testing.T) {
	app := newTestApp()
	mockService := NewMockOrganizationService()
	orgHandler := handler.NewOrganizationHandler(mockService)

	testUser := makeTestUser()
	orgID := uuid.New()
	serviceError := errors.New("database error")

	mockService.GetOrganizationsFunc = func(c fiber.Ctx, user *domain.User, activeOrganizationID uuid.UUID) ([]dto.MinimalOrganizationOutput, error) {
		return nil, serviceError
	}

	app.Use(setUserInContext(app, testUser))
	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/organizations", orgHandler.GetOrganizations)

	req, err := makeJSONRequest("GET", "/organizations", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}
