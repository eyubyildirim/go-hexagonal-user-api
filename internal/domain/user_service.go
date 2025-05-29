package domain

import "context"

type UserService interface {
	CreateUser(ctx context.Context, username, email, password string) (*User, error)
	UpdateUser(ctx context.Context, id UserID, username, email, password string) (*User, error)
	GetUserByID(ctx context.Context, id UserID) (*User, error)
	DeleteUser(ctx context.Context, id UserID) error
	ListUsers(ctx context.Context, page, pageSize int) ([]*User, error)
}
