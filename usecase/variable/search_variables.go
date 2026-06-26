package variable

import (
	"context"
	"encoding/json"
	"flowforge-api/domain/repository"
	variableDTO "flowforge-api/infrastructure/variable"
	"strings"

	"github.com/google/uuid"
)

type SearchVariablesPathUseCase struct {
	variableRepo *repository.VariableRepository
	workflowRepo *repository.WorkflowRepository
	stepRunRepo  *repository.StepRunRepository
}

func NewSearchVariablesPathUseCase(variableRepo *repository.VariableRepository, workflowRepo *repository.WorkflowRepository, stepRunRepo *repository.StepRunRepository) *SearchVariablesPathUseCase {
	return &SearchVariablesPathUseCase{
		variableRepo: variableRepo,
		workflowRepo: workflowRepo,
		stepRunRepo:  stepRunRepo,
	}
}

func (u *SearchVariablesPathUseCase) Execute(ctx context.Context, workflowID uuid.UUID, request variableDTO.SearchVariablesPathInput) ([]string, error) {
	stepRun, err := (*u.stepRunRepo).GetLatestCompletedByStepID(ctx, request.StepID)
	if err != nil {
		return nil, err
	}

	if stepRun.Response == "" {
		return []string{}, nil
	}

	var responseData interface{}
	if err := json.Unmarshal([]byte(stepRun.Response), &responseData); err != nil {
		return nil, err
	}

	paths := extractPaths(responseData, "", request.Query)

	return paths, nil
}

func extractPaths(data interface{}, currentPath string, query string) []string {
	var paths []string

	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			newPath := key
			if currentPath != "" {
				newPath = currentPath + "." + key
			}

			if query == "" || strings.Contains(strings.ToLower(key), strings.ToLower(query)) {
				paths = append(paths, newPath)
			}

			paths = append(paths, extractPaths(value, newPath, query)...)
		}
	case []interface{}:
		if len(v) > 0 {
			paths = append(paths, extractPaths(v[0], currentPath, query)...)
		}
	}

	return paths
}
