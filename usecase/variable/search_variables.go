package variable

import (
	"context"
	"encoding/json"
	"flowforge-api/domain/repository"
	variableDTO "flowforge-api/infrastructure/variable"
	"strconv"
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

func (u *SearchVariablesPathUseCase) Execute(ctx context.Context, workflowID uuid.UUID, request variableDTO.SearchVariablesPathInput) ([]string, int, error) {
	stepRun, err := (*u.stepRunRepo).GetLatestCompletedByStepID(ctx, request.StepID)
	if err != nil {
		return nil, 0, err
	}

	if stepRun.Response == "" {
		return []string{}, 0, nil
	}

	var responseData interface{}
	if err := json.Unmarshal([]byte(stepRun.Response), &responseData); err != nil {
		return nil, 0, err
	}

	allPaths := extractPaths(responseData, "", request.Query)
	total := len(allPaths)

	start := request.Offset()
	end := start + request.Limit

	if start >= total {
		return []string{}, total, nil
	}

	if end > total {
		end = total
	}

	paginatedPaths := allPaths[start:end]

	return paginatedPaths, total, nil
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
		maxElements := len(v)
		if maxElements > 5 {
			maxElements = 5
		}
		
		for i := 0; i < maxElements; i++ {
			indexPath := currentPath
			if currentPath == "" {
				indexPath = "[" + strconv.Itoa(i) + "]"
			} else {
				indexPath = currentPath + "[" + strconv.Itoa(i) + "]"
			}
			paths = append(paths, extractPaths(v[i], indexPath, query)...)
		}
	}

	return paths
}
