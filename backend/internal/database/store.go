package database

import (
	"backend/internal/models"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

type AuthSession struct {
	SessionToken string
	UserID       uuid.UUID
	ExpiresAt    time.Time
}

func GenerateSecureToken(length int) (string, error) {

	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil

}

func (s *Store) GetSession(ctx context.Context, token string) (*AuthSession, error) {
	query := `SELECT session_token, user_id, expires_at FROM auth_sessions WHERE session_token = $1`

	var session AuthSession
	err := s.db.QueryRowContext(ctx, query, token).Scan(
		&session.SessionToken,
		&session.UserID,
		&session.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *Store) CreateSession(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	query := `INSERT INTO auth_sessions (user_id, session_token, expires_at) VALUES ($1, $2, $3)`
	_, err := s.db.ExecContext(ctx, query, userID, token, expiresAt)
	return err
}

func (s *Store) DeleteSession(ctx context.Context, token string) error {
	query := `DELETE FROM auth_sessions WHERE session_token = $1`
	_, err := s.db.ExecContext(ctx, query, token)
	return err
}

func (s *Store) CreateUser(ctx context.Context, username, displayName string) (*models.User, error) {
	user := &models.User{
		ID:          uuid.New(),
		Username:    username,
		DisplayName: displayName,
	}

	query := `
        INSERT INTO users (id, username, display_name)
        VALUES ($1, $2, $3)
        RETURNING created_at, updated_at
    `

	err := s.db.QueryRowContext(ctx, query, user.ID, user.Username, user.DisplayName).
		Scan(&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, display_name, created_at, updated_at FROM users WHERE username = $1`

	err := s.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (s *Store) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, display_name, created_at, updated_at FROM users WHERE id = $1`

	err := s.db.GetContext(ctx, &user, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (s *Store) GetUserCredentials(ctx context.Context, userID uuid.UUID) ([]models.Credential, error) {
	var credentials []models.Credential
	query := `
        SELECT id, user_id, public_key, attestation_type, aaguid, sign_count, clone_warning, backup_eligible, backup_state, created_at, last_used_at
        FROM credentials
        WHERE user_id = $1
        ORDER BY created_at DESC
    `

	err := s.db.SelectContext(ctx, &credentials, query, userID)
	if err != nil {
		return nil, err
	}

	return credentials, nil
}

func (s *Store) AddCredential(ctx context.Context, cred *models.Credential) error {
	query := `
        INSERT INTO credentials (id, user_id, public_key, attestation_type, aaguid, sign_count, clone_warning, backup_eligible, backup_state)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `

	_, err := s.db.ExecContext(ctx, query,
		cred.ID,
		cred.UserID,
		cred.PublicKey,
		cred.AttestationType,
		cred.AAGUID,
		cred.SignCount,
		cred.CloneWarning,
		cred.BackupEligible,
		cred.BackupState,
	)

	return err
}

func (s *Store) UpdateCredential(ctx context.Context, credID []byte, signCount int, cloneWarning bool) error {
	query := `
        UPDATE credentials
        SET sign_count = $2, clone_warning = $3, last_used_at = CURRENT_TIMESTAMP
        WHERE id = $1
    `

	_, err := s.db.ExecContext(ctx, query, credID, signCount, cloneWarning)
	return err
}

func (s *Store) GetCredentialByID(ctx context.Context, credID []byte) (*models.Credential, error) {
	var cred models.Credential
	query := `
        SELECT id, user_id, public_key, attestation_type, aaguid, sign_count, clone_warning, backup_eligible, backup_state, created_at, last_used_at
        FROM credentials
        WHERE id = $1
    `

	err := s.db.GetContext(ctx, &cred, query, credID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &cred, nil
}

func (s *Store) IsSystemSetup(ctx context.Context) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users)`
	err := s.db.QueryRowContext(ctx, query).Scan(&exists)
	return exists, err
}
