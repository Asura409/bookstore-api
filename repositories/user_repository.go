package repositories

import (
	"bookstore-api/models"
	
	"time"
	"gorm.io/gorm"
)
type UserRepository struct {
	DB *gorm.DB
	}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

//Get user by email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	if err := r.DB.Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	if err := r.DB.First(user, id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByUsername retrieves a user by their username from the database.
// It returns the User model and an error if any occurs."
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	if err := r.DB.Where("username = ?", username).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
//delete a user by ID
func (r *UserRepository) DeleteUserByID(ID int) error {
	user := &models.User{}
	return r.DB.Where("ID = ?", ID).Delete(user).Error
	
}

// UpdatePassword updates the user's password in the database.
// It takes the user ID and the new password as parameters.
func (r *UserRepository) UpdatePassword(userID int, newPassword string) error {
	
	return r.DB.Model(&models.User{}).Where("id = ?", userID).Update("password", newPassword).Error
}
	
// CreateResetToken creates a new password reset token in the database.
// It takes a PasswordResetToken model as a parameter and returns an error if any occurs.
func (r *UserRepository) CreateResetToken(token *models.PasswordResetToken) error {
	return r.DB.Create(token).Error
}

// FindValidToken checks if a password reset token is valid.
func (r *UserRepository) FindValidToken(token string) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	err := r.DB.Where("token = ? AND used = ? AND expires_at > ?", token, false, time.Now()).First(&resetToken).Error
	return &resetToken, err
}

// MarkTokenAsUsed marks a password reset token as used in the database.
func (r *UserRepository) MarkTokenAsUsed(tokenID uint) error {
	return r.DB.Model(&models.PasswordResetToken{}).Where("id = ?", tokenID).Update("used", true).Error
}