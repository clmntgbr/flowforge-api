package endpoint

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	endpointDTO "flowforge-api/infrastructure/endpoint"
	"flowforge-api/usecase/tag"

	"github.com/google/uuid"
)

type CreateEndpointUseCase struct {
	endpointRepo          *repository.EndpointRepository
	getTagOrCreateUseCase *tag.GetTagOrCreateUseCase
}

func NewCreateEndpointUseCase(
	endpointRepo *repository.EndpointRepository,
	getTagOrCreateUseCase *tag.GetTagOrCreateUseCase,
) *CreateEndpointUseCase {
	return &CreateEndpointUseCase{
		endpointRepo:          endpointRepo,
		getTagOrCreateUseCase: getTagOrCreateUseCase,
	}
}

func (u *CreateEndpointUseCase) Execute(ctx context.Context, organizationID uuid.UUID, input endpointDTO.CreateEndpointInput) (entity.Endpoint, error) {

	tags := []entity.Tag{}
	for _, tag := range input.Tags {
		tagEntity, err := u.getTagOrCreateUseCase.Execute(ctx, organizationID, tag.ID, tag.Name, tag.Color)
		if err == nil {
			tags = append(tags, tagEntity)
		}
	}

	endpoint := &entity.Endpoint{
		Name:           input.Name,
		OrganizationID: organizationID,
		Description:    input.Description,
		BaseURI:        input.BaseURI,
		Path:           input.Path,
		Method:         input.Method,
		Timeout:        input.Timeout,
		Query:          input.Query,
		Header:         input.Header,
		Body:           input.Body,
		RetryOnFailure: input.RetryOnFailure,
		RetryCount:     input.RetryCount,
		RetryDelay:     input.RetryDelay,
		Tags:           tags,
	}

	if err := (*u.endpointRepo).Create(ctx, endpoint); err != nil {
		return entity.Endpoint{}, err
	}

	return *endpoint, nil
}
