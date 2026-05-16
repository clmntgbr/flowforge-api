package workflow

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type GetWorkflowUseCase struct {
	workflowRepo *repository.WorkflowRepository
}

func NewGetWorkflowUseCase(workflowRepo *repository.WorkflowRepository) *GetWorkflowUseCase {
	return &GetWorkflowUseCase{workflowRepo: workflowRepo}
}

func (u *GetWorkflowUseCase) Execute(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID) (entity.Workflow, error) {
	workflow, err := (*u.workflowRepo).GetByIDAndOrganizationID(ctx, organizationID, workflowID)
	if err != nil {
		return entity.Workflow{}, err
	}

	return workflow, nil
}
