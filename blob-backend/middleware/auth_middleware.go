package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type SessionDetails struct {
		Email				 string 	 `db:"email"`
    Expires      time.Time `db:"expires"`
}

func AuthMiddleware(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		sessionToken := parts[1]
		var session SessionDetails
		
		query := `SELECT u.email, s.expires FROM "Session" s JOIN "User" u ON s.user_id = u.id WHERE session_token = $1`
		err := db.Get(&session, query, sessionToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			c.Abort()
			return
		}

		if time.Now().After(session.Expires) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session has expired"})
			c.Abort()
			return
		}

		c.Set("email", session.Email)

		c.Next()
	}
}
