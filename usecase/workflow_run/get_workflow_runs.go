package workflow_run

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	"flowforge-api/infrastructure/paginate"

	"github.com/google/uuid"
)

type GetWorkflowRunsUseCase struct {
	workflowRepo    *repository.WorkflowRepository
	workflowRunRepo *repository.WorkflowRunRepository
}

func NewGetWorkflowRunsUseCase(
	workflowRepo *repository.WorkflowRepository,
	workflowRunRepo *repository.WorkflowRunRepository,
) *GetWorkflowRunsUseCase {
	return &GetWorkflowRunsUseCase{
		workflowRepo:    workflowRepo,
		workflowRunRepo: workflowRunRepo,
	}
}

func (u *GetWorkflowRunsUseCase) Execute(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID, query paginate.PaginateQuery) ([]entity.WorkflowRun, int64, error) {
	workflow, err := (*u.workflowRepo).GetByIDAndOrganizationID(ctx, organizationID, workflowID)
	if err != nil {
		return nil, 0, err
	}

	workflowRuns, total, err := (*u.workflowRunRepo).GetByWorkflowID(ctx, workflow.ID, query)
	if err != nil {
		return nil, 0, err
	}

	return workflowRuns, total, nil
}
