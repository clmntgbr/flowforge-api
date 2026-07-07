package endpoint

import (
	"flowforge-api/domain/types"
	"flowforge-api/infrastructure/paginate"
	"flowforge-api/infrastructure/tag"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"
)

type PaginateEndpointQuery struct {
	paginate.PaginateQuery
	Tags string `query:"tags"`
}

func (q PaginateEndpointQuery) TagIDs() []uuid.UUID {
	if q.Tags == "" {
		return nil
	}

	tagIDs := make([]uuid.UUID, 0)
	for _, part := range strings.Split(q.Tags, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		tagID, err := uuid.Parse(part)
		if err != nil {
			continue
		}

		tagIDs = append(tagIDs, tagID)
	}

	return tagIDs
}

type CreateEndpointInput struct {
	Name           string        `json:"name" validate:"required,min=2,max=255"`
	Description    string        `json:"description"`
	BaseURI        string        `json:"baseUri" validate:"required,url"`
	Path           string        `json:"path" validate:"required"`
	Method         string        `json:"method" validate:"required"`
	Timeout        int           `json:"timeout" validate:"required,min=1,max=60,number"`
	Query          types.Query   `json:"query" validate:"required,dive"`
	Header         types.Header  `json:"header" validate:"required,dive"`
	Body           types.Body    `json:"body" validate:"required,dive"`
	Tags           tag.TagInputs `json:"tags"`
	RetryOnFailure bool          `json:"retryOnFailure"`
	RetryCount     int           `json:"retryCount" validate:"min=0,max=10,number"`
	RetryDelay     int           `json:"retryDelay" validate:"min=0,max=600,number"`
}

type OpenAPI struct {
	OpenAPI string `json:"openapi"`
	Info    struct {
		Title   string `json:"title"`
		Version string `json:"version"`
	} `json:"info"`
	Servers []struct {
		URL string `json:"url"`
	} `json:"servers"`
	Paths map[string]map[string]struct {
		OperationID string       `json:"operationId"`
		Tags        []string     `json:"tags"`
		Summary     string       `json:"summary"`
		Description string       `json:"description"`
		Query       types.Query  `json:"query"`
		Header      types.Header `json:"header"`
		Body        types.Body   `json:"body"`
	} `json:"paths"`
}

type UpdateEndpointInput struct {
	CreateEndpointInput
}

type ImportEndpointsInput struct {
	BaseURL        string                `form:"baseUrl" validate:"required,url"`
	Timeout        int                   `form:"timeout" validate:"required,min=1,max=60,number"`
	RetryOnFailure bool                  `form:"retryOnFailure"`
	RetryCount     int                   `form:"retryCount" validate:"min=0,max=10,number"`
	RetryDelay     int                   `form:"retryDelay" validate:"min=0,max=600,number"`
	File           *multipart.FileHeader `form:"file" validate:"required"`
	Tags           tag.TagInputs         `form:"tags"`
	Query          types.Query           `form:"query"`
	Header         types.Header          `form:"header"`
	Body           types.Body            `form:"body"`
}
