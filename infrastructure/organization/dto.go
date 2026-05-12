package organization

type CreateOrganizationInput struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}

type UpdateOrganizationInput struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}
