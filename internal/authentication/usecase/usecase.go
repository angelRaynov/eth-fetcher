package usecase

import (
	"eth_fetcher/infrastructure/logger"
	"eth_fetcher/internal/authentication/repository"
	"eth_fetcher/internal/model"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

//todo use interface and make private
type AuthUseCase struct {
	l logger.ILogger
	ar *repository.AuthRepository
}

func NewAuthUseCase(l logger.ILogger, ar *repository.AuthRepository) *AuthUseCase {
	return &AuthUseCase{
		l: l,
		ar: ar,
	}
}

func (au *AuthUseCase)GenerateJWT(creds model.Credentials) (string, error) {
	hashedPW, err := au.ar.GetUserPassword(creds.Username)
	if err != nil {
		log.Fatal("get pass:",err)
	}

	// Compare the entered password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPW), []byte(creds.Password))
	if err != nil {
		return "", err
	}


	// Create a JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = creds.Username

	// TODO: Add additional claims or custom data to the token if needed

	// Sign the token with a secret key
	// Replace "secret" with your own secret key
	return token.SignedString([]byte("secret"))

}
