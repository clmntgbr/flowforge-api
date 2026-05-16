package workflow

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type ActivateWorkflowUseCase struct {
	workflowRepo *repository.WorkflowRepository
}

func NewActivateWorkflowUseCase(workflowRepo *repository.WorkflowRepository) *ActivateWorkflowUseCase {
	return &ActivateWorkflowUseCase{workflowRepo: workflowRepo}
}

func (u *ActivateWorkflowUseCase) Execute(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID) (entity.Workflow, error) {
	workflow, err := (*u.workflowRepo).GetByIDAndOrganizationID(ctx, organizationID, workflowID)
	if err != nil {
		return entity.Workflow{}, err
	}

	workflow.Status = enum.WorkflowStatusActive

	err = (*u.workflowRepo).Update(ctx, &workflow)
	if err != nil {
		return entity.Workflow{}, err
	}

	return workflow, nil
}
