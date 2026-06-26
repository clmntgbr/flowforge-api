package variable

import (
	"context"
	"encoding/json"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

type ReplaceVariablesUseCase struct {
	variableRepo *repository.VariableRepository
	stepRunRepo  *repository.StepRunRepository
}

func NewReplaceVariablesUseCase(variableRepo *repository.VariableRepository, stepRunRepo *repository.StepRunRepository) *ReplaceVariablesUseCase {
	return &ReplaceVariablesUseCase{
		variableRepo: variableRepo,
		stepRunRepo:  stepRunRepo,
	}
}

func (u *ReplaceVariablesUseCase) Execute(ctx context.Context, step *entity.Step, workflowRunID uuid.UUID) (*entity.Step, error) {
	variables, err := (*u.variableRepo).GetVariablesByWorkflowID(ctx, step.WorkflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get variables: %w", err)
	}

	if len(variables) == 0 {
		return step, nil
	}

	stepCopy := *step

	stepCopy.URL = u.replaceInString(ctx, stepCopy.URL, variables, workflowRunID)

	for i := range stepCopy.Header {
		stepCopy.Header[i].Value = u.replaceInString(ctx, stepCopy.Header[i].Value, variables, workflowRunID)
	}

	for i := range stepCopy.Query {
		stepCopy.Query[i].Value = u.replaceInString(ctx, stepCopy.Query[i].Value, variables, workflowRunID)
	}

	if len(stepCopy.Body) > 0 {
		bodyStr := string(stepCopy.Body)
		bodyStr = u.replaceInString(ctx, bodyStr, variables, workflowRunID)
		stepCopy.Body = []byte(bodyStr)
	}

	return &stepCopy, nil
}

func (u *ReplaceVariablesUseCase) replaceInString(ctx context.Context, input string, variables []entity.Variable, workflowRunID uuid.UUID) string {
	variablePattern := regexp.MustCompile(`\{\{([^}]+)\}\}`)

	return variablePattern.ReplaceAllStringFunc(input, func(match string) string {
		variableName := strings.TrimSpace(match[2 : len(match)-2])
		log.Printf("🔍 Attempting to replace variable: '%s' (match: '%s')", variableName, match)

		var targetVariable *entity.Variable
		for i := range variables {
			if variables[i].Name == variableName {
				targetVariable = &variables[i]
				break
			}
		}

		if targetVariable == nil {
			log.Printf("❌ Variable '%s' not found in workflow variables (available: %d variables)", variableName, len(variables))
			for _, v := range variables {
				log.Printf("   - Available variable: '%s' (StepID: %s, Path: %s)", v.Name, v.StepID, v.Path)
			}
			return match
		}

		log.Printf("✅ Variable '%s' found: StepID=%s, Path=%s", variableName, targetVariable.StepID, targetVariable.Path)

		stepRuns, err := (*u.stepRunRepo).GetAllByWorkflowRunID(ctx, workflowRunID)
		if err != nil {
			log.Printf("❌ Failed to get step runs: %v", err)
			return match
		}

		log.Printf("📦 Found %d step runs for workflow run %s", len(stepRuns), workflowRunID)

		var targetStepRun *entity.StepRun
		for i := range stepRuns {
			log.Printf("   - StepRun[%d]: StepID=%s, Status=%s, HasResponse=%v", i, stepRuns[i].StepID, stepRuns[i].Status, stepRuns[i].Response != "")
			if stepRuns[i].StepID == targetVariable.StepID {
				targetStepRun = &stepRuns[i]
				break
			}
		}

		if targetStepRun == nil {
			log.Printf("❌ No step run found for StepID %s", targetVariable.StepID)
			return match
		}

		if targetStepRun.Response == "" {
			log.Printf("❌ Step run %s has empty response", targetStepRun.ID)
			return match
		}

		log.Printf("📄 Step run response (first 200 chars): %s", truncateString(targetStepRun.Response, 200))

		value := u.extractValueFromResponse(targetStepRun.Response, targetVariable.Path)
		if value == "" {
			log.Printf("❌ Failed to extract value from response using path '%s'", targetVariable.Path)
			return match
		}

		log.Printf("✅ Successfully replaced '%s' with value (first 50 chars): %s", variableName, truncateString(value, 50))
		return value
	})
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func (u *ReplaceVariablesUseCase) extractValueFromResponse(response string, path string) string {
	var data interface{}
	if err := json.Unmarshal([]byte(response), &data); err != nil {
		return ""
	}

	// Split path by '.' but handle array indices like [0]
	current := data
	pathParts := splitPath(path)

	for _, part := range pathParts {
		// Check if this part is an array index like [0]
		if strings.HasPrefix(part, "[") && strings.HasSuffix(part, "]") {
			// Extract index
			indexStr := part[1 : len(part)-1]
			index := 0
			if _, err := fmt.Sscanf(indexStr, "%d", &index); err != nil {
				return ""
			}

			// Access array element
			switch v := current.(type) {
			case []interface{}:
				if index >= 0 && index < len(v) {
					current = v[index]
				} else {
					return ""
				}
			default:
				return ""
			}
		} else {
			// Normal object property
			switch v := current.(type) {
			case map[string]interface{}:
				var ok bool
				current, ok = v[part]
				if !ok {
					return ""
				}
			default:
				return ""
			}
		}
	}

	switch v := current.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%v", v)
	case bool:
		return fmt.Sprintf("%v", v)
	case nil:
		return ""
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		return string(jsonBytes)
	}
}

// splitPath splits a path like "data[0].user.name" into ["data", "[0]", "user", "name"]
func splitPath(path string) []string {
	var parts []string
	current := ""
	inBracket := false

	for i := 0; i < len(path); i++ {
		ch := path[i]

		if ch == '[' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
			inBracket = true
			current = "["
		} else if ch == ']' {
			current += "]"
			parts = append(parts, current)
			current = ""
			inBracket = false
		} else if ch == '.' && !inBracket {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts
}
