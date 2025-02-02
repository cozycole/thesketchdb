package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *pgxpool.Pool
}

type password struct {
	plaintext *string
	hash      []byte
}

type User struct {
	ID        int
	CreatedAt time.Time
	Username  string
	Email     string
	Password  password
	Activated bool
	Role      string
}

func (p *password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintext
	p.hash = hash

	return nil
}

func (p *password) Match(plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintext))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

type UserModelInterface interface {
	AddLike(userId, videoId int) error
	Authenticate(username, password string) (int, error)
	GetByUsername(username string) (*User, error)
	GetById(id int) (*User, error)
	Insert(user *User) error
	RemoveLike(userId, videoId int) error
}

func (m *UserModel) AddLike(userId, videoId int) error {
	stmt := `
		INSERT INTO likes (user_id, video_id)
		VALUES($1, $2)
	`

	_, err := m.DB.Exec(context.Background(), stmt, userId, videoId)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(username, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, password_hash FROM users WHERE username = $1"

	err := m.DB.QueryRow(context.Background(), stmt, username).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserModel) GetByUsername(username string) (*User, error) {
	query := `
		SELECT id, created_at, username, email, password_hash, activated, role
		FROM users
		WHERE username = $1
	`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Role,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) GetById(id int) (*User, error) {
	query := `
		SELECT id, created_at, username, email, password_hash, activated, role
		FROM users
		WHERE id = $1
	`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Role,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (username, email, password_hash, activated)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	args := []interface{}{
		user.Username,
		user.Email,
		user.Password.hash,
		user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, args...).Scan(
		&user.ID, &user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m *UserModel) RemoveLike(userId, videoId int) error {
	stmt := `
		DELETE FROM likes 
		WHERE user_id = $1 AND video_id = $2	
	`

	_, err := m.DB.Exec(context.Background(), stmt, userId, videoId)
	if err != nil {
		return err
	}
	return nil
}
