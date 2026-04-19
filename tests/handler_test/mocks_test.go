package handler_test

import (
	"forgeflow-api/domain"
	"forgeflow-api/dto"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type MockUserService struct {
	GetUserFunc func(user *domain.User) (*dto.UserOutput, error)
}

func (m *MockUserService) GetUser(user *domain.User) (*dto.UserOutput, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(user)
	}
	return nil, nil
}

func NewMockUserService() *MockUserService {
	return &MockUserService{}
}

type MockEndpointService struct {
	GetEndpointsFunc    func(c fiber.Ctx, organizationID uuid.UUID, query dto.PaginateQuery) (dto.PaginateResponse, error)
	CreateEndpointFunc  func(c fiber.Ctx, organizationID uuid.UUID, req dto.CreateEndpointInput) (dto.EndpointOutput, error)
	GetEndpointByIDFunc func(c fiber.Ctx, organizationID uuid.UUID, endpointID uuid.UUID) (dto.EndpointOutput, error)
	UpdateEndpointFunc  func(c fiber.Ctx, organizationID uuid.UUID, endpointID uuid.UUID, req dto.UpdateEndpointInput) (dto.EndpointOutput, error)
}

func (m *MockEndpointService) GetEndpoints(c fiber.Ctx, organizationID uuid.UUID, query dto.PaginateQuery) (dto.PaginateResponse, error) {
	if m.GetEndpointsFunc != nil {
		return m.GetEndpointsFunc(c, organizationID, query)
	}
	return dto.PaginateResponse{}, nil
}

func (m *MockEndpointService) CreateEndpoint(c fiber.Ctx, organizationID uuid.UUID, req dto.CreateEndpointInput) (dto.EndpointOutput, error) {
	if m.CreateEndpointFunc != nil {
		return m.CreateEndpointFunc(c, organizationID, req)
	}
	return dto.EndpointOutput{}, nil
}

func (m *MockEndpointService) GetEndpointByID(c fiber.Ctx, organizationID uuid.UUID, endpointID uuid.UUID) (dto.EndpointOutput, error) {
	if m.GetEndpointByIDFunc != nil {
		return m.GetEndpointByIDFunc(c, organizationID, endpointID)
	}
	return dto.EndpointOutput{}, nil
}

func (m *MockEndpointService) UpdateEndpoint(c fiber.Ctx, organizationID uuid.UUID, endpointID uuid.UUID, req dto.UpdateEndpointInput) (dto.EndpointOutput, error) {
	if m.UpdateEndpointFunc != nil {
		return m.UpdateEndpointFunc(c, organizationID, endpointID, req)
	}
	return dto.EndpointOutput{}, nil
}

func NewMockEndpointService() *MockEndpointService {
	return &MockEndpointService{}
}
