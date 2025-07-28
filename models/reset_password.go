package models

import "time"



type PasswordResetToken struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time`json:"expires_at"`   
	Used 	  bool    `json:"used"`
}
