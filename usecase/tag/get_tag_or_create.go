package tag

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetTagOrCreateUseCase struct {
	tagRepo *repository.TagRepository
}

func NewGetTagOrCreateUseCase(tagRepo *repository.TagRepository) *GetTagOrCreateUseCase {
	return &GetTagOrCreateUseCase{tagRepo: tagRepo}
}

func (u *GetTagOrCreateUseCase) Execute(ctx context.Context, organizationID uuid.UUID, tagID uuid.UUID, tagName string, tagColor string) (entity.Tag, error) {
	tag, err := (*u.tagRepo).Get(ctx, organizationID, tagID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return entity.Tag{}, err
	}

	if err == gorm.ErrRecordNotFound {
		tag = entity.Tag{
			ID:             tagID,
			OrganizationID: organizationID,
			Name:           tagName,
			Color:          tagColor,
		}
		err = (*u.tagRepo).Create(ctx, &tag)
	}

	return tag, nil
}
