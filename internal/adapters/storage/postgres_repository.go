package storage

import (
	"context"
	"errors"
	"hexa-user/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) (*PostgresUserRepository, error) {
	if err := db.Ping(context.Background()); err != nil {
		return nil, errors.New("failed to connect to the database: " + err.Error())
	}

	return &PostgresUserRepository{
		db: db,
	}, nil
}

func (po *PostgresUserRepository) CreateTableIfNotExists(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		username VARCHAR(30) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := po.db.Exec(ctx, query)
	if err != nil {
		return errors.New("failed to create users table: " + err.Error())
	}

	return nil
}

func (po *PostgresUserRepository) Save(ctx context.Context, user *domain.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	// Example SQL query to insert a user, adjust according to your schema
	query := `INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)`
	_, err := po.db.Exec(ctx, query, user.ID, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (po *PostgresUserRepository) FindByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	query := `SELECT id, username, email, password FROM users WHERE id = $1`
	row := po.db.QueryRow(ctx, query, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (po *PostgresUserRepository) FindAll(ctx context.Context, page int, pageSize int) ([]*domain.User, error) {
	if page < 1 || pageSize < 1 {
		return nil, errors.New("invalid pagination parameters")
	}

	offset := (page - 1) * pageSize
	query := `SELECT id, username, email, password FROM users ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := po.db.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (po *PostgresUserRepository) Delete(ctx context.Context, id domain.UserID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := po.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
