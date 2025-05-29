package storage

import (
	"context"
	"errors"
	"hexa-user/internal/domain"
	"sync"
)

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[domain.UserID]*domain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[domain.UserID]*domain.User),
	}
}

func (in *InMemoryUserRepository) Save(ctx context.Context, user *domain.User) error {
	in.mu.Lock()
	defer in.mu.Unlock()

	if user == nil {
		return errors.New("user cannot be nil")
	}

	in.users[user.ID] = user

	return nil
}

func (in *InMemoryUserRepository) FindByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	in.mu.RLock()
	defer in.mu.RUnlock()

	user, ok := in.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}

	return user, nil
}
func (in *InMemoryUserRepository) FindAll(ctx context.Context, page int, pageSize int) ([]*domain.User, error) {
	in.mu.RLock()
	defer in.mu.RUnlock()

	if page < 1 || pageSize < 1 {
		return nil, errors.New("invalid pagination parameters")
	}

	start := (page - 1) * pageSize
	if start >= len(in.users) {
		return nil, nil // No users to return
	}

	end := min(start+pageSize, len(in.users))

	users := make([]*domain.User, 0, end-start)
	for _, user := range in.users {
		if len(users) >= end-start {
			break
		}
		users = append(users, user)
	}

	return users, nil
}

func (in *InMemoryUserRepository) Delete(ctx context.Context, id domain.UserID) error {
	in.mu.Lock()
	defer in.mu.Unlock()

	if _, ok := in.users[id]; !ok {
		return errors.New("user not found")
	}

	delete(in.users, id)
	return nil
}
