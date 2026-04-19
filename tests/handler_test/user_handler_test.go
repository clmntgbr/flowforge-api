package handler_test

import (
	"errors"
	"forgeflow-api/domain"
	"forgeflow-api/dto"
	"forgeflow-api/handler"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_GetUser_Success(t *testing.T) {
	app := newTestApp()
	mockService := NewMockUserService()
	userHandler := handler.NewUserHandler(mockService)

	testUser := makeTestUser()
	expectedOutput := dto.NewUserOutput(*testUser)

	mockService.GetUserFunc = func(user *domain.User) (*dto.UserOutput, error) {
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, testUser.ClerkID, user.ClerkID)
		return &expectedOutput, nil
	}

	app.Use(setUserInContext(app, testUser))
	app.Get("/user", userHandler.GetUser)

	req, err := makeJSONRequest("GET", "/user", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result dto.UserOutput
	err = parseJSONResponse(resp, &result)
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput.ID, result.ID)
	assert.Equal(t, expectedOutput.ClerkID, result.ClerkID)
	assert.Equal(t, expectedOutput.FirstName, result.FirstName)
	assert.Equal(t, expectedOutput.LastName, result.LastName)
}

func TestUserHandler_GetUser_Unauthorized(t *testing.T) {
	app := newTestApp()
	mockService := NewMockUserService()
	userHandler := handler.NewUserHandler(mockService)

	app.Get("/user", userHandler.GetUser)

	req, err := makeJSONRequest("GET", "/user", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	var result map[string]interface{}
	err = parseJSONResponse(resp, &result)
	assert.NoError(t, err)
	assert.Contains(t, result["message"], "not authenticated")
}

func TestUserHandler_GetUser_ServiceError(t *testing.T) {
	app := newTestApp()
	mockService := NewMockUserService()
	userHandler := handler.NewUserHandler(mockService)

	testUser := makeTestUser()
	serviceError := errors.New("database connection failed")

	mockService.GetUserFunc = func(user *domain.User) (*dto.UserOutput, error) {
		return nil, serviceError
	}

	app.Use(setUserInContext(app, testUser))
	app.Get("/user", userHandler.GetUser)

	req, err := makeJSONRequest("GET", "/user", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var result map[string]interface{}
	err = parseJSONResponse(resp, &result)
	assert.NoError(t, err)
	assert.Equal(t, serviceError.Error(), result["message"])
}

func TestUserHandler_GetUser_ServiceReturnsNil(t *testing.T) {
	app := newTestApp()
	mockService := NewMockUserService()
	userHandler := handler.NewUserHandler(mockService)

	testUser := makeTestUser()

	mockService.GetUserFunc = func(user *domain.User) (*dto.UserOutput, error) {
		return nil, errors.New("user not found")
	}

	app.Use(setUserInContext(app, testUser))
	app.Get("/user", userHandler.GetUser)

	req, err := makeJSONRequest("GET", "/user", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
