package service

import (
	"context"
	"forgeflow-api/domain"
	"forgeflow-api/dto"
	"forgeflow-api/errors"
	"forgeflow-api/repository"

	"github.com/google/uuid"
)

type StepService struct {
	stepRepository     *repository.StepRepository
	endpointRepository *repository.EndpointRepository
}

func NewStepService(stepRepository *repository.StepRepository, endpointRepository *repository.EndpointRepository) *StepService {
	return &StepService{
		stepRepository:     stepRepository,
		endpointRepository: endpointRepository,
	}
}

func (s *StepService) UpsertSteps(ctx context.Context, workflowID uuid.UUID, stepsInput []dto.UpdateWorkflowStepInput) error {
	existingSteps, err := s.stepRepository.FindByWorkflowID(ctx, workflowID)
	if err != nil {
		return err
	}

	receivedStepIDs := make(map[uuid.UUID]bool)
	for _, stepInput := range stepsInput {
		stepUUID, err := uuid.Parse(stepInput.ID)
		if err != nil {
			return errors.ErrInvalidRequestBody
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
		if err := s.stepRepository.DeleteByIDs(ctx, stepsToDelete); err != nil {
			return err
		}
	}

	for _, stepInput := range stepsInput {
		stepUUID, err := uuid.Parse(stepInput.ID)
		if err != nil {
			return errors.ErrInvalidRequestBody
		}

		endpointUUID, err := uuid.Parse(stepInput.EndpointID)
		if err != nil {
			return errors.ErrInvalidRequestBody
		}

		endpoint, err := s.endpointRepository.FindByID(ctx, endpointUUID)
		if err != nil {
			return errors.ErrEndpointNotFound
		}

		index := stepInput.Index

		position := domain.Position{X: stepInput.Position.X, Y: stepInput.Position.Y}

		existingStep, err := s.stepRepository.FindByID(ctx, stepUUID)

		if existingStep == nil {
			newStep := &domain.Step{
				ID:          stepUUID,
				Name:        endpoint.Name,
				Description: endpoint.BaseURI + endpoint.Path,
				Timeout:     endpoint.Timeout,
				Query:       endpoint.Query,
				Position:    position,
				Index:       index,
				EndpointID:  endpointUUID,
				WorkflowID:  workflowID,
			}
			s.stepRepository.Create(newStep)
		} else {
			if err := s.stepRepository.UpdatePositionAndIndex(ctx, existingStep.ID, workflowID, position, index); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *StepService) GetStepByID(ctx context.Context, organizationID uuid.UUID, stepID uuid.UUID) (dto.StepOutput, error) {
	step, err := s.stepRepository.FindByOrganizationIDAndStepID(ctx, organizationID, stepID)
	if err != nil {
		return dto.StepOutput{}, errors.ErrStepNotFound
	}

	return dto.NewStepOutput(*step), nil
}

func (s *StepService) UpdateStep(ctx context.Context, organizationID uuid.UUID, stepID uuid.UUID, stepInput dto.UpdateStepInput) (dto.StepOutput, error) {
	step, err := s.stepRepository.FindByOrganizationIDAndStepID(ctx, organizationID, stepID)
	if err != nil {
		return dto.StepOutput{}, errors.ErrStepNotFound
	}

	step.Name = stepInput.Name
	step.Description = stepInput.Description
	step.Timeout = stepInput.Timeout
	step.Query = stepInput.Query

	if err := s.stepRepository.Update(step); err != nil {
		return dto.StepOutput{}, errors.ErrStepFailedToUpdate
	}

	return dto.NewStepOutput(*step), nil
}
