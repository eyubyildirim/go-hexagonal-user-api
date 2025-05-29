package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserID string

type User struct {
	ID        UserID    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(username, email, password string) (*User, error) {
	if username == "" || len(username) < 3 || len(username) > 20 {
		return nil, errors.New("username must be between 3 and 20 characters")
	}

	if email == "" || password == "" {
		return nil, errors.New("email and password cannot be empty")
	}

	return &User{
		ID:        UserID(uuid.NewString()),
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
