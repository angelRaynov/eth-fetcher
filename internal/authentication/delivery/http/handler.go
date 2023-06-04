package http

import (
	"eth_fetcher/infrastructure/logger"
	"eth_fetcher/internal/authentication"
	"eth_fetcher/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type authHandler struct {
	l           logger.ILogger
	authUseCase authentication.JWTGenerator
}

func NewAuthHandler(l logger.ILogger, authUsecase authentication.JWTGenerator) authentication.Authenticator {
	return &authHandler{
		l: l,
		authUseCase: authUsecase,
	}
}

// Authenticate handles the POST request to /lime/authenticate
func (ah *authHandler) Authenticate(c *gin.Context) {
	var creds model.Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenString, err := ah.authUseCase.GenerateJWT(creds)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			fmt.Println("err ", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})

			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	// Return the token in the response
	c.JSON(http.StatusOK, model.AuthResponse{Token: tokenString})
}


