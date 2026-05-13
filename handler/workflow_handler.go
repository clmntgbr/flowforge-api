package handler

import (
	"flowforge-api/handler/context"
	"flowforge-api/infrastructure/paginate"
	workflowDTO "flowforge-api/infrastructure/workflow"
	"flowforge-api/presenter"
	"flowforge-api/usecase/workflow"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type WorkflowHandler struct {
	listWorkflowsUseCase      *workflow.ListWorkflowsUseCase
	createWorkflowUseCase     *workflow.CreateWorkflowUseCase
	getWorkflowUseCase        *workflow.GetWorkflowUseCase
	updateWorkflowUseCase     *workflow.UpdateWorkflowUseCase
	activateWorkflowUseCase   *workflow.ActivateWorkflowUseCase
	deactivateWorkflowUseCase *workflow.DeactivateWorkflowUseCase
	upsertWorkflowUseCase     *workflow.UpsertWorkflowUseCase
}

func NewWorkflowHandler(
	listWorkflowsUseCase *workflow.ListWorkflowsUseCase,
	createWorkflowUseCase *workflow.CreateWorkflowUseCase,
	getWorkflowUseCase *workflow.GetWorkflowUseCase,
	updateWorkflowUseCase *workflow.UpdateWorkflowUseCase,
	activateWorkflowUseCase *workflow.ActivateWorkflowUseCase,
	deactivateWorkflowUseCase *workflow.DeactivateWorkflowUseCase,
	upsertWorkflowUseCase *workflow.UpsertWorkflowUseCase,
) *WorkflowHandler {
	return &WorkflowHandler{
		listWorkflowsUseCase:      listWorkflowsUseCase,
		createWorkflowUseCase:     createWorkflowUseCase,
		getWorkflowUseCase:        getWorkflowUseCase,
		updateWorkflowUseCase:     updateWorkflowUseCase,
		activateWorkflowUseCase:   activateWorkflowUseCase,
		deactivateWorkflowUseCase: deactivateWorkflowUseCase,
		upsertWorkflowUseCase:     upsertWorkflowUseCase,
	}
}

func (h *WorkflowHandler) GetWorkflows(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var query paginate.PaginateQuery
	if err := c.Bind().Query(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}
	query.Normalize()

	workflows, total, err := h.listWorkflowsUseCase.Execute(c.Context(), activeOrganizationID, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	return c.JSON(paginate.NewPaginateResponse(presenter.NewWorkflowListResponses(workflows), int(total), query))
}

func (h *WorkflowHandler) GetWorkflowByID(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	workflowID := c.Params("id")
	workflowUUID, err := uuid.Parse(workflowID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workflow ID",
		})
	}

	workflow, err := h.getWorkflowUseCase.Execute(c.Context(), activeOrganizationID, workflowUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get workflow",
		})
	}

	return c.JSON(presenter.NewWorkflowDetailResponse(workflow))
}

func (h *WorkflowHandler) CreateWorkflow(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var request workflowDTO.CreateWorkflowInput
	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.New().Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}

	_, err = h.createWorkflowUseCase.Execute(c.Context(), activeOrganizationID, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create workflow",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}

func (h *WorkflowHandler) UpdateWorkflow(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	workflowID := c.Params("id")
	workflowUUID, err := uuid.Parse(workflowID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workflow ID",
		})
	}

	var request workflowDTO.UpdateWorkflowInput
	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.New().Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}

	_, err = h.updateWorkflowUseCase.Execute(c.Context(), activeOrganizationID, workflowUUID, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update workflow",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}

func (h *WorkflowHandler) ActivateWorkflow(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	workflowID := c.Params("id")
	workflowUUID, err := uuid.Parse(workflowID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workflow ID",
		})
	}

	_, err = h.activateWorkflowUseCase.Execute(c.Context(), activeOrganizationID, workflowUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to activate workflow",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}

func (h *WorkflowHandler) DeactivateWorkflow(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	workflowID := c.Params("id")
	workflowUUID, err := uuid.Parse(workflowID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workflow ID",
		})
	}

	_, err = h.deactivateWorkflowUseCase.Execute(c.Context(), activeOrganizationID, workflowUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to deactivate workflow",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}

func (h *WorkflowHandler) UpsertWorkflow(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	workflowID := c.Params("id")
	workflowUUID, err := uuid.Parse(workflowID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workflow ID",
		})
	}

	var request workflowDTO.UpsertWorkflowInput
	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.New().Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}

	_, err = h.upsertWorkflowUseCase.Execute(c.Context(), activeOrganizationID, workflowUUID, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to upsert workflow",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
