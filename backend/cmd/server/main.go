package main

import (
	"fmt"
	"log"
	"os"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Fatal("Fehler beim Verbinden:", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db, dbName); err != nil {
		log.Fatal("Migrations-Fehler:", err)
	}

	fmt.Println("Backend-Server läuft und ist bereit...")

	store := database.NewStore(db)

	webAuthnConfig := config.WebAuthnConfig{
		RPDisplayName: "IncogniCloud", // Der Name deiner App
		RPID:          "localhost",    // Normalerweise die Domain deiner App
		RPOrigins:     []string{"http://localhost:5173"},
	}

	wa, err := config.NewWebAuthn(webAuthnConfig)
	if err != nil {
		log.Fatal("Failed to create WebAuthn:", err)
	}

	registrationHandler := handlers.NewRegistrationHandler(wa, store)
	authenticationHandler := handlers.NewAuthenticationHandler(wa, store)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	api := r.Group("/api")
	{
		api.POST("/register/begin", registrationHandler.BeginRegistration)
		api.POST("/register/finish", registrationHandler.FinishRegistration)
		api.POST("/login/begin", authenticationHandler.BeginAuthentication)
		api.POST("/login/finish", authenticationHandler.FinishAuthentication)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}

	/*mux := http.NewServeMux()

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

	protectedMux := handlers.EnableCORS(handlers.TailscaleOnly(mux))

	server := &http.Server{
		Addr:    ":8080",
		Handler: protectedMux, // Hier übergeben wir die Zwiebel an den Server
	}

	fmt.Println("Backend-Server läuft auf http://localhost:8080 ...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Fehler beim Starten des Servers:", err)
	}*/

}
