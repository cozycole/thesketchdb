package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ScopeActivation    = "activation"
	ScopePasswordReset = "password-reset"
	ScopeAPIKey        = "api-key"
)

type TokenModel struct {
	DB *pgxpool.Pool
}

type Token struct {
	ID        *int
	UserID    *int
	CreatedAt *time.Time
	Expiry    *time.Time
	Scope     *string
	Plaintext *string
	Hash      []byte
}

func generateToken() (plaintext string, hash []byte, err error) {
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return "", nil, err
	}
	plaintext = base64.URLEncoding.EncodeToString(b)
	sum := sha256.Sum256([]byte(plaintext))
	return plaintext, sum[:], nil
}

type TokenModelInterface interface {
	DeleteAllForUser(scope string, userID int) error
	Insert(token *Token) error
	New(userID int, ttl time.Duration, scope string) (*Token, error)
}

func (m *TokenModel) Insert(token *Token) error {
	args := []any{
		token.UserID,
		token.Scope,
		token.Hash,
		token.Expiry,
	}
	_, err := m.DB.Exec(context.Background(), `
        INSERT INTO token (user_id, scope, token_hash, ttl)
        VALUES ($1, $2, $3, $4)
    `, args...)
	if err != nil {
		return err
	}

	return nil
}
func (m TokenModel) DeleteAllForUser(scope string, userID int) error {
	query := `
	DELETE FROM token
	WHERE scope = $1 AND user_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.Exec(ctx, query, scope, userID)
	return err
}

func (m TokenModel) New(userID int, ttl time.Duration, scope string) (*Token, error) {
	plaintext, hash, err := generateToken()
	if err != nil {
		return nil, err
	}

	var ttlTime *time.Time
	if ttl > 0 {
		t := time.Now().Add(ttl)
		ttlTime = &t
	}

	token := &Token{
		UserID:    &userID,
		Scope:     &scope,
		Expiry:    ttlTime,
		Hash:      hash,
		Plaintext: &plaintext,
	}

	err = m.Insert(token)
	return token, err
}
