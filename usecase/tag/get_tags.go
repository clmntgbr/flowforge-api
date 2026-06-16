package tag

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type GetTagsUseCase struct {
	tagRepo *repository.TagRepository
}

func NewGetTagsUseCase(tagRepo *repository.TagRepository) *GetTagsUseCase {
	return &GetTagsUseCase{tagRepo: tagRepo}
}

func (u *GetTagsUseCase) Execute(ctx context.Context, organizationID uuid.UUID) ([]entity.Tag, error) {
	tags, err := (*u.tagRepo).List(ctx, organizationID)
	if err != nil {
		return []entity.Tag{}, err
	}

	return tags, nil
}
