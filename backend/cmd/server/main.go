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

	frontendURL := os.Getenv("FRONTEND_URL")
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = "localhost"
	}

	allowedOrigins := []string{
		"http://localhost:5173",
		"http://localhost:3000",
	}

	if frontendURL != "" {
		allowedOrigins = append(allowedOrigins, frontendURL)
	}

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
		RPID:          domain,         // RPID muss zwingend mit der Domain übereinstimmen!
		RPOrigins:     allowedOrigins,
	}

	wa, err := config.NewWebAuthn(webAuthnConfig)
	if err != nil {
		log.Fatal("Failed to create WebAuthn:", err)
	}

	registrationHandler := handlers.NewRegistrationHandler(wa, store)
	authenticationHandler := handlers.NewAuthenticationHandler(wa, store)

	r := gin.Default()

	r.Use(func(c *gin.Context) {

		origin := c.Request.Header.Get("Origin")

		for _, allowed := range allowedOrigins {
			if origin == allowed {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

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

}
