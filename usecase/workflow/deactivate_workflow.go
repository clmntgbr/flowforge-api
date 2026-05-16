package workflow

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type DeactivateWorkflowUseCase struct {
	workflowRepo *repository.WorkflowRepository
}

func NewDeactivateWorkflowUseCase(workflowRepo *repository.WorkflowRepository) *DeactivateWorkflowUseCase {
	return &DeactivateWorkflowUseCase{workflowRepo: workflowRepo}
}

func (u *DeactivateWorkflowUseCase) Execute(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID) (entity.Workflow, error) {
	workflow, err := (*u.workflowRepo).GetByIDAndOrganizationID(ctx, organizationID, workflowID)
	if err != nil {
		return entity.Workflow{}, err
	}

	workflow.Status = enum.WorkflowStatusInactive

	err = (*u.workflowRepo).Update(ctx, &workflow)
	if err != nil {
		return entity.Workflow{}, err
	}

	return workflow, nil
}
