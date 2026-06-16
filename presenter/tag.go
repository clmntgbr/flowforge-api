package presenter

import (
	"flowforge-api/domain/entity"
)

type TagResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func NewTagResponses(tags []entity.Tag) []TagResponse {
	responses := make([]TagResponse, len(tags))
	for i, tag := range tags {
		responses[i] = NewTagDetailResponse(tag)
	}
	return responses
}

func NewTagDetailResponse(tag entity.Tag) TagResponse {
	return TagResponse{
		ID:    tag.ID.String(),
		Name:  tag.Name,
		Color: tag.Color,
	}
}
