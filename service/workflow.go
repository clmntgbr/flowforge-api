package service

import (
	"forgeflow-api/dto"
	"forgeflow-api/errors"
	"forgeflow-api/repository"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type WorkflowService struct {
	workflowRepository *repository.WorkflowRepository
}

func NewWorkflowService(workflowRepository *repository.WorkflowRepository) *WorkflowService {
	return &WorkflowService{
		workflowRepository: workflowRepository,
	}
}

func (s *WorkflowService) GetWorkflows(c fiber.Ctx, organizationID uuid.UUID, query dto.PaginateQuery) (dto.PaginateResponse, error) {
	workflows, total, err := s.workflowRepository.FindAllByOrganizationID(c, organizationID, query)
	if err != nil {
		return dto.PaginateResponse{}, errors.ErrWorkflowsNotFound
	}

	outputs := dto.NewWorkflowsOutput(workflows)
	return dto.NewPaginateResponse(outputs, int(total), query), nil
}
