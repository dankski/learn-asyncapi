package store

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq" // Postgres driver
)

type UserStore struct {
	db *sqlx.DB
}

type User struct {
	Id                   uuid.UUID `db:"id"`
	Username             string    `db:"username"`
	Email                string    `db:"email"`
	HashedPasswordBase64 string    `db:"hash_password"`
	CreatedAt            time.Time `db:"created_at"`
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: sqlx.NewDb(db, "postgres"),
	}
}

func (u *User) ComparePassword(password string) error {
	hashedPassword, err := base64.StdEncoding.DecodeString(u.HashedPasswordBase64)
	if err != nil {
		return fmt.Errorf("failed to decode hashed password: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	return nil
}

func (us *UserStore) ByEmail(ctx context.Context, email string) (*User, error) {
	const query = `
		SELECT id, username, email, hash_password, created_at
		FROM users
		WHERE email = $1
	`

	var user User
	err := us.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (us *UserStore) ById(ctx context.Context, userId uuid.UUID) (*User, error) {
	const query = `
		SELECT id, username, email, hash_password, created_at
		FROM users
		WHERE id = $1
	`

	var user User
	err := us.db.GetContext(ctx, &user, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (us *UserStore) CreateUser(ctx context.Context, username, email, password string) (*User, error) {
	const dml = `
		INSERT INTO users (username, email, hash_password)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, hash_password, created_at
	`

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	HashedPasswordBase64 := base64.StdEncoding.EncodeToString(bytes)

	var user User
	err = us.db.GetContext(ctx, &user, dml, username, email, HashedPasswordBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to insert the user: %w", err)
	}

	return &user, nil
}
