package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"bookstore-api/models"
	"bookstore-api/repositories"
)

type AuthService struct {
	userRepo      *repositories.UserRepository
	jwtSecret     string
	resetTokenExp time.Duration
}

func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		resetTokenExp: 1 * time.Hour, // Tokens expire in 1 hour
	}
}

// Generate a JWT token for password reset
func (s *AuthService) GenerateResetToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"exp":      time.Now().Add(s.resetTokenExp).Unix(),
		"purpose":  "password_reset",
	})

	return token.SignedString([]byte(s.jwtSecret))
}

// Validate a reset token
func (s *AuthService) ValidateResetToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["purpose"] != "password_reset" {
			return 0, errors.New("invalid token purpose")
		}
		return uint(claims["user_id"].(float64)), nil
	}

	return 0, errors.New("invalid token")
}

// Hash a password
func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Initiate password reset
func (s *AuthService) RequestPasswordReset(email string) error {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		// Don't reveal if user exists (security)
		return nil
	}

	token, err := s.GenerateResetToken(uint(user.ID))
	if err != nil {
		return err
	}

	resetToken := &models.PasswordResetToken{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(s.resetTokenExp),
	}

	return s.userRepo.CreateResetToken(resetToken)
}

// Complete password reset
func (s *AuthService) ResetPassword(token, newPassword string) error {
	// Validate token structure first
	userID, err := s.ValidateResetToken(token)
	if err != nil {
		return errors.New("invalid token")
	}

	// Check if token exists in DB and is unused
	resetToken, err := s.userRepo.FindValidToken(token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	// Hash new password
	hashedPassword, err := s.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Update user password
	if err := s.userRepo.UpdatePassword(int(userID), hashedPassword); err != nil {
		return errors.New("failed to update password")
	}

	// Mark token as used
	return s.userRepo.MarkTokenAsUsed(uint(resetToken.ID))
}
