package service

import (
	"forgeflow-api/domain"
	"forgeflow-api/dto"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UserServiceInterface interface {
	GetUser(user *domain.User) (*dto.UserOutput, error)
}

type EndpointServiceInterface interface {
	GetEndpoints(c fiber.Ctx, organizationID uuid.UUID, query dto.PaginateQuery) (dto.PaginateResponse, error)
	CreateEndpoint(c fiber.Ctx, organizationID uuid.UUID, req dto.CreateEndpointInput) (dto.EndpointOutput, error)
	GetEndpointByID(c fiber.Ctx, organizationID uuid.UUID, endpointID uuid.UUID) (dto.EndpointOutput, error)
	UpdateEndpoint(c fiber.Ctx, organizationID uuid.UUID, endpointID uuid.UUID, req dto.UpdateEndpointInput) (dto.EndpointOutput, error)
}
