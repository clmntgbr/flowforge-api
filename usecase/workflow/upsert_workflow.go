package workflow

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	workflowDTO "flowforge-api/infrastructure/workflow"

	"github.com/google/uuid"
)

type UpsertWorkflowUseCase struct {
	workflowRepo repository.WorkflowRepository
}

func NewUpsertWorkflowUseCase(workflowRepo repository.WorkflowRepository) *UpsertWorkflowUseCase {
	return &UpsertWorkflowUseCase{workflowRepo: workflowRepo}
}

func (u *UpsertWorkflowUseCase) Execute(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID, request workflowDTO.UpsertWorkflowInput) (entity.Workflow, error) {
	workflow, err := u.workflowRepo.GetByIDAndOrganizationID(ctx, organizationID, workflowID)
	if err != nil {
		return entity.Workflow{}, err
	}

	return workflow, nil
}
