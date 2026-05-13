package step

import (
	"context"
	"strconv"
	"strings"
)

type CalculateExecutionOrderUseCase struct{}

func NewCalculateExecutionOrderUseCase() *CalculateExecutionOrderUseCase {
	return &CalculateExecutionOrderUseCase{}
}

func (u *CalculateExecutionOrderUseCase) Execute(ctx context.Context, index string) int {
	parts := strings.Split(index, ".")

	if len(parts) == 1 {
		level, _ := strconv.Atoi(parts[0])
		return level
	}

	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	return major*1000 + minor
}
