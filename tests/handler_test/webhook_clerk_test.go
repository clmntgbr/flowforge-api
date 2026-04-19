package handler_test

import (
	"encoding/json"
	"errors"
	"forgeflow-api/domain"
	"forgeflow-api/dto"
	"forgeflow-api/handler"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setClerkEventInContext(app *fiber.App, event dto.ClerkEvent) fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Locals("payload", event)
		return c.Next()
	}
}

func TestWebhookClerkHandler_Handle_UserCreated_Success(t *testing.T) {
	app := newTestApp()
	mockUserService := NewMockUserService()
	mockOrgService := NewMockOrganizationService()
	mockUserRepo := NewMockUserRepository()
	webhookHandler := handler.NewWebhookClerkHandler(mockUserService, mockOrgService, mockUserRepo)

	banned := false
	userData := dto.ClerkUserCreated{
		ID:        "clerk_123",
		FirstName: "John",
		LastName:  "Doe",
		Banned:    &banned,
	}
	dataJSON, _ := json.Marshal(userData)

	clerkEvent := dto.ClerkEvent{
		Type:       "user.created",
		InstanceID: "ins_123",
		Object:     "event",
		Timestamp:  1234567890,
		Data:       dataJSON,
	}

	newUser := &domain.User{
		ID:        uuid.New(),
		ClerkID:   userData.ID,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
	}

	newOrg := dto.OrganizationOutput{
		MinimalOrganizationOutput: dto.MinimalOrganizationOutput{
			ID:   uuid.New().String(),
			Name: "Default Organization",
		},
	}

	mockUserRepo.FindByClerkIDFunc = func(clerkID string) (*domain.User, error) {
		assert.Equal(t, userData.ID, clerkID)
		return nil, nil
	}

	mockUserService.CreateUserFunc = func(c fiber.Ctx, id string, firstName string, lastName string, banned bool) (*domain.User, error) {
		assert.Equal(t, userData.ID, id)
		assert.Equal(t, userData.FirstName, firstName)
		return newUser, nil
	}

	mockOrgService.CreateOrganizationFunc = func(c fiber.Ctx, user *domain.User, name string) (dto.OrganizationOutput, error) {
		assert.Equal(t, newUser.ID, user.ID)
		assert.Equal(t, "Default Organization", name)
		return newOrg, nil
	}

	mockUserRepo.UpdateFunc = func(user *domain.User) error {
		assert.NotNil(t, user.ActiveOrganizationID)
		return nil
	}

	app.Use(setClerkEventInContext(app, clerkEvent))
	app.Post("/webhook", webhookHandler.Handle)

	req, err := makeJSONRequest("POST", "/webhook", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestWebhookClerkHandler_Handle_UserCreated_UserAlreadyExists(t *testing.T) {
	app := newTestApp()
	mockUserService := NewMockUserService()
	mockOrgService := NewMockOrganizationService()
	mockUserRepo := NewMockUserRepository()
	webhookHandler := handler.NewWebhookClerkHandler(mockUserService, mockOrgService, mockUserRepo)

	banned := false
	userData := dto.ClerkUserCreated{
		ID:        "clerk_123",
		FirstName: "John",
		LastName:  "Doe",
		Banned:    &banned,
	}
	dataJSON, _ := json.Marshal(userData)

	clerkEvent := dto.ClerkEvent{
		Type:       "user.created",
		InstanceID: "ins_123",
		Object:     "event",
		Timestamp:  1234567890,
		Data:       dataJSON,
	}

	existingUser := makeTestUser()
	existingUser.ClerkID = userData.ID

	mockUserRepo.FindByClerkIDFunc = func(clerkID string) (*domain.User, error) {
		return existingUser, nil
	}

	app.Use(setClerkEventInContext(app, clerkEvent))
	app.Post("/webhook", webhookHandler.Handle)

	req, err := makeJSONRequest("POST", "/webhook", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestWebhookClerkHandler_Handle_UserCreated_InvalidJSON(t *testing.T) {
	app := newTestApp()
	mockUserService := NewMockUserService()
	mockOrgService := NewMockOrganizationService()
	mockUserRepo := NewMockUserRepository()
	webhookHandler := handler.NewWebhookClerkHandler(mockUserService, mockOrgService, mockUserRepo)

	clerkEvent := dto.ClerkEvent{
		Type:       "user.created",
		InstanceID: "ins_123",
		Object:     "event",
		Timestamp:  1234567890,
		Data:       json.RawMessage(`{invalid json}`),
	}

	app.Use(setClerkEventInContext(app, clerkEvent))
	app.Post("/webhook", webhookHandler.Handle)

	req, err := makeJSONRequest("POST", "/webhook", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestWebhookClerkHandler_Handle_UserCreated_ValidationError(t *testing.T) {
	app := newTestApp()
	mockUserService := NewMockUserService()
	mockOrgService := NewMockOrganizationService()
	mockUserRepo := NewMockUserRepository()
	webhookHandler := handler.NewWebhookClerkHandler(mockUserService, mockOrgService, mockUserRepo)

	userData := dto.ClerkUserCreated{
		ID:        "",
		FirstName: "John",
		LastName:  "Doe",
		Banned:    nil,
	}
	dataJSON, _ := json.Marshal(userData)

	clerkEvent := dto.ClerkEvent{
		Type:       "user.created",
		InstanceID: "ins_123",
		Object:     "event",
		Timestamp:  1234567890,
		Data:       dataJSON,
	}

	app.Use(setClerkEventInContext(app, clerkEvent))
	app.Post("/webhook", webhookHandler.Handle)

	req, err := makeJSONRequest("POST", "/webhook", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusInternalServerError)
}

func TestWebhookClerkHandler_Handle_UserCreated_CreateUserError(t *testing.T) {
	app := newTestApp()
	mockUserService := NewMockUserService()
	mockOrgService := NewMockOrganizationService()
	mockUserRepo := NewMockUserRepository()
	webhookHandler := handler.NewWebhookClerkHandler(mockUserService, mockOrgService, mockUserRepo)

	banned := false
	userData := dto.ClerkUserCreated{
		ID:        "clerk_123",
		FirstName: "John",
		LastName:  "Doe",
		Banned:    &banned,
	}
	dataJSON, _ := json.Marshal(userData)

	clerkEvent := dto.ClerkEvent{
		Type:       "user.created",
		InstanceID: "ins_123",
		Object:     "event",
		Timestamp:  1234567890,
		Data:       dataJSON,
	}

	mockUserRepo.FindByClerkIDFunc = func(clerkID string) (*domain.User, error) {
		return nil, nil
	}

	mockUserService.CreateUserFunc = func(c fiber.Ctx, id string, firstName string, lastName string, banned bool) (*domain.User, error) {
		return nil, errors.New("failed to create user")
	}

	app.Use(setClerkEventInContext(app, clerkEvent))
	app.Post("/webhook", webhookHandler.Handle)

	req, err := makeJSONRequest("POST", "/webhook", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestWebhookClerkHandler_Handle_UserCreated_CreateOrgError(t *testing.T) {
	app := newTestApp()
	mockUserService := NewMockUserService()
	mockOrgService := NewMockOrganizationService()
	mockUserRepo := NewMockUserRepository()
	webhookHandler := handler.NewWebhookClerkHandler(mockUserService, mockOrgService, mockUserRepo)

	banned := false
	userData := dto.ClerkUserCreated{
		ID:        "clerk_123",
		FirstName: "John",
		LastName:  "Doe",
		Banned:    &banned,
	}
	dataJSON, _ := json.Marshal(userData)

	clerkEvent := dto.ClerkEvent{
		Type:       "user.created",
		InstanceID: "ins_123",
		Object:     "event",
		Timestamp:  1234567890,
		Data:       dataJSON,
	}

	newUser := makeTestUser()

	mockUserRepo.FindByClerkIDFunc = func(clerkID string) (*domain.User, error) {
		return nil, nil
	}

	mockUserService.CreateUserFunc = func(c fiber.Ctx, id string, firstName string, lastName string, banned bool) (*domain.User, error) {
		return newUser, nil
	}

	mockOrgService.CreateOrganizationFunc = func(c fiber.Ctx, user *domain.User, name string) (dto.OrganizationOutput, error) {
		return dto.OrganizationOutput{}, errors.New("failed to create organization")
	}

	app.Use(setClerkEventInContext(app, clerkEvent))
	app.Post("/webhook", webhookHandler.Handle)

	req, err := makeJSONRequest("POST", "/webhook", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestWebhookClerkHandler_Handle_UserCreated_UpdateUserError(t *testing.T) {
	app := newTestApp()
	mockUserService := NewMockUserService()
	mockOrgService := NewMockOrganizationService()
	mockUserRepo := NewMockUserRepository()
	webhookHandler := handler.NewWebhookClerkHandler(mockUserService, mockOrgService, mockUserRepo)

	banned := false
	userData := dto.ClerkUserCreated{
		ID:        "clerk_123",
		FirstName: "John",
		LastName:  "Doe",
		Banned:    &banned,
	}
	dataJSON, _ := json.Marshal(userData)

	clerkEvent := dto.ClerkEvent{
		Type:       "user.created",
		InstanceID: "ins_123",
		Object:     "event",
		Timestamp:  1234567890,
		Data:       dataJSON,
	}

	newUser := makeTestUser()
	newOrg := dto.OrganizationOutput{
		MinimalOrganizationOutput: dto.MinimalOrganizationOutput{
			ID:   uuid.New().String(),
			Name: "Default Organization",
		},
	}

	mockUserRepo.FindByClerkIDFunc = func(clerkID string) (*domain.User, error) {
		return nil, nil
	}

	mockUserService.CreateUserFunc = func(c fiber.Ctx, id string, firstName string, lastName string, banned bool) (*domain.User, error) {
		return newUser, nil
	}

	mockOrgService.CreateOrganizationFunc = func(c fiber.Ctx, user *domain.User, name string) (dto.OrganizationOutput, error) {
		return newOrg, nil
	}

	mockUserRepo.UpdateFunc = func(user *domain.User) error {
		return errors.New("failed to update user")
	}

	app.Use(setClerkEventInContext(app, clerkEvent))
	app.Post("/webhook", webhookHandler.Handle)

	req, err := makeJSONRequest("POST", "/webhook", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestWebhookClerkHandler_Handle_UserCreated_InvalidOrgUUID(t *testing.T) {
	app := newTestApp()
	mockUserService := NewMockUserService()
	mockOrgService := NewMockOrganizationService()
	mockUserRepo := NewMockUserRepository()
	webhookHandler := handler.NewWebhookClerkHandler(mockUserService, mockOrgService, mockUserRepo)

	banned := false
	userData := dto.ClerkUserCreated{
		ID:        "clerk_123",
		FirstName: "John",
		LastName:  "Doe",
		Banned:    &banned,
	}
	dataJSON, _ := json.Marshal(userData)

	clerkEvent := dto.ClerkEvent{
		Type:       "user.created",
		InstanceID: "ins_123",
		Object:     "event",
		Timestamp:  1234567890,
		Data:       dataJSON,
	}

	newUser := makeTestUser()
	newOrg := dto.OrganizationOutput{
		MinimalOrganizationOutput: dto.MinimalOrganizationOutput{
			ID:   "invalid-uuid",
			Name: "Default Organization",
		},
	}

	mockUserRepo.FindByClerkIDFunc = func(clerkID string) (*domain.User, error) {
		return nil, nil
	}

	mockUserService.CreateUserFunc = func(c fiber.Ctx, id string, firstName string, lastName string, banned bool) (*domain.User, error) {
		return newUser, nil
	}

	mockOrgService.CreateOrganizationFunc = func(c fiber.Ctx, user *domain.User, name string) (dto.OrganizationOutput, error) {
		return newOrg, nil
	}

	app.Use(setClerkEventInContext(app, clerkEvent))
	app.Post("/webhook", webhookHandler.Handle)

	req, err := makeJSONRequest("POST", "/webhook", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestWebhookClerkHandler_Handle_UserCreated_FindUserError(t *testing.T) {
	app := newTestApp()
	mockUserService := NewMockUserService()
	mockOrgService := NewMockOrganizationService()
	mockUserRepo := NewMockUserRepository()
	webhookHandler := handler.NewWebhookClerkHandler(mockUserService, mockOrgService, mockUserRepo)

	banned := false
	userData := dto.ClerkUserCreated{
		ID:        "clerk_123",
		FirstName: "John",
		LastName:  "Doe",
		Banned:    &banned,
	}
	dataJSON, _ := json.Marshal(userData)

	clerkEvent := dto.ClerkEvent{
		Type:       "user.created",
		InstanceID: "ins_123",
		Object:     "event",
		Timestamp:  1234567890,
		Data:       dataJSON,
	}

	mockUserRepo.FindByClerkIDFunc = func(clerkID string) (*domain.User, error) {
		return nil, errors.New("database error")
	}

	app.Use(setClerkEventInContext(app, clerkEvent))
	app.Post("/webhook", webhookHandler.Handle)

	req, err := makeJSONRequest("POST", "/webhook", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestWebhookClerkHandler_Handle_UserUpdated_Success(t *testing.T) {
	app := newTestApp()
	mockUserService := NewMockUserService()
	mockOrgService := NewMockOrganizationService()
	mockUserRepo := NewMockUserRepository()
	webhookHandler := handler.NewWebhookClerkHandler(mockUserService, mockOrgService, mockUserRepo)

	banned := true
	userData := dto.ClerkUserUpdated{
		ID:        "clerk_123",
		FirstName: "Jane",
		LastName:  "Smith",
		Banned:    &banned,
	}
	dataJSON, _ := json.Marshal(userData)

	clerkEvent := dto.ClerkEvent{
		Type:       "user.updated",
		InstanceID: "ins_123",
		Object:     "event",
		Timestamp:  1234567890,
		Data:       dataJSON,
	}

	existingUser := makeTestUser()
	existingUser.ClerkID = userData.ID

	mockUserRepo.FindByClerkIDFunc = func(clerkID string) (*domain.User, error) {
		return existingUser, nil
	}

	mockUserService.UpdateUserFunc = func(c fiber.Ctx, id string, firstName string, lastName string, banned bool) error {
		assert.Equal(t, userData.ID, id)
		assert.Equal(t, userData.FirstName, firstName)
		return nil
	}

	app.Use(setClerkEventInContext(app, clerkEvent))
	app.Post("/webhook", webhookHandler.Handle)

	req, err := makeJSONRequest("POST", "/webhook", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}

func TestWebhookClerkHandler_Handle_UserUpdated_InvalidJSON(t *testing.T) {
	app := newTestApp()
	mockUserService := NewMockUserService()
	mockOrgService := NewMockOrganizationService()
	mockUserRepo := NewMockUserRepository()
	webhookHandler := handler.NewWebhookClerkHandler(mockUserService, mockOrgService, mockUserRepo)

	clerkEvent := dto.ClerkEvent{
		Type:       "user.updated",
		InstanceID: "ins_123",
		Object:     "event",
		Timestamp:  1234567890,
		Data:       json.RawMessage(`{invalid json}`),
	}

	app.Use(setClerkEventInContext(app, clerkEvent))
	app.Post("/webhook", webhookHandler.Handle)

	req, err := makeJSONRequest("POST", "/webhook", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestWebhookClerkHandler_Handle_UserUpdated_ValidationError(t *testing.T) {
	app := newTestApp()
	mockUserService := NewMockUserService()
	mockOrgService := NewMockOrganizationService()
	mockUserRepo := NewMockUserRepository()
	webhookHandler := handler.NewWebhookClerkHandler(mockUserService, mockOrgService, mockUserRepo)

	userData := dto.ClerkUserUpdated{
		ID:        "",
		FirstName: "Jane",
		LastName:  "Smith",
		Banned:    nil,
	}
	dataJSON, _ := json.Marshal(userData)

	clerkEvent := dto.ClerkEvent{
		Type:       "user.updated",
		InstanceID: "ins_123",
		Object:     "event",
		Timestamp:  1234567890,
		Data:       dataJSON,
	}

	app.Use(setClerkEventInContext(app, clerkEvent))
	app.Post("/webhook", webhookHandler.Handle)

	req, err := makeJSONRequest("POST", "/webhook", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusInternalServerError)
}

func TestWebhookClerkHandler_Handle_UserUpdated_UserNotFound(t *testing.T) {
	app := newTestApp()
	mockUserService := NewMockUserService()
	mockOrgService := NewMockOrganizationService()
	mockUserRepo := NewMockUserRepository()
	webhookHandler := handler.NewWebhookClerkHandler(mockUserService, mockOrgService, mockUserRepo)

	banned := false
	userData := dto.ClerkUserUpdated{
		ID:        "clerk_123",
		FirstName: "Jane",
		LastName:  "Smith",
		Banned:    &banned,
	}
	dataJSON, _ := json.Marshal(userData)

	clerkEvent := dto.ClerkEvent{
		Type:       "user.updated",
		InstanceID: "ins_123",
		Object:     "event",
		Timestamp:  1234567890,
		Data:       dataJSON,
	}

	mockUserRepo.FindByClerkIDFunc = func(clerkID string) (*domain.User, error) {
		return nil, nil
	}

	app.Use(setClerkEventInContext(app, clerkEvent))
	app.Post("/webhook", webhookHandler.Handle)

	req, err := makeJSONRequest("POST", "/webhook", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestWebhookClerkHandler_Handle_UserUpdated_FindUserError(t *testing.T) {
	app := newTestApp()
	mockUserService := NewMockUserService()
	mockOrgService := NewMockOrganizationService()
	mockUserRepo := NewMockUserRepository()
	webhookHandler := handler.NewWebhookClerkHandler(mockUserService, mockOrgService, mockUserRepo)

	banned := false
	userData := dto.ClerkUserUpdated{
		ID:        "clerk_123",
		FirstName: "Jane",
		LastName:  "Smith",
		Banned:    &banned,
	}
	dataJSON, _ := json.Marshal(userData)

	clerkEvent := dto.ClerkEvent{
		Type:       "user.updated",
		InstanceID: "ins_123",
		Object:     "event",
		Timestamp:  1234567890,
		Data:       dataJSON,
	}

	mockUserRepo.FindByClerkIDFunc = func(clerkID string) (*domain.User, error) {
		return nil, errors.New("database error")
	}

	app.Use(setClerkEventInContext(app, clerkEvent))
	app.Post("/webhook", webhookHandler.Handle)

	req, err := makeJSONRequest("POST", "/webhook", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestWebhookClerkHandler_Handle_UserUpdated_UpdateError(t *testing.T) {
	app := newTestApp()
	mockUserService := NewMockUserService()
	mockOrgService := NewMockOrganizationService()
	mockUserRepo := NewMockUserRepository()
	webhookHandler := handler.NewWebhookClerkHandler(mockUserService, mockOrgService, mockUserRepo)

	banned := false
	userData := dto.ClerkUserUpdated{
		ID:        "clerk_123",
		FirstName: "Jane",
		LastName:  "Smith",
		Banned:    &banned,
	}
	dataJSON, _ := json.Marshal(userData)

	clerkEvent := dto.ClerkEvent{
		Type:       "user.updated",
		InstanceID: "ins_123",
		Object:     "event",
		Timestamp:  1234567890,
		Data:       dataJSON,
	}

	existingUser := makeTestUser()

	mockUserRepo.FindByClerkIDFunc = func(clerkID string) (*domain.User, error) {
		return existingUser, nil
	}

	mockUserService.UpdateUserFunc = func(c fiber.Ctx, id string, firstName string, lastName string, banned bool) error {
		return errors.New("failed to update user")
	}

	app.Use(setClerkEventInContext(app, clerkEvent))
	app.Post("/webhook", webhookHandler.Handle)

	req, err := makeJSONRequest("POST", "/webhook", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestWebhookClerkHandler_Handle_UserDeleted_Success(t *testing.T) {
	app := newTestApp()
	mockUserService := NewMockUserService()
	mockOrgService := NewMockOrganizationService()
	mockUserRepo := NewMockUserRepository()
	webhookHandler := handler.NewWebhookClerkHandler(mockUserService, mockOrgService, mockUserRepo)

	userData := dto.ClerkUserDeleted{
		ID: "clerk_123",
	}
	dataJSON, _ := json.Marshal(userData)

	clerkEvent := dto.ClerkEvent{
		Type:       "user.deleted",
		InstanceID: "ins_123",
		Object:     "event",
		Timestamp:  1234567890,
		Data:       dataJSON,
	}

	mockUserService.DeleteUserFunc = func(c fiber.Ctx, id string) error {
		assert.Equal(t, userData.ID, id)
		return nil
	}

	app.Use(setClerkEventInContext(app, clerkEvent))
	app.Post("/webhook", webhookHandler.Handle)

	req, err := makeJSONRequest("POST", "/webhook", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}

func TestWebhookClerkHandler_Handle_UserDeleted_InvalidJSON(t *testing.T) {
	app := newTestApp()
	mockUserService := NewMockUserService()
	mockOrgService := NewMockOrganizationService()
	mockUserRepo := NewMockUserRepository()
	webhookHandler := handler.NewWebhookClerkHandler(mockUserService, mockOrgService, mockUserRepo)

	clerkEvent := dto.ClerkEvent{
		Type:       "user.deleted",
		InstanceID: "ins_123",
		Object:     "event",
		Timestamp:  1234567890,
		Data:       json.RawMessage(`{invalid json}`),
	}

	app.Use(setClerkEventInContext(app, clerkEvent))
	app.Post("/webhook", webhookHandler.Handle)

	req, err := makeJSONRequest("POST", "/webhook", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestWebhookClerkHandler_Handle_UserDeleted_ValidationError(t *testing.T) {
	app := newTestApp()
	mockUserService := NewMockUserService()
	mockOrgService := NewMockOrganizationService()
	mockUserRepo := NewMockUserRepository()
	webhookHandler := handler.NewWebhookClerkHandler(mockUserService, mockOrgService, mockUserRepo)

	userData := dto.ClerkUserDeleted{
		ID: "",
	}
	dataJSON, _ := json.Marshal(userData)

	clerkEvent := dto.ClerkEvent{
		Type:       "user.deleted",
		InstanceID: "ins_123",
		Object:     "event",
		Timestamp:  1234567890,
		Data:       dataJSON,
	}

	app.Use(setClerkEventInContext(app, clerkEvent))
	app.Post("/webhook", webhookHandler.Handle)

	req, err := makeJSONRequest("POST", "/webhook", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusInternalServerError)
}

func TestWebhookClerkHandler_Handle_UserDeleted_DeleteError(t *testing.T) {
	app := newTestApp()
	mockUserService := NewMockUserService()
	mockOrgService := NewMockOrganizationService()
	mockUserRepo := NewMockUserRepository()
	webhookHandler := handler.NewWebhookClerkHandler(mockUserService, mockOrgService, mockUserRepo)

	userData := dto.ClerkUserDeleted{
		ID: "clerk_123",
	}
	dataJSON, _ := json.Marshal(userData)

	clerkEvent := dto.ClerkEvent{
		Type:       "user.deleted",
		InstanceID: "ins_123",
		Object:     "event",
		Timestamp:  1234567890,
		Data:       dataJSON,
	}

	mockUserService.DeleteUserFunc = func(c fiber.Ctx, id string) error {
		return errors.New("failed to delete user")
	}

	app.Use(setClerkEventInContext(app, clerkEvent))
	app.Post("/webhook", webhookHandler.Handle)

	req, err := makeJSONRequest("POST", "/webhook", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestWebhookClerkHandler_Handle_UnknownEventType(t *testing.T) {
	app := newTestApp()
	mockUserService := NewMockUserService()
	mockOrgService := NewMockOrganizationService()
	mockUserRepo := NewMockUserRepository()
	webhookHandler := handler.NewWebhookClerkHandler(mockUserService, mockOrgService, mockUserRepo)

	clerkEvent := dto.ClerkEvent{
		Type:       "user.unknown",
		InstanceID: "ins_123",
		Object:     "event",
		Timestamp:  1234567890,
		Data:       json.RawMessage(`{}`),
	}

	app.Use(setClerkEventInContext(app, clerkEvent))
	app.Post("/webhook", webhookHandler.Handle)

	req, err := makeJSONRequest("POST", "/webhook", nil)
	assert.NoError(t, err)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
