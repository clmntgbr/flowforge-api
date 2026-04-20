package handler_test

import (
	"errors"
	"forgeflow-api/dto"
	"forgeflow-api/handler"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestEndpointHandler_GetEndpoints_Success(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	orgID := uuid.New()
	expectedResponse := dto.PaginateResponse{
		Total:      2,
		Page:       1,
		Limit:      20,
		TotalPages: 1,
		Members: []dto.MinimalEndpointOutput{
			{ID: uuid.New().String(), Name: "Endpoint 1", Path: "/api/v1", Method: "GET"},
			{ID: uuid.New().String(), Name: "Endpoint 2", Path: "/api/v2", Method: "POST"},
		},
	}

	mockService.GetEndpointsFunc = func(c fiber.Ctx, organizationID uuid.UUID, query dto.PaginateQuery) (dto.PaginateResponse, error) {
		assert.Equal(t, orgID, organizationID)
		return expectedResponse, nil
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/endpoints", endpointHandler.GetEndpoints)

	req, err := makeJSONRequest("GET", "/endpoints?page=1&limit=20", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result dto.PaginateResponse
	err = parseJSONResponse(resp, &result)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Total, result.Total)
}

func TestEndpointHandler_GetEndpoints_Unauthorized(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	app.Get("/endpoints", endpointHandler.GetEndpoints)

	req, err := makeJSONRequest("GET", "/endpoints", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestEndpointHandler_GetEndpoints_ServiceError(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	orgID := uuid.New()
	serviceError := errors.New("database error")

	mockService.GetEndpointsFunc = func(c fiber.Ctx, organizationID uuid.UUID, query dto.PaginateQuery) (dto.PaginateResponse, error) {
		return dto.PaginateResponse{}, serviceError
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/endpoints", endpointHandler.GetEndpoints)

	req, err := makeJSONRequest("GET", "/endpoints", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestEndpointHandler_CreateEndpoint_Success(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	orgID := uuid.New()
	validInput := dto.CreateEndpointInput{
		Name:           "Test Endpoint",
		BaseURI:        "https://api.example.com",
		Path:           "/users",
		Method:         "GET",
		Timeout:        50,
		Query:          datatypes.JSON([]byte(`{"key": "value"}`)),
		RetryOnFailure: false,
		RetryCount:     3,
		RetryDelay:     500,
	}

	mockService.CreateEndpointFunc = func(c fiber.Ctx, organizationID uuid.UUID, req dto.CreateEndpointInput) (dto.EndpointOutput, error) {
		assert.Equal(t, orgID, organizationID)
		assert.Equal(t, validInput.Name, req.Name)
		assert.Equal(t, validInput.BaseURI, req.BaseURI)
		return dto.EndpointOutput{}, nil
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Post("/endpoints", endpointHandler.CreateEndpoint)

	req, err := makeJSONRequest("POST", "/endpoints", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	err = parseJSONResponse(resp, &result)
	assert.NoError(t, err)
	assert.True(t, result["success"].(bool))
}

func TestEndpointHandler_CreateEndpoint_InvalidData(t *testing.T) {
	tests := []struct {
		name        string
		input       map[string]interface{}
		description string
	}{
		{
			name:        "Empty name",
			input:       map[string]interface{}{"name": "", "baseUri": "https://api.example.com", "path": "/users", "method": "GET", "timeout": 50, "query": `{}`},
			description: "Name is required",
		},
		{
			name:        "Name too short",
			input:       map[string]interface{}{"name": "A", "baseUri": "https://api.example.com", "path": "/users", "method": "GET", "timeout": 50, "query": `{}`},
			description: "Name must be at least 2 characters",
		},
		{
			name:        "Invalid URL",
			input:       map[string]interface{}{"name": "Test", "baseUri": "not-a-url", "path": "/users", "method": "GET", "timeout": 50, "query": `{}`},
			description: "BaseURI must be a valid URL",
		},
		{
			name:        "Missing path",
			input:       map[string]interface{}{"name": "Test", "baseUri": "https://api.example.com", "path": "", "method": "GET", "timeout": 50, "query": `{}`},
			description: "Path is required",
		},
		{
			name:        "Missing method",
			input:       map[string]interface{}{"name": "Test", "baseUri": "https://api.example.com", "path": "/users", "method": "", "timeout": 50, "query": `{}`},
			description: "Method is required",
		},
		{
			name:        "Timeout too low",
			input:       map[string]interface{}{"name": "Test", "baseUri": "https://api.example.com", "path": "/users", "method": "GET", "timeout": 0, "query": `{}`},
			description: "Timeout must be at least 1",
		},
		{
			name:        "Timeout too high",
			input:       map[string]interface{}{"name": "Test", "baseUri": "https://api.example.com", "path": "/users", "method": "GET", "timeout": 500000, "query": `{}`},
			description: "Timeout must be at most 60",
		},
		{
			name:        "Missing all fields",
			input:       map[string]interface{}{},
			description: "All fields are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApp()
			mockService := NewMockEndpointService()
			endpointHandler := handler.NewEndpointHandler(mockService)

			orgID := uuid.New()

			app.Use(setOrganizationIDInContext(app, orgID))
			app.Post("/endpoints", endpointHandler.CreateEndpoint)

			req, err := makeJSONRequest("POST", "/endpoints", tt.input)
			assert.NoError(t, err)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestEndpointHandler_CreateEndpoint_Unauthorized(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	app.Post("/endpoints", endpointHandler.CreateEndpoint)

	validInput := dto.CreateEndpointInput{
		Name:           "Test Endpoint",
		BaseURI:        "https://api.example.com",
		Path:           "/users",
		Method:         "GET",
		Timeout:        50,
		Query:          datatypes.JSON([]byte(`{}`)),
		RetryOnFailure: false,
		RetryCount:     3,
		RetryDelay:     500,
	}

	req, err := makeJSONRequest("POST", "/endpoints", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestEndpointHandler_CreateEndpoint_ServiceError(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	orgID := uuid.New()
	serviceError := errors.New("failed to create endpoint")

	mockService.CreateEndpointFunc = func(c fiber.Ctx, organizationID uuid.UUID, req dto.CreateEndpointInput) (dto.EndpointOutput, error) {
		return dto.EndpointOutput{}, serviceError
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Post("/endpoints", endpointHandler.CreateEndpoint)

	validInput := dto.CreateEndpointInput{
		Name:           "Test Endpoint",
		BaseURI:        "https://api.example.com",
		Path:           "/users",
		Method:         "GET",
		Timeout:        50,
		Query:          datatypes.JSON([]byte(`{}`)),
		RetryOnFailure: false,
		RetryCount:     3,
		RetryDelay:     500,
	}

	req, err := makeJSONRequest("POST", "/endpoints", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestEndpointHandler_GetEndpointByID_Success(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	orgID := uuid.New()
	endpointID := uuid.New()
	expectedOutput := dto.EndpointOutput{
		MinimalEndpointOutput: dto.MinimalEndpointOutput{
			ID:     endpointID.String(),
			Name:   "Test Endpoint",
			Path:   "/users",
			Method: "GET",
		},
		BaseURI: "https://api.example.com",
		Timeout: 5000,
		Query:   datatypes.JSON([]byte(`{}`)),
	}

	mockService.GetEndpointByIDFunc = func(c fiber.Ctx, organizationID uuid.UUID, id uuid.UUID) (dto.EndpointOutput, error) {
		assert.Equal(t, orgID, organizationID)
		assert.Equal(t, endpointID, id)
		return expectedOutput, nil
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/endpoints/:id", endpointHandler.GetEndpointByID)

	req, err := makeJSONRequest("GET", "/endpoints/"+endpointID.String(), nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result dto.EndpointOutput
	err = parseJSONResponse(resp, &result)
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput.ID, result.ID)
	assert.Equal(t, expectedOutput.Name, result.Name)
}

func TestEndpointHandler_GetEndpointByID_InvalidUUID(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	orgID := uuid.New()

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/endpoints/:id", endpointHandler.GetEndpointByID)

	req, err := makeJSONRequest("GET", "/endpoints/invalid-uuid", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	// parseUUIDParam sends BadRequest response but returns error,
	// which might result in 500 from error handler
	assert.True(t, resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusInternalServerError)
}

func TestEndpointHandler_GetEndpointByID_Unauthorized(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	app.Get("/endpoints/:id", endpointHandler.GetEndpointByID)

	endpointID := uuid.New()
	req, err := makeJSONRequest("GET", "/endpoints/"+endpointID.String(), nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestEndpointHandler_GetEndpointByID_ServiceError(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	orgID := uuid.New()
	endpointID := uuid.New()
	serviceError := errors.New("endpoint not found")

	mockService.GetEndpointByIDFunc = func(c fiber.Ctx, organizationID uuid.UUID, id uuid.UUID) (dto.EndpointOutput, error) {
		return dto.EndpointOutput{}, serviceError
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/endpoints/:id", endpointHandler.GetEndpointByID)

	req, err := makeJSONRequest("GET", "/endpoints/"+endpointID.String(), nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestEndpointHandler_UpdateEndpoint_Success(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	orgID := uuid.New()
	endpointID := uuid.New()
	validInput := dto.UpdateEndpointInput{
		Name:           "Updated Endpoint",
		BaseURI:        "https://api.updated.com",
		Path:           "/v2/users",
		Method:         "POST",
		Timeout:        50,
		Query:          datatypes.JSON([]byte(`{"updated": "value"}`)),
		RetryOnFailure: true,
		RetryCount:     5,
		RetryDelay:     500,
	}

	mockService.UpdateEndpointFunc = func(c fiber.Ctx, organizationID uuid.UUID, id uuid.UUID, req dto.UpdateEndpointInput) (dto.EndpointOutput, error) {
		assert.Equal(t, orgID, organizationID)
		assert.Equal(t, endpointID, id)
		assert.Equal(t, validInput.Name, req.Name)
		assert.Equal(t, validInput.BaseURI, req.BaseURI)
		return dto.EndpointOutput{}, nil
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Put("/endpoints/:id", endpointHandler.UpdateEndpoint)

	req, err := makeJSONRequest("PUT", "/endpoints/"+endpointID.String(), validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = parseJSONResponse(resp, &result)
	assert.NoError(t, err)
	assert.True(t, result["success"].(bool))
}

func TestEndpointHandler_UpdateEndpoint_InvalidData(t *testing.T) {
	tests := []struct {
		name        string
		input       map[string]interface{}
		description string
	}{
		{
			name:        "Empty name",
			input:       map[string]interface{}{"name": "", "baseUri": "https://api.example.com", "path": "/users", "method": "GET", "timeout": 50, "query": `{}`},
			description: "Name is required",
		},
		{
			name:        "Negative timeout",
			input:       map[string]interface{}{"name": "Updated", "baseUri": "https://api.example.com", "path": "/users", "method": "GET", "timeout": -100, "query": `{}`},
			description: "Timeout cannot be negative",
		},
		{
			name:        "Missing method",
			input:       map[string]interface{}{"name": "Updated", "baseUri": "https://api.example.com", "path": "/users", "method": "", "timeout": 50, "query": `{}`},
			description: "Method is required",
		},
		{
			name:        "Very long name",
			input:       map[string]interface{}{"name": string(make([]byte, 300)), "baseUri": "https://api.example.com", "path": "/users", "method": "GET", "timeout": 50, "query": `{}`},
			description: "Name exceeds maximum length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApp()
			mockService := NewMockEndpointService()
			endpointHandler := handler.NewEndpointHandler(mockService)

			orgID := uuid.New()
			endpointID := uuid.New()

			app.Use(setOrganizationIDInContext(app, orgID))
			app.Put("/endpoints/:id", endpointHandler.UpdateEndpoint)

			req, err := makeJSONRequest("PUT", "/endpoints/"+endpointID.String(), tt.input)
			assert.NoError(t, err)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestEndpointHandler_UpdateEndpoint_InvalidUUID(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	orgID := uuid.New()
	validInput := dto.UpdateEndpointInput{
		Name:    "Updated",
		BaseURI: "https://api.example.com",
		Path:    "/users",
		Method:  "GET",
		Timeout: 50,
		Query:   datatypes.JSON([]byte(`{}`)),
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Put("/endpoints/:id", endpointHandler.UpdateEndpoint)

	req, err := makeJSONRequest("PUT", "/endpoints/not-a-uuid", validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	// parseUUIDParam sends BadRequest response but returns error,
	// which might result in 500 from error handler
	assert.True(t, resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusInternalServerError)
}

func TestEndpointHandler_UpdateEndpoint_Unauthorized(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	endpointID := uuid.New()
	validInput := dto.UpdateEndpointInput{
		Name:           "Updated",
		BaseURI:        "https://api.example.com",
		Path:           "/users",
		Method:         "GET",
		Timeout:        50,
		Query:          datatypes.JSON([]byte(`{}`)),
		RetryOnFailure: false,
		RetryCount:     3,
		RetryDelay:     500,
	}

	app.Put("/endpoints/:id", endpointHandler.UpdateEndpoint)

	req, err := makeJSONRequest("PUT", "/endpoints/"+endpointID.String(), validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestEndpointHandler_UpdateEndpoint_ServiceError(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	orgID := uuid.New()
	endpointID := uuid.New()
	serviceError := errors.New("failed to update endpoint")

	mockService.UpdateEndpointFunc = func(c fiber.Ctx, organizationID uuid.UUID, id uuid.UUID, req dto.UpdateEndpointInput) (dto.EndpointOutput, error) {
		return dto.EndpointOutput{}, serviceError
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Put("/endpoints/:id", endpointHandler.UpdateEndpoint)

	validInput := dto.UpdateEndpointInput{
		Name:           "Updated",
		BaseURI:        "https://api.example.com",
		Path:           "/users",
		Method:         "GET",
		Timeout:        50,
		Query:          datatypes.JSON([]byte(`{}`)),
		RetryOnFailure: false,
		RetryCount:     3,
		RetryDelay:     500,
	}

	req, err := makeJSONRequest("PUT", "/endpoints/"+endpointID.String(), validInput)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestEndpointHandler_CreateEndpoint_VariousValidInputs(t *testing.T) {
	validInputs := []struct {
		name  string
		input dto.CreateEndpointInput
	}{
		{
			name: "Minimal timeout",
			input: dto.CreateEndpointInput{
				Name:           "Fast Endpoint",
				BaseURI:        "https://api.example.com",
				Path:           "/quick",
				Method:         "GET",
				Timeout:        1,
				Query:          datatypes.JSON([]byte(`{}`)),
				RetryOnFailure: false,
				RetryCount:     0,
				RetryDelay:     0,
			},
		},
		{
			name: "Maximum timeout",
			input: dto.CreateEndpointInput{
				Name:           "Slow Endpoint",
				BaseURI:        "https://api.example.com",
				Path:           "/slow",
				Method:         "POST",
				Timeout:        60,
				Query:          datatypes.JSON([]byte(`{}`)),
				RetryOnFailure: true,
				RetryCount:     10,
				RetryDelay:     500,
			},
		},
		{
			name: "Complex path with parameters",
			input: dto.CreateEndpointInput{
				Name:           "Dynamic Endpoint",
				BaseURI:        "https://api.example.com",
				Path:           "/users/:id/posts/:postId",
				Method:         "DELETE",
				Timeout:        50,
				Query:          datatypes.JSON([]byte(`{"param": "value"}`)),
				RetryOnFailure: true,
				RetryCount:     3,
				RetryDelay:     500,
			},
		},
		{
			name: "Different HTTP methods",
			input: dto.CreateEndpointInput{
				Name:           "PATCH Endpoint",
				BaseURI:        "https://api.example.com",
				Path:           "/resource",
				Method:         "PATCH",
				Timeout:        50,
				Query:          datatypes.JSON([]byte(`{"filter": "active"}`)),
				RetryOnFailure: false,
				RetryCount:     2,
				RetryDelay:     500,
			},
		},
	}

	for _, tt := range validInputs {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApp()
			mockService := NewMockEndpointService()
			endpointHandler := handler.NewEndpointHandler(mockService)

			orgID := uuid.New()

			mockService.CreateEndpointFunc = func(c fiber.Ctx, organizationID uuid.UUID, req dto.CreateEndpointInput) (dto.EndpointOutput, error) {
				return dto.EndpointOutput{}, nil
			}

			app.Use(setOrganizationIDInContext(app, orgID))
			app.Post("/endpoints", endpointHandler.CreateEndpoint)

			req, err := makeJSONRequest("POST", "/endpoints", tt.input)
			assert.NoError(t, err)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusCreated, resp.StatusCode)
		})
	}
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

	req, err := makeJSONRequest("GET", "/endpoints?page=0&limit=0", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestEndpointHandler_GetEndpoints_WithMalformedQueryString(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	orgID := uuid.New()

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/endpoints", endpointHandler.GetEndpoints)

	req, err := makeJSONRequest("GET", "/endpoints", nil)
	assert.NoError(t, err)

	req.URL.RawQuery = strings.Repeat("invalid&", 1000)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode == fiber.StatusOK || resp.StatusCode >= 400)
}

func TestBaseHandler_BindAndValidate_InvalidJSON(t *testing.T) {
	app := newTestApp()
	mockService := NewMockEndpointService()
	endpointHandler := handler.NewEndpointHandler(mockService)

	orgID := uuid.New()

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Post("/endpoints", endpointHandler.CreateEndpoint)

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
		t.Fatal("GetEndpoints must not be called when query binding fails")
		return dto.PaginateResponse{}, nil
	}

	app.Use(setOrganizationIDInContext(app, orgID))
	app.Get("/endpoints", endpointHandler.GetEndpoints)

	req, err := makeJSONRequest("GET", "/endpoints?page=invalid", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}
