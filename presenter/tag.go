package presenter

import (
	"flowforge-api/domain/entity"
	"time"
)

type TagResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewTagResponse(tags []entity.Tag) []TagResponse {
	responses := make([]TagResponse, len(tags))
	for i, tag := range tags {
		responses[i] = NewTagDetailResponse(tag)
	}
	return responses
}

func NewTagDetailResponse(tag entity.Tag) TagResponse {
	return TagResponse{
		ID:        tag.ID.String(),
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}
}
