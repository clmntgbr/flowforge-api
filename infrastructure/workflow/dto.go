package workflow

type CreateWorkflowInput struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description" validate:"omitempty,min=2,max=255"`
}

type UpdateWorkflowInput struct {
	CreateWorkflowInput
}
