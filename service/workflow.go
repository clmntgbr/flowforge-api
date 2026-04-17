package service

import (
	"forgeflow-api/domain"
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

	outputs := dto.NewMinimalWorkflowsOutput(workflows)
	return dto.NewPaginateResponse(outputs, int(total), query), nil
}

func (s *WorkflowService) CreateWorkflow(c fiber.Ctx, organizationID uuid.UUID, req dto.CreateWorkflowInput) (dto.WorkflowOutput, error) {
	workflow := &domain.Workflow{
		Name:           req.Name,
		OrganizationID: organizationID,
		Description:    req.Description,
	}

	err := s.workflowRepository.Create(workflow)
	if err != nil {
		return dto.WorkflowOutput{}, err
	}
	return dto.NewWorkflowOutput(*workflow), nil
}

func (s *WorkflowService) UpdateWorkflow(c fiber.Ctx, organizationID uuid.UUID, workflowID uuid.UUID, req dto.UpdateWorkflowInput) (dto.WorkflowOutput, error) {
	workflow, err := s.workflowRepository.FindByOrganizationIDAndWorkflowID(c, organizationID, workflowID)
	if err != nil {
		return dto.WorkflowOutput{}, errors.ErrWorkflowNotFound
	}

	workflow.Name = req.Name
	workflow.Description = req.Description

	if err := s.workflowRepository.Update(&workflow); err != nil {
		return dto.WorkflowOutput{}, errors.ErrWorkflowFailedToUpdate
	}

	return dto.NewWorkflowOutput(workflow), nil
}

func (s *WorkflowService) GetWorkflowByID(c fiber.Ctx, organizationID uuid.UUID, workflowID uuid.UUID) (dto.WorkflowOutput, error) {
	workflow, err := s.workflowRepository.FindByOrganizationIDAndWorkflowID(c, organizationID, workflowID)
	if err != nil {
		return dto.WorkflowOutput{}, errors.ErrWorkflowNotFound
	}

	return dto.NewWorkflowOutput(workflow), nil
}
