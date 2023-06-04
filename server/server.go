package server

import (
	"eth_fetcher/infrastructure/api"
	"eth_fetcher/infrastructure/config"
	"eth_fetcher/infrastructure/database"
	"eth_fetcher/infrastructure/logger"
	httpHandler "eth_fetcher/internal/transaction/delivery/http"
	"eth_fetcher/internal/transaction/repository"
	"eth_fetcher/internal/transaction/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
)

func Run() {
	cfg := config.New()
	l := logger.Init(cfg.AppMode)
	db := database.Init(cfg, l)

	alchemy := api.NewAlchemyAPI(cfg, l)
	tr := repository.NewTransactionRepository(db, l)
	tuc := usecase.NewTransactionUseCase(alchemy, tr, l)
	h := httpHandler.NewTransactionHandler(tuc, l)
	router := gin.Default()

	router.GET("/lime/all", h.ExploreAllTransactions)
	router.GET("/lime/eth/:rlphex", h.ExploreTransactionsByRLP)

	router.POST("/lime/authenticate",Authenticate)
	l.Infow("listening on port", "port", cfg.APIPort)

	port := fmt.Sprintf(":%s", cfg.APIPort)
	log.Fatal(router.Run(port))
}

// Authenticate handles the POST request to /lime/authenticate
func Authenticate(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the username and password are valid
	password, exists := users[creds.Username]
	if !exists || password != creds.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create a JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = creds.Username

	// TODO: Add additional claims or custom data to the token if needed

	// Sign the token with a secret key
	// Replace "secret" with your own secret key
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return the token in the response
	c.JSON(http.StatusOK, AuthResponse{Token: tokenString})
}
// User represents a user with username and password
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Credentials represents the request body for authentication
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents the response format for authentication
type AuthResponse struct {
	Token string `json:"token"`
}

var users = map[string]string{
	"alice":  "alice",
	"bob":    "bob",
	"carol":  "carol",
	"dave":   "dave",
}
