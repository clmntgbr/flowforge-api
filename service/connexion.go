package service

import (
	"context"
	"forgeflow-api/domain"
	"forgeflow-api/dto"
	"forgeflow-api/errors"
	"forgeflow-api/repository"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ConnexionService struct {
	connexionRepository *repository.ConnexionRepository
}

func NewConnexionService(connexionRepository *repository.ConnexionRepository) *ConnexionService {
	return &ConnexionService{
		connexionRepository: connexionRepository,
	}
}

func (s *ConnexionService) CreateConnexion(c fiber.Ctx, workflowID uuid.UUID, req dto.CreateConnexionInput) (dto.ConnexionOutput, error) {
	fromStepID, err := uuid.Parse(req.From)
	if err != nil {
		return dto.ConnexionOutput{}, errors.ErrInvalidStepID
	}

	toStepID, err := uuid.Parse(req.To)
	if err != nil {
		return dto.ConnexionOutput{}, errors.ErrInvalidStepID
	}

	connexions, err := s.connexionRepository.FindByFromTo(c.Context(), fromStepID, toStepID)
	if err != nil {
		return dto.ConnexionOutput{}, err
	}

	if len(connexions) > 0 {
		return dto.ConnexionOutput{}, nil
	}

	connexion := &domain.Connexion{
		ID:         uuid.New(),
		WorkflowID: workflowID,
		FromStepID: fromStepID,
		ToStepID:   toStepID,
	}

	if err := s.connexionRepository.Create(connexion); err != nil {
		return dto.ConnexionOutput{}, err
	}

	return dto.NewConnexionOutput(*connexion), nil
}

func (s *ConnexionService) DeleteConnexion(ctx context.Context, connexionID uuid.UUID) (dto.ConnexionOutput, error) {
	connexion, err := s.connexionRepository.FindByID(ctx, connexionID)
	if err != nil {
		return dto.ConnexionOutput{}, errors.ErrConnexionNotFound
	}

	if err := s.connexionRepository.Delete(connexion); err != nil {
		return dto.ConnexionOutput{}, err
	}

	return dto.NewConnexionOutput(*connexion), nil
}
