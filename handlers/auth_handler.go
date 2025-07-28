package handlers

import (
	

	"github.com/gofiber/fiber/v2"
	"bookstore-api/services"
)

type AuthHandler struct {
	authService *services.AuthService
	emailService *services.EmailService
}

func NewAuthHandler(authService *services.AuthService, emailService *services.EmailService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		emailService: emailService,
	}
}

// Request password reset (POST /request-password-reset)
func (h *AuthHandler) RequestPasswordReset(c *fiber.Ctx) error {
	var request struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	if err := h.authService.RequestPasswordReset(request.Email); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process request",
		})
	}

	// In a real app, you'd send an email here
	// token := ... (retrieve from DB if needed)
	// h.emailService.SendResetLink(request.Email, token)

	return c.JSON(fiber.Map{
		"message": "If this email exists, a reset link has been sent",
	})
}

// Reset password (POST /reset-password)
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var request struct {
		Token       string `json:"token" validate:"required"`
		NewPassword string `json:"newPassword" validate:"required,min=8"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	if err := h.authService.ResetPassword(request.Token, request.NewPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password updated successfully",
	})
}