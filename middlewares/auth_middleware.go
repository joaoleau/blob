package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	oidc "github.com/coreos/go-oidc"
)

func AuthMiddleware(verifier *oidc.IDTokenVerifier, ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(401, gin.H{
				"error":            "Unauthorized",
				"error_description": "Authorization header is missing",
			})
			c.Abort()
			return
		}

		token := header[len(BearerSchema):]
		idToken, err := verifier.Verify(ctx, token)
		if err != nil {
			c.JSON(401, gin.H{
				"error":            "Unauthorized",
				"error_description": "Token verification failed: " + err.Error(),
			})
			c.Abort()
			return
		}

		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			c.JSON(401, gin.H{
				"error":            "Unauthorized",
				"error_description": "Failed to parse claims: " + err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
