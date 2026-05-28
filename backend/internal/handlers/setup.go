package handlers

import (
	"backend/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SetupHandler struct {
	store *database.Store
}

func NewSetupHandler(store *database.Store) *SetupHandler {
	return &SetupHandler{store: store}
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
