package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// Das ist dein echtes User-Abbild aus der Datenbank
type User struct {
	ID       int
	Username string
}

// Und hier hängen wir die Funktionen an das Struct an, damit WebAuthn zufrieden ist:
func (u User) WebAuthnID() []byte {
	// WebAuthn braucht die ID als Byte-Array, nicht als Zahl.
	// Später wandeln wir hier u.ID in Bytes um.
	return []byte(fmt.Sprintf("%d", u.ID))
}

func (u User) WebAuthnName() string {
	return u.Username
}

func (u User) WebAuthnDisplayName() string {
	return u.Username
}

func (u User) WebAuthnCredentials() []webauthn.Credential {
	// Später laden wir hier die Passkeys aus der Datenbank.
	// Für Schritt 1 (Registrierung) reicht erstmal eine leere Liste.
	return []webauthn.Credential{}
}

type WebAuthnCredential struct {
	WebAuthnID          []byte
	WebAuthnName        string
	WebAuthnDisplayName string
	WebAuthnCredentials []webauthn.Credential
}

func WebAuthnInitialize() (*webauthn.WebAuthn, error) {

	config := &webauthn.Config{
		RPDisplayName: "IncogniCloud",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost:5173"},
	}

	return webauthn.New(config)

}

var sessionExamplePasskey *webauthn.SessionData

func saveSessionExamplePasskey(sessionData *webauthn.SessionData) {
	sessionExamplePasskey = sessionData
}

func loadSessionExamplePasskey() (*webauthn.SessionData, error) {
	if sessionExamplePasskey == nil {
		return nil, fmt.Errorf("Keine Session-Daten gefunden")
	}

	return sessionExamplePasskey, nil
}

/*func loadUserExamplePasskey(rawID []byte, userHandle []byte) (user webauthn.User, err error) {
	// Crude / Abstract example of retrieving the user for the rawID/userHandle value.

	//SQL Befehl später reinschreiben, um den User anhand der rawID/userHandle aus der Datenbank zu laden.
	return LoadUserByHandle(userHandle)
}*/

func HandlerPasskeyCreateChallenge(w *webauthn.WebAuthn) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Hier würdest du normalerweise den User aus der Datenbank laden.
		// Für dieses Beispiel erstellen wir einfach einen Dummy-User.
		testUser := User{
			ID:       1,
			Username: "AdminMax",
		}

		options, sessionData, err := w.BeginRegistration(testUser)
		if err != nil {
			http.Error(rw, "Fehler bei der Erstellung der Challenge: "+err.Error(), http.StatusInternalServerError)
			return
		}

		saveSessionExamplePasskey(sessionData)

		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(options)

	}
}

func HandlerPasskeyValidateCreateChallengeResponse(w *webauthn.WebAuthn) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		// Für später
		/*user, err := LoadUser()
		if err != nil {
			http.Error(rw, "Fehler beim Laden des Users: "+err.Error(), http.StatusInternalServerError)
			return
		}*/

		// Hier würdest du normalerweise den User aus der Datenbank laden.
		// Für dieses Beispiel erstellen wir einfach einen Dummy-User.
		testUser := User{
			ID:       1,
			Username: "AdminMax",
		}

		sessionData, err := loadSessionExamplePasskey()
		if err != nil {
			http.Error(rw, "Fehler beim Laden der Session-Daten: "+err.Error(), http.StatusInternalServerError)
			return
		}

		credential, err := w.FinishRegistration(testUser, *sessionData, r)
		if err != nil {
			http.Error(rw, "Fehler bei der Validierung der Challenge-Antwort: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Damit credential keinen unused error raushaut
		fmt.Print("Passkey erfolgreich erstellt: ")
		fmt.Printf("%+v\n", credential)

		// Für später
		// user.credentials = append(user.credentials, *credential)
		/*
			if err := SaveUser(user); err != nil {
				http.Error(rw, "Fehler beim Speichern des Users: "+err.Error(), http.StatusInternalServerError)
				return
			}
		*/

		rw.WriteHeader(http.StatusOK)

	}
}

func HandlerPasskeyLoginChallenge(w *webauthn.WebAuthn) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		assertion, s, err := w.BeginDiscoverableMediatedLogin(protocol.MediationDefault)
		if err != nil {
			http.Error(rw, "Fehler bei der Erstellung der Login-Challenge: "+err.Error(), http.StatusInternalServerError)
			return
		}

		saveSessionExamplePasskey(s)

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)

		encoder := json.NewEncoder(rw)
		if err := encoder.Encode(assertion); err != nil {
			http.Error(rw, "Fehler beim Senden der Login-Challenge: "+err.Error(), http.StatusInternalServerError)
			return
		}

	}

}

func HandlerPasskeyLoginChallengeResponse(w *webauthn.WebAuthn) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		/*sessionData, err := loadSessionExamplePasskey()
		if err != nil {
			http.Error(rw, "Fehler beim Laden der Session-Daten: "+err.Error(), http.StatusInternalServerError)
			return
		}

		validatedUser, validatedCredential, err := w.FinishPasskeyLogin(loadUserExamplePasskey, *sessionData, r)
		if err != nil {
			http.Error(rw, "Fehler bei der Validierung der Login-Challenge-Antwort: "+err.Error(), http.StatusBadRequest)
			return
		}

		user, ok := validatedUser.(*User)
		if !ok {
			http.Error(rw, "Fehler beim Konvertieren des validierten Users", http.StatusInternalServerError)
			return
		}

		var found bool

		for i, credential := range user.credentials {
			if bytes.Equal(credential.ID, validatedCredential.ID) {

				user.WebAuthnCredentials()[i] = *validatedCredential

				if err = SaveUser(user); err != nil {
					http.Error(rw, "Fehler beim Aktualisieren des Users: "+err.Error(), http.StatusInternalServerError)
					return
				}

				found = true

				break
			}
		}

		if !found {
			http.Error(rw, "Ungültige Anmeldeinformationen", http.StatusUnauthorized)
			return
		}*/

		rw.WriteHeader(http.StatusOK)

	}
}
