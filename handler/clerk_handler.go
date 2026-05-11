package handler

import (
	"encoding/json"
	"flowforge-api/domain/repository"
	"flowforge-api/presenter"
	"flowforge-api/usecase/organization"
	"flowforge-api/usecase/user"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClerkHandler struct {
	userRepository            repository.UserRepository
	createUserUseCase         *user.CreateUserUseCase
	createOrganizationUseCase *organization.CreateOrganizationUseCase
	updateUserUseCase         *user.UpdateUserUseCase
	deleteUserUseCase         *user.DeleteUserUseCase
}

func NewClerkHandler(userRepository repository.UserRepository, createUserUseCase *user.CreateUserUseCase, createOrganizationUseCase *organization.CreateOrganizationUseCase, updateUserUseCase *user.UpdateUserUseCase, deleteUserUseCase *user.DeleteUserUseCase) *ClerkHandler {
	return &ClerkHandler{
		userRepository:            userRepository,
		createUserUseCase:         createUserUseCase,
		createOrganizationUseCase: createOrganizationUseCase,
		updateUserUseCase:         updateUserUseCase,
		deleteUserUseCase:         deleteUserUseCase,
	}
}

func (h *ClerkHandler) Execute(c fiber.Ctx) error {
	clerkEvent := c.Locals("payload").(presenter.ClerkEvent)
	validate := validator.New()

	fmt.Printf("Clerk event type: %s", clerkEvent.Type)

	switch clerkEvent.Type {
	case "user.created":
		var data presenter.ClerkUserCreated
		if err := json.Unmarshal(clerkEvent.Data, &data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid request body",
			})
		}

		if err := validate.Struct(data); err != nil {
			return err
		}

		if err := h.CreateUser(c, data); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to create user",
			})
		}

		return c.SendStatus(fiber.StatusCreated)

	case "user.updated":
		var data presenter.ClerkUserUpdated
		if err := json.Unmarshal(clerkEvent.Data, &data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid request body",
			})
		}

		if err := validate.Struct(data); err != nil {
			return err
		}

		if err := h.UpdateUser(c, data); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to update user",
			})
		}

		return c.SendStatus(fiber.StatusNoContent)

	case "user.deleted":
		var data presenter.ClerkUserDeleted
		if err := json.Unmarshal(clerkEvent.Data, &data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid request body",
			})
		}

		if err := validate.Struct(data); err != nil {
			return err
		}

		if err := h.DeleteUser(c, data); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to delete user",
			})
		}

		return c.SendStatus(fiber.StatusNoContent)

	default:
		log.Printf("Unhandled event type: %s", clerkEvent.Type)
		return c.SendStatus(fiber.StatusOK)
	}
}

func (h *ClerkHandler) CreateUser(c fiber.Ctx, data presenter.ClerkUserCreated) error {
	user, err := h.userRepository.GetByClerkID(c.Context(), data.ID)
	if err != nil {
		log.Printf("Error finding user by Clerk ID %s: %v", data.ID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
		})
	}

	if user != nil {
		log.Printf("User with Clerk ID %s already exists, skipping creation", data.ID)
		return nil
	}

	txFunc := func(tx *gorm.DB) error {
		user, err = h.createUserUseCase.Execute(c.Context(), data.ID, data.FirstName, data.LastName, *data.Banned)
		if err != nil {
			log.Printf("Error creating user with Clerk ID %s: %v", data.ID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to create user",
			})
		}

		organization, err := h.createOrganizationUseCase.Execute(c.Context(), user, "Default Organization")
		if err != nil {
			log.Printf("Error creating default organization for user %s: %v", user.ID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to create default organization",
			})
		}

		organizationID, err := uuid.Parse(organization.ID)
		if err != nil {
			log.Printf("Error parsing organization UUID %s: %v", organization.ID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to parse organization UUID",
			})
		}

		user.ActiveOrganizationID = &organizationID
		if err := h.userRepository.Update(c.Context(), user); err != nil {
			log.Printf("Error updating user %s with active organization: %v", user.ID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to update user",
			})
		}

		log.Printf("Successfully created user with Clerk ID %s and organization %s", data.ID, organizationID)
		return nil
	}

	return txFunc(nil)
}

func (h *ClerkHandler) UpdateUser(c fiber.Ctx, data presenter.ClerkUserUpdated) error {
	user, err := h.userRepository.GetByClerkID(c.Context(), data.ID)
	if err != nil {
		log.Printf("Error finding user by Clerk ID %s: %v", data.ID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to find user",
		})
	}

	if user == nil {
		log.Printf("User with Clerk ID %s not found for update", data.ID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	user.FirstName = data.FirstName
	user.LastName = data.LastName
	user.Banned = *data.Banned

	if err := h.updateUserUseCase.Execute(c.Context(), user); err != nil {
		log.Printf("Error updating user with Clerk ID %s: %v", data.ID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update user",
		})
	}

	log.Printf("Successfully updated user with Clerk ID %s", data.ID)
	return nil
}

func (h *ClerkHandler) DeleteUser(c fiber.Ctx, data presenter.ClerkUserDeleted) error {
	if err := h.deleteUserUseCase.Execute(c.Context(), data.ID); err != nil {
		log.Printf("Error deleting user with Clerk ID %s: %v", data.ID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete user",
		})
	}

	log.Printf("Successfully deleted user with Clerk ID %s", data.ID)
	return nil
}
