package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"encoding/json"

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

	fmt.Println("Backend-Server läuft auf http://localhost:8080 ...")
	log.Fatal(http.ListenAndServe(":8080", mux))

}