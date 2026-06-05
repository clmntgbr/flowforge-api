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

const indexBase = 100
const indexMaxDepth = 4

func (u *CalculateExecutionOrderUseCase) Execute(ctx context.Context, index string) int {
	parts := strings.Split(index, ".")
	result := 0
	for _, part := range parts {
		val, _ := strconv.Atoi(part)
		result = result*indexBase + val
	}
	for i := len(parts); i < indexMaxDepth; i++ {
		result *= indexBase
	}
	return result
}
