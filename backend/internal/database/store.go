package database

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
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
        SELECT id, user_id, public_key, attestation_type, aaguid, sign_count, clone_warning, created_at, last_used_at
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
        INSERT INTO credentials (id, user_id, public_key, attestation_type, aaguid, sign_count, clone_warning)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	_, err := s.db.ExecContext(ctx, query,
		cred.ID,
		cred.UserID,
		cred.PublicKey,
		cred.AttestationType,
		cred.AAGUID,
		cred.SignCount,
		cred.CloneWarning,
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
        SELECT id, user_id, public_key, attestation_type, aaguid, sign_count, clone_warning, created_at, last_used_at
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
