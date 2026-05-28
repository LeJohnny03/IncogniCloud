package handlers

import (
	"fmt"
	"net/http"
	"time"

	"backend/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

type AuthenticationHandler struct {
	webAuthn *webauthn.WebAuthn
	store    *database.Store
	sessions map[string]*webauthn.SessionData
}

func NewAuthenticationHandler(wa *webauthn.WebAuthn, store *database.Store) *AuthenticationHandler {

	return &AuthenticationHandler{
		webAuthn: wa,
		store:    store,
		sessions: make(map[string]*webauthn.SessionData),
	}

}

func (h *AuthenticationHandler) BeginAuthentication(c *gin.Context) {

	var req struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.store.GetUserByUsername(c.Request.Context(), req.Username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	credentials, err := h.store.GetUserCredentials(c.Request.Context(), user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load credentials"})
		return
	}

	user.Credentials = credentials

	options, sessionData, err := h.webAuthn.BeginLogin(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin authentication"})
		return
	}

	sessionID := uuid.New().String()
	h.sessions[sessionID] = sessionData

	c.SetCookie("webauthn_challenge", sessionID, 300, "/", "", true, true)

	c.JSON(http.StatusOK, options)

}

func (h *AuthenticationHandler) FinishAuthentication(c *gin.Context) {

	sessionID, err := c.Cookie("webauthn_challenge")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No authentication session"})
		return
	}

	sessionData, exists := h.sessions[sessionID]

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		return
	}

	defer delete(h.sessions, sessionID)

	userID, err := uuid.Parse(string(sessionData.UserID))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.store.GetUserByID(c.Request.Context(), userID)

	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	credentials, err := h.store.GetUserCredentials(c.Request.Context(), user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load credentials"})
		return
	}

	user.Credentials = credentials

	credential, err := h.webAuthn.FinishLogin(user, *sessionData, c.Request)

	if err != nil {
		fmt.Printf("WebAuthn FinishLogin Error: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed", "details": err.Error()})
		return

	}

	if credential.Authenticator.CloneWarning {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  "Credential may be cloned",
			"action": "contact_support",
		})
		return

	}

	err = h.store.UpdateCredential(
		c.Request.Context(),
		credential.ID,
		int(credential.Authenticator.SignCount),
		credential.Authenticator.CloneWarning,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update credential"})
		return
	}

	c.SetCookie("webauthn_challenge", "", -1, "/", "", true, true)

	// Create a secure database session

	sessionToken, err := database.GenerateSecureToken(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session token"})
		return
	}

	expiresAt := time.Now().UTC().Add(2 * time.Hour)
	maxAgeSeconds := int(2 * 60 * 60) // 2 hours in seconds

	err = h.store.CreateSession(c.Request.Context(), user.ID, sessionToken, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session to database"})
		return
	}

	c.SetCookie("auth_session", sessionToken, maxAgeSeconds, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"user_id":  user.ID,
		"username": user.Username,
	})

}

func (h *AuthenticationHandler) Logout(c *gin.Context) {
	token, err := c.Cookie("auth_session")

	if err == nil && token != "" {
		_ = h.store.DeleteSession(c.Request.Context(), token)
	}
	c.SetCookie("auth_session", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"success": true})
}
