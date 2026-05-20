package endpoint

import (
	"context"
	"encoding/json"
	"flowforge-api/domain/repository"
	"flowforge-api/infrastructure/endpoint"
	"io"
	"log"
	"mime/multipart"

	"github.com/google/uuid"
)

type ImportFromOpenAPIUseCase struct {
	endpointRepo *repository.EndpointRepository
}

func NewImportFromOpenAPIUseCase(endpointRepo *repository.EndpointRepository) *ImportFromOpenAPIUseCase {
	return &ImportFromOpenAPIUseCase{endpointRepo: endpointRepo}
}

func (u *ImportFromOpenAPIUseCase) Execute(ctx context.Context, organizationID uuid.UUID, baseURL string, file *multipart.FileHeader) error {
	f, _ := file.Open()
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
