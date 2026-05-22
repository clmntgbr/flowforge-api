package endpoint

import (
	"context"
	"encoding/json"
	"errors"
	"flowforge-api/domain/repository"
	"flowforge-api/domain/types"
	"flowforge-api/infrastructure/endpoint"
	endpointDTO "flowforge-api/infrastructure/endpoint"
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

var openAPIHTTPMethods = map[string]struct{}{
	"get": {}, "put": {}, "post": {}, "delete": {},
	"options": {}, "head": {}, "patch": {}, "trace": {},
}

type ImportFromOpenAPIUseCase struct {
	endpointRepo          *repository.EndpointRepository
	createEndpointUseCase *CreateEndpointUseCase
}

func NewImportFromOpenAPIUseCase(
	endpointRepo *repository.EndpointRepository,
	createEndpointUseCase *CreateEndpointUseCase,
) *ImportFromOpenAPIUseCase {
	return &ImportFromOpenAPIUseCase{
		endpointRepo:          endpointRepo,
		createEndpointUseCase: createEndpointUseCase,
	}
}

func (u *ImportFromOpenAPIUseCase) Execute(ctx context.Context, organizationID uuid.UUID, input endpointDTO.ImportEndpointsInput) error {
	parsedURL, err := url.Parse(input.BaseURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return errors.New("baseUrl must be a valid URL (e.g., https://api.example.com)")
	}

	baseURI := strings.TrimSuffix(parsedURL.String(), "/")

	f, _ := input.File.Open()
	defer f.Close()

	data, _ := io.ReadAll(f)

	log.Println(string(data))

	var openAPI endpoint.OpenAPI
	json.Unmarshal(data, &openAPI)

	for path, methods := range openAPI.Paths {
		for method, op := range methods {
			if _, ok := openAPIHTTPMethods[strings.ToLower(method)]; !ok {
				continue
			}

			u.createEndpointUseCase.Execute(ctx, organizationID, endpointDTO.CreateEndpointInput{
				Name:           op.OperationID,
				Description:    op.Description,
				BaseURI:        baseURI,
				Path:           path,
				Method:         strings.ToUpper(method),
				Timeout:        input.Timeout,
				RetryOnFailure: input.RetryOnFailure,
				RetryCount:     input.RetryCount,
				RetryDelay:     input.RetryDelay,
				Query:          types.Query{},
				Header:         types.Header{},
				Body:           types.Body{},
			})
		}
	}

	return nil
}
