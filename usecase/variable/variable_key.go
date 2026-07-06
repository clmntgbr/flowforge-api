package variable

import (
	"flowforge-api/domain/entity"
	"regexp"
)

func stringUsesVariableKey(input, key string) bool {
	if input == "" {
		return false
	}

	pattern := regexp.MustCompile(`\{\{\s*` + regexp.QuoteMeta(key) + `\s*\}\}`)
	return pattern.MatchString(input)
}

func stepUsesVariableKey(step entity.Step, key string) bool {
	if stringUsesVariableKey(step.URL, key) {
		return true
	}

	for _, header := range step.Header {
		if stringUsesVariableKey(header.Value, key) {
			return true
		}
	}

	for _, query := range step.Query {
		if stringUsesVariableKey(query.Value, key) {
			return true
		}
	}

	if len(step.Body) > 0 && stringUsesVariableKey(string(step.Body), key) {
		return true
	}

	return false
}

func findStepsUsingVariableKey(steps []entity.Step, key string) []entity.Step {
	var usedSteps []entity.Step

	for _, step := range steps {
		if stepUsesVariableKey(step, key) {
			usedSteps = append(usedSteps, step)
		}
	}

	return usedSteps
}
