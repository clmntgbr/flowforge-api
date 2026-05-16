package workflow

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	workflowDTO "flowforge-api/infrastructure/workflow"

	"github.com/google/uuid"
)

type UpdateWorkflowUseCase struct {
	workflowRepo *repository.WorkflowRepository
}

func NewUpdateWorkflowUseCase(workflowRepo *repository.WorkflowRepository) *UpdateWorkflowUseCase {
	return &UpdateWorkflowUseCase{workflowRepo: workflowRepo}
}

func (u *UpdateWorkflowUseCase) Execute(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID, request workflowDTO.UpdateWorkflowInput) (entity.Workflow, error) {
	workflow, err := (*u.workflowRepo).GetByIDAndOrganizationID(ctx, organizationID, workflowID)
	if err != nil {
		return entity.Workflow{}, err
	}

	workflow.Name = request.Name
	workflow.Description = request.Description

	err = (*u.workflowRepo).Update(ctx, &workflow)
	if err != nil {
		return entity.Workflow{}, err
	}

	return workflow, nil
}
