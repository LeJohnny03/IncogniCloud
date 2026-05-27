package models

import (
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Email       string    `db:"email" json:"email"`
	Username    string    `db:"username" json:"username"`
	DisplayName string    `db:"display_name" json:"display_name"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`

	Credentials []Credential `db:"-" json:"-"`
}

// WebAuthnID returns the user's ID as bytes
func (u *User) WebAuthnID() []byte {
	return []byte(u.ID.String())
}

// WebAuthnName returns the username for display
func (u *User) WebAuthnName() string {
	return u.Username
}

// WebAuthnDisplayName returns the display name
func (u *User) WebAuthnDisplayName() string {
	return u.DisplayName
}

// WebAuthnIcon returns the user icon URL (optional)
func (u *User) WebAuthnIcon() string {
	return ""
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {

	credentials := make([]webauthn.Credential, len(u.Credentials))

	for i, cred := range u.Credentials {
		credentials[i] = webauthn.Credential{
			ID:              cred.ID,
			PublicKey:       cred.PublicKey,
			AttestationType: cred.AttestationType,
			Authenticator: webauthn.Authenticator{
				AAGUID:       cred.AAGUID,
				SignCount:    uint32(cred.SignCount),
				CloneWarning: cred.CloneWarning,
			},
			Flags: webauthn.CredentialFlags{
				BackupEligible: cred.BackupEligible,
				BackupState:    cred.BackupState,
			},
		}
	}

	return credentials
}

type Credential struct {
	ID              []byte     `db:"id" json:"id"`
	UserID          uuid.UUID  `db:"user_id" json:"user_id"`
	PublicKey       []byte     `db:"public_key" json:"-"`
	AttestationType string     `db:"attestation_type" json:"attestation_type"`
	AAGUID          []byte     `db:"aaguid" json:"aaguid"`
	SignCount       int        `db:"sign_count" json:"sign_count"`
	CloneWarning    bool       `db:"clone_warning" json:"clone_warning"`
	BackupEligible  bool       `db:"backup_eligible" json:"backup_eligible"`
	BackupState     bool       `db:"backup_state" json:"backup_state"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	LastUsedAt      *time.Time `db:"last_used_at" json:"last_used_at,omitempty"`
}
