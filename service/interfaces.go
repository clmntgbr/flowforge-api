package service

import (
	"context"
	"forgeflow-api/domain"
	"forgeflow-api/dto"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UserServiceInterface interface {
	GetUser(user *domain.User) (*dto.UserOutput, error)
	CreateUser(c fiber.Ctx, id string, firstName string, lastName string, banned bool) (*domain.User, error)
	UpdateUser(c fiber.Ctx, id string, firstName string, lastName string, banned bool) error
	DeleteUser(c fiber.Ctx, id string) error
}

type EndpointServiceInterface interface {
	GetEndpoints(c fiber.Ctx, organizationID uuid.UUID, query dto.PaginateQuery) (dto.PaginateResponse, error)
	CreateEndpoint(c fiber.Ctx, organizationID uuid.UUID, req dto.CreateEndpointInput) (dto.EndpointOutput, error)
	GetEndpointByID(c fiber.Ctx, organizationID uuid.UUID, endpointID uuid.UUID) (dto.EndpointOutput, error)
	UpdateEndpoint(c fiber.Ctx, organizationID uuid.UUID, endpointID uuid.UUID, req dto.UpdateEndpointInput) (dto.EndpointOutput, error)
}

type OrganizationServiceInterface interface {
	GetOrganizations(c fiber.Ctx, user *domain.User, activeOrganizationID uuid.UUID) ([]dto.MinimalOrganizationOutput, error)
	GetOrganizationByID(c fiber.Ctx, user *domain.User, organizationID uuid.UUID) (dto.OrganizationOutput, error)
	CreateOrganization(c fiber.Ctx, user *domain.User, name string) (dto.OrganizationOutput, error)
	UpdateOrganization(c fiber.Ctx, user *domain.User, organizationID uuid.UUID, req dto.UpdateOrganizationInput) (dto.OrganizationOutput, error)
	ActivateOrganization(ctx context.Context, userID uuid.UUID, organizationID uuid.UUID) (dto.OrganizationOutput, error)
}

type WorkflowServiceInterface interface {
	GetWorkflows(c fiber.Ctx, organizationID uuid.UUID, query dto.PaginateQuery) (dto.PaginateResponse, error)
	CreateWorkflow(c fiber.Ctx, organizationID uuid.UUID, req dto.CreateWorkflowInput) (dto.WorkflowOutput, error)
	GetWorkflowByID(c fiber.Ctx, organizationID uuid.UUID, workflowID uuid.UUID) (dto.WorkflowOutput, error)
	UpdateWorkflow(c fiber.Ctx, organizationID uuid.UUID, workflowID uuid.UUID, req dto.UpdateWorkflowInput) (dto.WorkflowOutput, error)
}

type StepServiceInterface interface {
	GetStepByID(ctx context.Context, organizationID uuid.UUID, stepID uuid.UUID) (dto.StepOutput, error)
	UpdateStep(ctx context.Context, organizationID uuid.UUID, stepID uuid.UUID, req dto.UpdateStepInput) (dto.StepOutput, error)
	UpsertSteps(ctx context.Context, workflowID uuid.UUID, steps []dto.UpdateWorkflowStepInput) error
}

type ConnexionServiceInterface interface {
	CreateConnexion(c fiber.Ctx, workflowID uuid.UUID, req dto.CreateConnexionInput) (dto.ConnexionOutput, error)
	DeleteConnexion(ctx context.Context, connexionID uuid.UUID) (dto.ConnexionOutput, error)
}
