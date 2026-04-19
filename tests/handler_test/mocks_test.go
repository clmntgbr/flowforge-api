package handler_test

import (
	"context"
	"forgeflow-api/domain"
	"forgeflow-api/dto"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type MockUserService struct {
	GetUserFunc    func(user *domain.User) (*dto.UserOutput, error)
	CreateUserFunc func(c fiber.Ctx, id string, firstName string, lastName string, banned bool) (*domain.User, error)
	UpdateUserFunc func(c fiber.Ctx, id string, firstName string, lastName string, banned bool) error
	DeleteUserFunc func(c fiber.Ctx, id string) error
}

func (m *MockUserService) GetUser(user *domain.User) (*dto.UserOutput, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(user)
	}
	return nil, nil
}

func (m *MockUserService) CreateUser(c fiber.Ctx, id string, firstName string, lastName string, banned bool) (*domain.User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(c, id, firstName, lastName, banned)
	}
	return nil, nil
}

func (m *MockUserService) UpdateUser(c fiber.Ctx, id string, firstName string, lastName string, banned bool) error {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(c, id, firstName, lastName, banned)
	}
	return nil
}

func (m *MockUserService) DeleteUser(c fiber.Ctx, id string) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(c, id)
	}
	return nil
}

func NewMockUserService() *MockUserService {
	return &MockUserService{}
}

type MockUserRepository struct {
	FindByClerkIDFunc func(clerkID string) (*domain.User, error)
	UpdateFunc        func(user *domain.User) error
}

func (m *MockUserRepository) FindByClerkID(clerkID string) (*domain.User, error) {
	if m.FindByClerkIDFunc != nil {
		return m.FindByClerkIDFunc(clerkID)
	}
	return nil, nil
}

func (m *MockUserRepository) Update(user *domain.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(user)
	}
	return nil
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{}
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

type MockOrganizationService struct {
	GetOrganizationsFunc    func(c fiber.Ctx, user *domain.User, activeOrganizationID uuid.UUID) ([]dto.MinimalOrganizationOutput, error)
	GetOrganizationByIDFunc func(c fiber.Ctx, user *domain.User, organizationID uuid.UUID) (dto.OrganizationOutput, error)
	CreateOrganizationFunc  func(c fiber.Ctx, user *domain.User, name string) (dto.OrganizationOutput, error)
	UpdateOrganizationFunc  func(c fiber.Ctx, user *domain.User, organizationID uuid.UUID, req dto.UpdateOrganizationInput) (dto.OrganizationOutput, error)
	ActivateOrganizationFunc func(ctx context.Context, userID uuid.UUID, organizationID uuid.UUID) (dto.OrganizationOutput, error)
}

func (m *MockOrganizationService) GetOrganizations(c fiber.Ctx, user *domain.User, activeOrganizationID uuid.UUID) ([]dto.MinimalOrganizationOutput, error) {
	if m.GetOrganizationsFunc != nil {
		return m.GetOrganizationsFunc(c, user, activeOrganizationID)
	}
	return nil, nil
}

func (m *MockOrganizationService) GetOrganizationByID(c fiber.Ctx, user *domain.User, organizationID uuid.UUID) (dto.OrganizationOutput, error) {
	if m.GetOrganizationByIDFunc != nil {
		return m.GetOrganizationByIDFunc(c, user, organizationID)
	}
	return dto.OrganizationOutput{}, nil
}

func (m *MockOrganizationService) CreateOrganization(c fiber.Ctx, user *domain.User, name string) (dto.OrganizationOutput, error) {
	if m.CreateOrganizationFunc != nil {
		return m.CreateOrganizationFunc(c, user, name)
	}
	return dto.OrganizationOutput{}, nil
}

func (m *MockOrganizationService) UpdateOrganization(c fiber.Ctx, user *domain.User, organizationID uuid.UUID, req dto.UpdateOrganizationInput) (dto.OrganizationOutput, error) {
	if m.UpdateOrganizationFunc != nil {
		return m.UpdateOrganizationFunc(c, user, organizationID, req)
	}
	return dto.OrganizationOutput{}, nil
}

func (m *MockOrganizationService) ActivateOrganization(ctx context.Context, userID uuid.UUID, organizationID uuid.UUID) (dto.OrganizationOutput, error) {
	if m.ActivateOrganizationFunc != nil {
		return m.ActivateOrganizationFunc(ctx, userID, organizationID)
	}
	return dto.OrganizationOutput{}, nil
}

func NewMockOrganizationService() *MockOrganizationService {
	return &MockOrganizationService{}
}

