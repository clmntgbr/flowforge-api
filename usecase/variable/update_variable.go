package variable

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	variableDTO "flowforge-api/infrastructure/variable"
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

type UpdateVariableUseCase struct {
	variableRepo *repository.VariableRepository
	stepRepo     *repository.StepRepository
}

func NewUpdateVariableUseCase(variableRepo *repository.VariableRepository, stepRepo *repository.StepRepository) *UpdateVariableUseCase {
	return &UpdateVariableUseCase{
		variableRepo: variableRepo,
		stepRepo:     stepRepo,
	}
}

func (u *UpdateVariableUseCase) Execute(ctx context.Context, workflowID uuid.UUID, variableID uuid.UUID, request variableDTO.UpdateVariableInput) (entity.Variable, error) {

	variable, err := (*u.variableRepo).GetVariableByIDAndWorkflowID(ctx, workflowID, variableID)
	if err != nil {
		return entity.Variable{}, err
	}

	oldKey := variable.Key

	variable.Name = request.Name
	variable.Key = request.Key
	variable.Path = request.Path
	variable.Description = request.Description
	variable.StepID = request.StepID

	if oldKey != request.Key {
		if err := u.replaceKeyInWorkflowSteps(ctx, workflowID, oldKey, request.Key); err != nil {
			return entity.Variable{}, fmt.Errorf("failed to replace key in workflow steps: %w", err)
		}
	}

	err = (*u.variableRepo).Update(ctx, &variable)
	if err != nil {
		return entity.Variable{}, err
	}

	return variable, nil
}

func (u *UpdateVariableUseCase) replaceKeyInWorkflowSteps(ctx context.Context, workflowID uuid.UUID, oldKey, newKey string) error {
	steps, err := (*u.stepRepo).GetByWorkflowID(ctx, workflowID)
	if err != nil {
		return err
	}

	for i := range steps {
		step := &steps[i]
		modified := false

		step.URL = replaceKeyInString(step.URL, oldKey, newKey)
		if step.URL != steps[i].URL {
			modified = true
		}

		for j := range step.Header {
			originalValue := step.Header[j].Value
			step.Header[j].Value = replaceKeyInString(step.Header[j].Value, oldKey, newKey)
			if step.Header[j].Value != originalValue {
				modified = true
			}
		}

		for j := range step.Query {
			originalValue := step.Query[j].Value
			step.Query[j].Value = replaceKeyInString(step.Query[j].Value, oldKey, newKey)
			if step.Query[j].Value != originalValue {
				modified = true
			}
		}

		if len(step.Body) > 0 {
			originalBody := string(step.Body)
			newBody := replaceKeyInString(originalBody, oldKey, newKey)
			if newBody != originalBody {
				step.Body = []byte(newBody)
				modified = true
			}
		}

		if modified {
			if err := (*u.stepRepo).Update(ctx, step); err != nil {
				return fmt.Errorf("failed to update step %s: %w", step.ID, err)
			}
		}
	}

	return nil
}

func replaceKeyInString(input, oldKey, newKey string) string {
	pattern := regexp.MustCompile(`\{\{\s*` + regexp.QuoteMeta(oldKey) + `\s*\}\}`)
	return pattern.ReplaceAllString(input, "{{ "+newKey+" }}")
}
