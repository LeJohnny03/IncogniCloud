package handlers

import (
	"backend/internal/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(store *database.Store) gin.HandlerFunc {
	return func(c *gin.Context) {

		token, err := c.Cookie("auth_session")
		if err != nil || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - no cookie found"})
			return
		}

		session, err := store.GetSession(c.Request.Context(), token)
		if err != nil || session == nil {
			c.SetCookie("auth_session", "", -1, "/", "", true, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - invalid session"})
			return
		}

		if time.Now().UTC().After(session.ExpiresAt) {
			_ = store.DeleteSession(c.Request.Context(), token)
			c.SetCookie("auth_session", "", -1, "/", "", true, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - session expired"})
			return
		}

		c.Set("user_id", session.UserID)

		c.Next()

	}
}