type MockWorkflowService struct {
	GetWorkflowsFunc    func(c fiber.Ctx, organizationID uuid.UUID, query dto.PaginateQuery) (dto.PaginateResponse, error)
	CreateWorkflowFunc  func(c fiber.Ctx, organizationID uuid.UUID, req dto.CreateWorkflowInput) (dto.WorkflowOutput, error)
	GetWorkflowByIDFunc func(c fiber.Ctx, organizationID uuid.UUID, workflowID uuid.UUID) (dto.WorkflowOutput, error)
	UpdateWorkflowFunc  func(c fiber.Ctx, organizationID uuid.UUID, workflowID uuid.UUID, req dto.UpdateWorkflowInput) (dto.WorkflowOutput, error)
}

func (m *MockWorkflowService) GetWorkflows(c fiber.Ctx, organizationID uuid.UUID, query dto.PaginateQuery) (dto.PaginateResponse, error) {
	if m.GetWorkflowsFunc != nil {
		return m.GetWorkflowsFunc(c, organizationID, query)
	}
	return dto.PaginateResponse{}, nil
}

func (m *MockWorkflowService) CreateWorkflow(c fiber.Ctx, organizationID uuid.UUID, req dto.CreateWorkflowInput) (dto.WorkflowOutput, error) {
	if m.CreateWorkflowFunc != nil {
		return m.CreateWorkflowFunc(c, organizationID, req)
	}
	return dto.WorkflowOutput{}, nil
}

func (m *MockWorkflowService) GetWorkflowByID(c fiber.Ctx, organizationID uuid.UUID, workflowID uuid.UUID) (dto.WorkflowOutput, error) {
	if m.GetWorkflowByIDFunc != nil {
		return m.GetWorkflowByIDFunc(c, organizationID, workflowID)
	}
	return dto.WorkflowOutput{}, nil
}

func (m *MockWorkflowService) UpdateWorkflow(c fiber.Ctx, organizationID uuid.UUID, workflowID uuid.UUID, req dto.UpdateWorkflowInput) (dto.WorkflowOutput, error) {
	if m.UpdateWorkflowFunc != nil {
		return m.UpdateWorkflowFunc(c, organizationID, workflowID, req)
	}
	return dto.WorkflowOutput{}, nil
}

func NewMockWorkflowService() *MockWorkflowService {
	return &MockWorkflowService{}
}

type MockStepService struct {
	GetStepByIDFunc func(ctx context.Context, organizationID uuid.UUID, stepID uuid.UUID) (dto.StepOutput, error)
	UpdateStepFunc  func(ctx context.Context, organizationID uuid.UUID, stepID uuid.UUID, req dto.UpdateStepInput) (dto.StepOutput, error)
	UpsertStepsFunc func(ctx context.Context, workflowID uuid.UUID, steps []dto.UpdateWorkflowStepInput) error
}

func (m *MockStepService) GetStepByID(ctx context.Context, organizationID uuid.UUID, stepID uuid.UUID) (dto.StepOutput, error) {
	if m.GetStepByIDFunc != nil {
		return m.GetStepByIDFunc(ctx, organizationID, stepID)
	}
	return dto.StepOutput{}, nil
}

func (m *MockStepService) UpdateStep(ctx context.Context, organizationID uuid.UUID, stepID uuid.UUID, req dto.UpdateStepInput) (dto.StepOutput, error) {
	if m.UpdateStepFunc != nil {
		return m.UpdateStepFunc(ctx, organizationID, stepID, req)
	}
	return dto.StepOutput{}, nil
}

func (m *MockStepService) UpsertSteps(ctx context.Context, workflowID uuid.UUID, steps []dto.UpdateWorkflowStepInput) error {
	if m.UpsertStepsFunc != nil {
		return m.UpsertStepsFunc(ctx, workflowID, steps)
	}
	return nil
}

func NewMockStepService() *MockStepService {
	return &MockStepService{}
}

type MockConnexionService struct {
	CreateConnexionFunc func(c fiber.Ctx, workflowID uuid.UUID, req dto.CreateConnexionInput) (dto.ConnexionOutput, error)
	DeleteConnexionFunc func(ctx context.Context, connexionID uuid.UUID) (dto.ConnexionOutput, error)
}

func (m *MockConnexionService) CreateConnexion(c fiber.Ctx, workflowID uuid.UUID, req dto.CreateConnexionInput) (dto.ConnexionOutput, error) {
	if m.CreateConnexionFunc != nil {
		return m.CreateConnexionFunc(c, workflowID, req)
	}
	return dto.ConnexionOutput{}, nil
}

func (m *MockConnexionService) DeleteConnexion(ctx context.Context, connexionID uuid.UUID) (dto.ConnexionOutput, error) {
	if m.DeleteConnexionFunc != nil {
		return m.DeleteConnexionFunc(ctx, connexionID)
	}
	return dto.ConnexionOutput{}, nil
}

func NewMockConnexionService() *MockConnexionService {
	return &MockConnexionService{}
}
