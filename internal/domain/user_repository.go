package domain

import "context"

type UserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id UserID) (*User, error)
	FindAll(ctx context.Context, page, pageSize int) ([]*User, error)
	Delete(ctx context.Context, id UserID) error
}
