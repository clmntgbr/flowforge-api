package workflow

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	workflowDTO "flowforge-api/infrastructure/workflow"
	usecase "flowforge-api/usecase/step"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpsertWorkflowUseCase struct {
	workflowRepo                   repository.WorkflowRepository
	stepRepo                       repository.StepRepository
	endpointRepo                   repository.EndpointRepository
	calculateExecutionOrderUseCase usecase.CalculateExecutionOrderUseCase
}

func NewUpsertWorkflowUseCase(
	workflowRepo repository.WorkflowRepository,
	stepRepo repository.StepRepository,
	endpointRepo repository.EndpointRepository,
	calculateExecutionOrderUseCase usecase.CalculateExecutionOrderUseCase,
) *UpsertWorkflowUseCase {
	return &UpsertWorkflowUseCase{
		workflowRepo:                   workflowRepo,
		stepRepo:                       stepRepo,
		endpointRepo:                   endpointRepo,
		calculateExecutionOrderUseCase: calculateExecutionOrderUseCase,
	}
}

func (u *UpsertWorkflowUseCase) Execute(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID, request workflowDTO.UpsertWorkflowInput) error {
	workflow, err := u.workflowRepo.GetByIDAndOrganizationID(ctx, organizationID, workflowID)
	if err != nil {
		return err
	}

	return u.workflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		existingSteps, err := u.stepRepo.GetByWorkflowID(ctx, workflow.ID)
		if err != nil {
			return err
		}

		receivedStepIDs := make(map[uuid.UUID]bool)
		for _, stepInput := range request.Steps {
			stepUUID, err := uuid.Parse(stepInput.ID)
			if err != nil {
				return err
			}
			receivedStepIDs[stepUUID] = true
		}

		stepsToDelete := make([]uuid.UUID, 0)
		for _, existingStep := range existingSteps {
			if !receivedStepIDs[existingStep.ID] {
				stepsToDelete = append(stepsToDelete, existingStep.ID)
			}
		}

		if len(stepsToDelete) > 0 {
			if err := u.stepRepo.DeleteByIDs(ctx, stepsToDelete); err != nil {
				return err
			}
		}

		for _, stepInput := range request.Steps {
			stepUUID, err := uuid.Parse(stepInput.ID)
			if err != nil {
				return err
			}

			endpointUUID, err := uuid.Parse(stepInput.EndpointID)
			if err != nil {
				return err
			}

			endpoint, err := u.endpointRepo.GetByID(ctx, endpointUUID)
			if err != nil {
				return err
			}

			index := stepInput.Index
			executionOrder := u.calculateExecutionOrderUseCase.Execute(ctx, index)

			position := entity.Position{X: stepInput.Position.X, Y: stepInput.Position.Y}

			existingStep, _ := u.stepRepo.GetByID(ctx, stepUUID)

			if existingStep == nil {
				newStep := &entity.Step{
					ID:             stepUUID,
					Name:           endpoint.Name,
					Description:    endpoint.BaseURI + endpoint.Path,
					Timeout:        endpoint.Timeout,
					Query:          endpoint.Query,
					Header:         endpoint.Header,
					Body:           endpoint.Body,
					Position:       position,
					Index:          index,
					ExecutionOrder: executionOrder,
					EndpointID:     endpointUUID,
					WorkflowID:     workflowID,
					RetryOnFailure: endpoint.RetryOnFailure,
					RetryCount:     endpoint.RetryCount,
					RetryDelay:     endpoint.RetryDelay,
				}
				if err := u.stepRepo.Create(ctx, newStep); err != nil {
					return err
				}
			} else {
				if err := u.stepRepo.UpdatePositionAndIndex(ctx, existingStep.ID, workflowID, position, index, executionOrder); err != nil {
					return err
				}
			}
		}

		return nil
	})
}
