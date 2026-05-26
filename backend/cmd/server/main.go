package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/internal/auth"
	"backend/internal/database"

	"github.com/joho/godotenv"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {

	// Lädt die .env Datei (nur für lokale Entwicklung nötig, auf dem NAS macht das später Docker/K8s)
	// Ignoriere den Fehler, falls die Datei nicht da ist (z.B. in der Produktion)
	_ = godotenv.Load("../.env")

	// Lese die Variablen dynamisch aus
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Baue den dynamischen Verbindungs-String
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := database.Connect(dsn)
	if err != nil {
		log.Fatalf("Fehler beim Verbinden: %v\n", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db, dbName); err != nil {
		log.Fatalf("Migrations-Fehler: %v\n", err)
	}

	fmt.Println("Backend-Server läuft und ist bereit...")

	webAuthInstance, err := auth.WebAuthnInitialize()
	if err != nil {
		log.Fatalf("Fehler bei WebAuthn-Initialisierung: %v\n", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		// CORS-Header setzen (WICHTIG für die lokale Entwicklung!)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Content-Type", "application/json")

		// Daten vorbereiten
		response := HealthResponse{
			Status:  "ok",
			Message: "Hallo von Go! Deine REST-API steht.",
		}

		// Als JSON an das Frontend schicken
		json.NewEncoder(w).Encode(response)
	})

	mux.HandleFunc("/webauthn/register/start", auth.HandlerPasskeyCreateChallenge(webAuthInstance))
	mux.HandleFunc("/webauthn/register/finish", auth.HandlerPasskeyValidateCreateChallengeResponse(webAuthInstance))
	mux.HandleFunc("/webauthn/login/start", auth.HandlerPasskeyLoginChallenge(webAuthInstance))
	mux.HandleFunc("/webauthn/login/finish", auth.HandlerPasskeyLoginChallengeResponse(webAuthInstance))

	protectedMux := auth.EnableCORS(auth.TailscaleOnly(mux))

	server := &http.Server{
		Addr:    ":8080",
		Handler: protectedMux, // Hier übergeben wir die Zwiebel an den Server
	}

	fmt.Println("Backend-Server läuft auf http://localhost:8080 ...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Fehler beim Starten des Servers: %v\n", err)
	}

}
