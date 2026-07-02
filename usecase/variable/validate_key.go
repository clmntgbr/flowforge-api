package variable

import (
	"context"
	"errors"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

func validateUniqueKey(ctx context.Context, variableRepo *repository.VariableRepository, workflowID uuid.UUID, key string, excludeVariableID *uuid.UUID) error {
	existing, err := (*variableRepo).GetVariableByWorkflowIDAndKey(ctx, workflowID, key)
	if err != nil {
		return err
	}

	if existing == nil {
		return nil
	}

	if excludeVariableID != nil && existing.ID == *excludeVariableID {
		return nil
	}

	return errors.New("variable key already exists in workflow")
}
