package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenModel struct {
	DB *pgxpool.Pool
}

type Token struct {
	ID         int
	UserID     int
	CreatedAt  time.Time
	LastUsedAt *time.Time
	secret     password // reusing your existing struct
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
	Insert(userId int) (string, error)
}

func (m *TokenModel) Insert(userID int) (string, error) {
	plaintext, hash, err := generateToken()
	if err != nil {
		return "", err
	}

	_, err = m.DB.Exec(context.Background(), `
        INSERT INTO token (user_id, token_hash)
        VALUES ($1, $2)
    `, userID, hash)
	if err != nil {
		return "", err
	}

	return plaintext, nil
}
