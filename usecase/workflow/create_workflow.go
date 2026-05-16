package workflow

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	workflowDTO "flowforge-api/infrastructure/workflow"

	"github.com/google/uuid"
)

type CreateWorkflowUseCase struct {
	workflowRepo *repository.WorkflowRepository
}

func NewCreateWorkflowUseCase(workflowRepo *repository.WorkflowRepository) *CreateWorkflowUseCase {
	return &CreateWorkflowUseCase{workflowRepo: workflowRepo}
}

func (u *CreateWorkflowUseCase) Execute(ctx context.Context, organizationID uuid.UUID, request workflowDTO.CreateWorkflowInput) (entity.Workflow, error) {
	workflow := &entity.Workflow{
		Name:           request.Name,
		OrganizationID: organizationID,
		Description:    request.Description,
	}

	err := (*u.workflowRepo).Create(ctx, workflow)
	if err != nil {
		return entity.Workflow{}, err
	}

	return *workflow, nil
}
