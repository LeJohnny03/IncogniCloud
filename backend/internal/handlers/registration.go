package handlers

import (
	"net/http"

	"backend/internal/database"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

type RegistrationHandler struct {
	webAuthn *webauthn.WebAuthn
	store    *database.Store
	sessions map[string]*webauthn.SessionData
}

func NewRegistrationHandler(wa *webauthn.WebAuthn, store *database.Store) *RegistrationHandler {
	return &RegistrationHandler{
		webAuthn: wa,
		store:    store,
		sessions: make(map[string]*webauthn.SessionData),
	}
}

func (h *RegistrationHandler) BeginRegistration(c *gin.Context) {
	var req struct {
		Username    string `json:"username" binding:"required"`
		DisplayName string `json:"display_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	existingUser, err := h.store.GetUserByUsername(c.Request.Context(), req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	user, err := h.store.CreateUser(c.Request.Context(), req.Username, req.DisplayName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	options, sessionData, err := h.webAuthn.BeginRegistration(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin registration"})
		return
	}

	sessionID := uuid.New().String()
	h.sessions[sessionID] = sessionData

	c.SetCookie("registration_session", sessionID, 300, "/", "", true, true)
	c.JSON(http.StatusOK, options)
}

func (h *RegistrationHandler) FinishRegistration(c *gin.Context) {
	sessionID, err := c.Cookie("registration_session")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No registration session"})
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

	credential, err := h.webAuthn.FinishRegistration(user, *sessionData, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to verify credential"})
		return
	}

	dbCredential := &models.Credential{
		ID:              credential.ID,
		UserID:          user.ID,
		PublicKey:       credential.PublicKey,
		AttestationType: credential.AttestationType,
		AAGUID:          credential.Authenticator.AAGUID,
		SignCount:       int(credential.Authenticator.SignCount),
		CloneWarning:    false,
	}

	if err := h.store.AddCredential(c.Request.Context(), dbCredential); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store credential"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user_id": user.ID,
	})
}
