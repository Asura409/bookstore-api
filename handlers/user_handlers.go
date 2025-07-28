package handlers

import (
	"bookstore-api/middleware"
	"bookstore-api/models"
	"bookstore-api/repositories"
	"time"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	repo *repositories.UserRepository
}

func NewUserHandler(repo *repositories.UserRepository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

func (h *UserHandler) CreateUserHandler(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	hash, crypterr := middleware.HashPassword(user.Password)
	if crypterr!= nil{
		c.Status(fiber.StatusInternalServerError)
	}
	user.Password = hash
	if err := h.repo.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})	
	}	
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) LoginUserHandler(c *fiber.Ctx) error{
	username := c.FormValue("Username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error" : "email and password required!!"})
	}
	user, err := h.repo.GetUserByUsername(username)
	if	 err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user"})
}
if check := middleware.ComparePassword(user.Password, password); check != nil{
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error" : "Invalid Credentials",
	})
}
//GENERATE TOKEN AND STORE IN COOKIE
token,terr := middleware.GenerateToken(user.ID, username,"secret key")
if terr != nil {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get token"})
}
 c.Cookie(&fiber.Cookie{
	Name : "jwt",
	Value : token,
	Expires: time.Now().Add(time.Minute*30),
	HTTPOnly: true,
	Secure : true,
	SameSite : "Lax",
	Path : "/",
 })
 
return c.JSON(fiber.Map{
	"message" : "Login Successful",
})

}



func (h *UserHandler) GetUserByIDHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	user, err := h.repo.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}

func (h *UserHandler) GetUserByUsernameHandler(c *fiber.Ctx) error {
	username := c.Params("username")

	user, err := h.repo.GetUserByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}


type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func(h *UserHandler) ResetPasswordHandler(c *fiber.Ctx) error {
	
		// Parse request
		req := new(ForgotPasswordRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}

		// Validate email format
		if !isValidEmail(req.Email) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Please enter a valid email address",
			}) // :cite[5]
		}

		// Check if user exists
		var user models.User
		if err := h.repo.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Don't reveal whether user exists for security
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"message": "If the email exists, a reset link has been sent",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}

		// Generate reset token (JWT)
		token, err := middleware.GenerateToken(user.ID, user.Username, "reset-password-secret")

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate token",
			})
		}

		

		// Store token in database (or cache)
		resetToken := models.PasswordResetToken{
			Email:     user.Email,
			Token:     token,
			ExpiresAt: time.Now().Add(time.Hour * 1),
		}
		if err := h.repo.DB.Create(&resetToken).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save token",
			})
		}

		// Send email (in production you'd use a mail service)
		go sendResetEmail(user.Email, token)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "If the email exists, a reset link has been sent",
		})
	}


func isValidEmail(email string) bool {
	// Simple email validation - implement proper validation
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func sendResetEmail(to, token string) {
	// Implement actual email sending logic
	// This would typically use SMTP or a service like SendGrid
	// For testing, you might just log it
	fmt.Printf("Reset link for %s: http://yourapp.com/reset-password?token=%s\n", to, token)
}