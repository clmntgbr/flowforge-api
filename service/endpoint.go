package service

import (
	"forgeflow-api/domain"
	"forgeflow-api/dto"
	"forgeflow-api/errors"
	"forgeflow-api/repository"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type EndpointService struct {
	endpointRepository *repository.EndpointRepository
}

func NewEndpointService(endpointRepository *repository.EndpointRepository) *EndpointService {
	return &EndpointService{
		endpointRepository: endpointRepository,
	}
}

func (s *EndpointService) GetEndpoints(c fiber.Ctx, organizationID uuid.UUID, query dto.PaginateQuery) (dto.PaginateResponse, error) {
	endpoints, total, err := s.endpointRepository.FindAllByOrganizationID(c, organizationID, query)
	if err != nil {
		return dto.PaginateResponse{}, errors.ErrEndpointsNotFound
	}

	outputs := dto.NewEndpointsOutput(endpoints)
	return dto.NewPaginateResponse(outputs, int(total), query), nil
}

func (s *EndpointService) CreateEndpoint(c fiber.Ctx, organizationID uuid.UUID, req dto.CreateEndpointInput) (dto.EndpointOutput, error) {
	endpoint := &domain.Endpoint{
		Name:           req.Name,
		OrganizationID: organizationID,
		BaseURI:        req.BaseURI,
		Path:           req.Path,
		Method:         req.Method,
		Timeout:        req.Timeout,
	}

	err := s.endpointRepository.Create(endpoint)
	if err != nil {
		return dto.EndpointOutput{}, err
	}
	return dto.NewEndpointOutput(*endpoint), nil
}

func (s *EndpointService) UpdateEndpoint(c fiber.Ctx, organizationID uuid.UUID, endpointID uuid.UUID, req dto.UpdateEndpointInput) (dto.EndpointOutput, error) {
	endpoint, err := s.endpointRepository.FindByOrganizationIDAndEndpointID(c, organizationID, endpointID)
	if err != nil {
		return dto.EndpointOutput{}, errors.ErrEndpointNotFound
	}

	endpoint.Name = req.Name
	endpoint.BaseURI = req.BaseURI
	endpoint.Path = req.Path
	endpoint.Method = req.Method
	endpoint.Timeout = req.Timeout

	if err := s.endpointRepository.Update(&endpoint); err != nil {
		return dto.EndpointOutput{}, errors.ErrEndpointFailedToUpdate
	}

	return dto.NewEndpointOutput(endpoint), nil
}

func (s *EndpointService) GetEndpointByID(c fiber.Ctx, organizationID uuid.UUID, endpointID uuid.UUID) (dto.EndpointOutput, error) {
	endpoint, err := s.endpointRepository.FindByOrganizationIDAndEndpointID(c, organizationID, endpointID)
	if err != nil {
		return dto.EndpointOutput{}, errors.ErrEndpointNotFound
	}

	return dto.NewEndpointOutput(endpoint), nil
}
