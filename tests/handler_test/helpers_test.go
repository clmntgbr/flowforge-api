package handler_test

import (
	"bytes"
	"encoding/json"
	"forgeflow-api/ctxutil"
	"forgeflow-api/domain"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func newTestApp() *fiber.App {
	return fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		},
	})
}

func makeTestUser() *domain.User {
	return &domain.User{
		ID:        uuid.New(),
		ClerkID:   "clerk_test_123",
		FirstName: "John",
		LastName:  "Doe",
		Banned:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func setUserInContext(app *fiber.App, user *domain.User) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctxutil.SetUser(c, *user)
		return c.Next()
	}
}

func setOrganizationIDInContext(app *fiber.App, organizationID uuid.UUID) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctxutil.SetOrganizationID(c, organizationID)
		return c.Next()
	}
}

func makeJSONRequest(method, url string, body interface{}) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req := httptest.NewRequest(method, url, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func parseJSONResponse(resp *http.Response, target interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}
