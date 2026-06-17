package endpoint

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	endpointDTO "flowforge-api/infrastructure/endpoint"
	"flowforge-api/usecase/tag"

	"github.com/google/uuid"
)

type UpdateEndpointUseCase struct {
	endpointRepo          *repository.EndpointRepository
	getTagOrCreateUseCase *tag.GetTagOrCreateUseCase
}

func NewUpdateEndpointUseCase(
	endpointRepo *repository.EndpointRepository,
	getTagOrCreateUseCase *tag.GetTagOrCreateUseCase,
) *UpdateEndpointUseCase {
	return &UpdateEndpointUseCase{
		endpointRepo:          endpointRepo,
		getTagOrCreateUseCase: getTagOrCreateUseCase,
	}
}

func (u *UpdateEndpointUseCase) Execute(ctx context.Context, organizationID uuid.UUID, endpointID uuid.UUID, input endpointDTO.UpdateEndpointInput) (entity.Endpoint, error) {
	endpoint, err := (*u.endpointRepo).GetByIDAndOrganizationID(ctx, endpointID, organizationID)
	if err != nil {
		return entity.Endpoint{}, err
	}

	tags := []entity.Tag{}
	for _, tag := range input.Tags {
		tagEntity, err := u.getTagOrCreateUseCase.Execute(ctx, organizationID, tag.ID, tag.Name, tag.Color)
		if err == nil {
			tags = append(tags, tagEntity)
		}
	}

	endpoint.Name = input.Name
	endpoint.BaseURI = input.BaseURI
	endpoint.Path = input.Path
	endpoint.Method = input.Method
	endpoint.Timeout = input.Timeout
	endpoint.Query = input.Query
	endpoint.Header = input.Header
	endpoint.Body = input.Body
	endpoint.RetryOnFailure = input.RetryOnFailure
	endpoint.RetryCount = input.RetryCount
	endpoint.RetryDelay = input.RetryDelay
	endpoint.Tags = tags

	if err := (*u.endpointRepo).Update(ctx, &endpoint); err != nil {
		return entity.Endpoint{}, err
	}

	return endpoint, nil
}
