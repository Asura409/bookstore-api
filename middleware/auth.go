package middleware

import (
	
	"time"
	
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)

// Custom claims struct
type Claims struct {
    UserID   int    `json:"userId"`
    Username string `json:"username"`
    jwt.RegisteredClaims  // Embed standard claims (exp, iss, etc.)
}
// hash password
func HashPassword(password string) (string, error) {
	hash := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(hash, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
// compare password passwors and hash
func ComparePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}
// generate JWt
func GenerateToken(userID int, username string, secret string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),// Expires in 24 hours
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

