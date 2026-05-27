package config

import (
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

type WebAuthnConfig struct {
	RPDisplayName string
	RPID          string
	RPOrigins     []string
}

func NewWebAuthn(cfg WebAuthnConfig) (*webauthn.WebAuthn, error) {

	wconfig := &webauthn.Config{
		RPDisplayName: cfg.RPDisplayName,
		RPID:          cfg.RPID,
		RPOrigins:     cfg.RPOrigins,

		AuthenticatorSelection: protocol.AuthenticatorSelection{
			RequireResidentKey: protocol.ResidentKeyNotRequired(),
			ResidentKey:        protocol.ResidentKeyRequirementPreferred,
			UserVerification:   protocol.VerificationPreferred,
		},

		AttestationPreference: protocol.PreferNoAttestation,

		Timeouts: webauthn.TimeoutsConfig{
			Login: webauthn.TimeoutConfig{
				Enforce: true,
				Timeout: 60000,
			},
			Registration: webauthn.TimeoutConfig{
				Enforce: true,
				Timeout: 60000,
			},
		},
	}

	return webauthn.New(wconfig)

}
