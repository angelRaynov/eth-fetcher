package server

import (
	"eth_fetcher/infrastructure/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
)

func AuthMiddleware(cfg *config.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.GetHeader("AUTH_TOKEN")
		if authToken != "" {
			// Verify and parse the JWT token
			token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				c.Abort()
				return
			}

			// Extract the username claim from the token's payload
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				username := claims["username"].(string)
				// Store the username in the Gin context
				c.Set("username", username)
			}
		}

		// Token is valid, proceed with the next handler
		c.Next()
	}
}
