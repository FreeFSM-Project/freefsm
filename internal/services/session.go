package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionService struct {
	db *pgxpool.Pool
}

func NewSessionService(db *pgxpool.Pool) *SessionService {
	return &SessionService{db: db}
}

func (s *SessionService) Create(ctx context.Context, userID int64) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	token := hex.EncodeToString(b)
	hash := sha256.Sum256([]byte(token))

	_, err := s.db.Exec(ctx,
		`INSERT INTO sessions (token_hash, user_id, expires_at) VALUES ($1, $2, $3)`,
		hex.EncodeToString(hash[:]), userID, time.Now().Add(7*24*time.Hour),
	)
	if err != nil {
		return "", fmt.Errorf("save session: %w", err)
	}
	return token, nil
}

func (s *SessionService) Validate(ctx context.Context, token string) (int64, error) {
	if token == "" {
		return 0, fmt.Errorf("empty token")
	}
	hash := sha256.Sum256([]byte(token))
	var userID int64
	err := s.db.QueryRow(ctx,
		`SELECT user_id FROM sessions WHERE token_hash = $1 AND expires_at > NOW()`,
		hex.EncodeToString(hash[:]),
	).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("invalid session")
	}
	return userID, nil
}

func (s *SessionService) Delete(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	hash := sha256.Sum256([]byte(token))
	_, err := s.db.Exec(ctx, `DELETE FROM sessions WHERE token_hash = $1`, hex.EncodeToString(hash[:]))
	return err
}
