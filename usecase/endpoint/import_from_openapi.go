package endpoint

import (
	"context"
	"encoding/json"
	"errors"
	"flowforge-api/domain/repository"
	"flowforge-api/infrastructure/endpoint"
	endpointDTO "flowforge-api/infrastructure/endpoint"
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

type ImportFromOpenAPIUseCase struct {
	endpointRepo *repository.EndpointRepository
}

func NewImportFromOpenAPIUseCase(endpointRepo *repository.EndpointRepository) *ImportFromOpenAPIUseCase {
	return &ImportFromOpenAPIUseCase{endpointRepo: endpointRepo}
}

func (u *ImportFromOpenAPIUseCase) Execute(ctx context.Context, organizationID uuid.UUID, input endpointDTO.ImportEndpointsInput) error {
	parsedURL, err := url.Parse(input.BaseURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return errors.New("baseUrl must be a valid URL (e.g., https://api.example.com)")
	}

	_ = strings.TrimSuffix(parsedURL.String(), "/")

	f, _ := input.File.Open()
	defer f.Close()

	data, _ := io.ReadAll(f)

	var openAPI endpoint.OpenAPI
	json.Unmarshal(data, &openAPI)

	for path, methods := range openAPI.Paths {
		for method, op := range methods {
			log.Println(path, method, op.OperationID)
		}
	}

	return nil
}
