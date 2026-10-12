package models

import (
	"context"
	"crypto/sha256"
	"errors"
	"strings"
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
	ID           *int
	CreatedAt    *time.Time
	Username     *string
	Email        *string
	Password     password
	Activated    *bool
	Role         *string
	ProfileImage *string
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

type UserSketchInfo struct {
	Rating    *int
	Timestamp *time.Time
}

type UserModelInterface interface {
	AddRating(userId, sketchId, rating int) error
	AddLike(userId, sketchId int) error
	Authenticate(username, password string) (int, error)
	DeleteRating(userId, sketchId int) error
	GetById(id int) (*User, error)
	GetByToken(plaintext, scope string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	GetUserSketchInfo(userId, sketchId int) (*UserSketchInfo, error)
	Insert(user *User) error
	RemoveLike(userId, sketchId int) error
	Update(user *User) error
	UpdateRating(userId, sketchId, rating int) error
}

func (m *UserModel) AddLike(userId, sketchId int) error {
	stmt := `
		INSERT INTO likes (user_id, sketch_id)
		VALUES($1, $2)
	`

	_, err := m.DB.Exec(context.Background(), stmt, userId, sketchId)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) AddRating(userId, sketchId, rating int) error {
	stmt := `
		INSERT INTO sketch_rating (user_id, sketch_id, rating)
		VALUES($1, $2, $3)
	`

	_, err := m.DB.Exec(context.Background(), stmt, userId, sketchId, rating)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) GetUserSketchInfo(userId, sketchId int) (*UserSketchInfo, error) {
	stmt := `
		SELECT rating, created_at 
		FROM sketch_rating 
		WHERE user_id = $1 AND sketch_id = $2
	`

	row := m.DB.QueryRow(context.Background(), stmt, userId, sketchId)

	r := UserSketchInfo{}

	if err := row.Scan(&r.Rating, &r.Timestamp); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &r, ErrNoRecord
		}
		return &r, err
	}
	return &r, nil
}

func (m *UserModel) UpdateRating(userId, sketchId, rating int) error {
	stmt := `
		UPDATE sketch_rating SET rating = $1
		WHERE user_id = $2 AND sketch_id = $3
	`

	_, err := m.DB.Exec(context.Background(), stmt, rating, userId, sketchId)
	if err != nil {
		return err
	}
	return nil

}

func (m *UserModel) DeleteRating(userId, sketchId int) error {
	stmt := `
		DELETE FROM sketch_rating
		WHERE user_id = $1 AND sketch_id = $2
	`

	_, err := m.DB.Exec(context.Background(), stmt, userId, sketchId)
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

func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, created_at, username, email, password_hash, activated, role, profile_image
		FROM users
		WHERE email = $1
	`

	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	password := password{}
	err := m.DB.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&password.hash,
		&user.Activated,
		&user.Role,
		&user.ProfileImage,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	user.Password = password

	return &user, nil
}

func (m *UserModel) GetByUsername(username string) (*User, error) {
	query := `
		SELECT id, created_at, username, email, password_hash, activated, role, profile_image
		FROM users
		WHERE username = $1
	`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	password := password{}
	err := m.DB.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&password.hash,
		&user.Activated,
		&user.Role,
		&user.ProfileImage,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	user.Password = password

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

	password := password{}
	err := m.DB.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&password.hash,
		&user.Activated,
		&user.Role,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	user.Password = password

	return &user, nil
}

func (m *UserModel) GetByToken(plaintext, scope string) (*User, error) {
	sum := sha256.Sum256([]byte(plaintext))
	hash := sum[:]

	query := `
		SELECT u.id, u.created_at, u.username, u.email, u.password_hash, u.activated, u.role
		FROM users as u
		JOIN token as t
		ON u.id = t.user_id
		WHERE t.token_hash = $1
			AND t.scope = $2
			AND (t.ttl IS NULL OR t.ttl > $3)
	`

	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	password := password{}
	err := m.DB.QueryRow(ctx, query, hash, scope, time.Now()).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&password.hash,
		&user.Activated,
		&user.Role,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	user.Password = password

	return &user, nil
}

func (m *UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	args := []any{
		user.Username,
		user.Email,
		user.Password.hash,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, args...).Scan(
		&user.ID, &user.CreatedAt)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), `violates unique constraint "users_email_key"`):
			return ErrDuplicateEmail
		case strings.Contains(err.Error(), `violates unique constraint "users_username_key"`):
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (m *UserModel) RemoveLike(userId, sketchId int) error {
	stmt := `
		DELETE FROM likes 
		WHERE user_id = $1 AND sketch_id = $2	
	`

	_, err := m.DB.Exec(context.Background(), stmt, userId, sketchId)
	if err != nil {
		return err
	}
	return nil
}

func (m UserModel) Update(user *User) error {
	stmt := `
		UPDATE users
		SET username = $1, email = $2, password_hash = $3, activated = $4, role = $5,
		profile_image = $6
		WHERE id = $7
		`
	args := []any{
		user.Username,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.Role,
		user.ProfileImage,
		user.ID,
	}

	_, err := m.DB.Exec(context.Background(), stmt, args...)
	if err != nil {
		return err
	}

	return nil
}
