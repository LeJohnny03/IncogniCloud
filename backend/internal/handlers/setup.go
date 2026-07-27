package handlers

import (
	"backend/internal/config"
	"backend/internal/database"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SetupHandler struct {
	store *database.Store
}

func NewSetupHandler(store *database.Store) *SetupHandler {
	return &SetupHandler{store: store}
}

// Struct für den Request-Body aus dem Frontend
type SetupFinishRequest struct {
	config.ServerPreferences
}

func (h *SetupHandler) Status(c *gin.Context) {
	isSetup, err := h.store.IsSystemSetup(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"needs_setup": !isSetup,
	})

}

func (h *SetupHandler) Finish(c *gin.Context) {
	var req SetupFinishRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültige Eingabedaten"})
		return
	}

	prefsBytes, err := json.Marshal(req.ServerPreferences)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while processing Data"})
		return
	}

	err = h.store.SaveUserPreferences(c.Request.Context(), prefsBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error with saving server preferences."})
		return
	}

	fmt.Print(req.StorageDrives, req.StorageFolder, req.EnableBackups, req.BackupDrive, req.BackupType, req.CloudFolderName)

	c.JSON(http.StatusOK, gin.H{"message": "Setup successfully saved into server preferences."})

}
