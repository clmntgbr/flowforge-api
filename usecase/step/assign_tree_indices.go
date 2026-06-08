package step

import (
	"context"
	"flowforge-api/domain/repository"
	"sort"

	"github.com/google/uuid"
)

type AssignTreeIndicesUseCase struct {
	stepRepo      *repository.StepRepository
	connexionRepo *repository.ConnexionRepository
}

func NewAssignTreeIndicesUseCase(stepRepo *repository.StepRepository, connexionRepo *repository.ConnexionRepository) *AssignTreeIndicesUseCase {
	return &AssignTreeIndicesUseCase{stepRepo: stepRepo, connexionRepo: connexionRepo}
}

// Execute computes connected components of the workflow step graph (undirected)
// and assigns a 1-based tree_index to each step, sorted by the minimum
// execution_order within each component so that the "first" tree runs first.
func (u *AssignTreeIndicesUseCase) Execute(ctx context.Context, workflowID uuid.UUID) error {
	steps, err := (*u.stepRepo).GetByWorkflowID(ctx, workflowID)
	if err != nil {
		return err
	}
	if len(steps) == 0 {
		return nil
	}

	connexions, err := (*u.connexionRepo).GetByWorkflowID(ctx, workflowID)
	if err != nil {
		return err
	}

	parent := make(map[uuid.UUID]uuid.UUID, len(steps))
	for _, s := range steps {
		parent[s.ID] = s.ID
	}

	var find func(uuid.UUID) uuid.UUID
	find = func(x uuid.UUID) uuid.UUID {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}

	union := func(x, y uuid.UUID) {
		rx, ry := find(x), find(y)
		if rx != ry {
			parent[rx] = ry
		}
	}

	for _, c := range connexions {
		if _, ok := parent[c.FromStepID]; !ok {
			continue
		}
		if _, ok := parent[c.ToStepID]; !ok {
			continue
		}
		union(c.FromStepID, c.ToStepID)
	}

	// Group steps by root and track minimum execution_order per component
	type component struct {
		root     uuid.UUID
		minOrder int
		stepIDs  []uuid.UUID
	}
	comps := make(map[uuid.UUID]*component)
	for _, s := range steps {
		root := find(s.ID)
		comp, exists := comps[root]
		if !exists {
			comps[root] = &component{root: root, minOrder: s.ExecutionOrder, stepIDs: []uuid.UUID{s.ID}}
		} else {
			comp.stepIDs = append(comp.stepIDs, s.ID)
			if s.ExecutionOrder < comp.minOrder {
				comp.minOrder = s.ExecutionOrder
			}
		}
	}

	sorted := make([]*component, 0, len(comps))
	for _, c := range comps {
		sorted = append(sorted, c)
	}
	sort.Slice(sorted, func(i, j int) bool {
		si, sj := len(sorted[i].stepIDs), len(sorted[j].stepIDs)
		if si != sj {
			return si > sj
		}
		return sorted[i].minOrder < sorted[j].minOrder
	})

	indices := make(map[uuid.UUID]int, len(steps))
	for treeIdx, comp := range sorted {
		for _, id := range comp.stepIDs {
			indices[id] = treeIdx + 1
		}
	}

	return (*u.stepRepo).UpdateTreeIndices(ctx, indices)
}
