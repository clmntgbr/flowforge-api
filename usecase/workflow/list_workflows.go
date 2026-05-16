package workflow

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	"flowforge-api/infrastructure/paginate"

	"github.com/google/uuid"
)

type ListWorkflowsUseCase struct {
	workflowRepo *repository.WorkflowRepository
}

func NewListWorkflowsUseCase(workflowRepo *repository.WorkflowRepository) *ListWorkflowsUseCase {
	return &ListWorkflowsUseCase{workflowRepo: workflowRepo}
}

func (u *ListWorkflowsUseCase) Execute(ctx context.Context, organizationID uuid.UUID, query paginate.PaginateQuery) ([]entity.Workflow, int64, error) {
	workflows, total, err := (*u.workflowRepo).List(ctx, organizationID, query)
	if err != nil {
		return []entity.Workflow{}, 0, err
	}

	return workflows, total, nil
}
