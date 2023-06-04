package usecase

import (
	"eth_fetcher/infrastructure/config"
	"eth_fetcher/infrastructure/logger"
	"eth_fetcher/internal/authentication"
	"eth_fetcher/internal/model"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type authUseCase struct {
	l logger.ILogger
	ar authentication.PasswordGetter
	cfg *config.Application
}

func NewAuthUseCase(cfg *config.Application,l logger.ILogger, ar authentication.PasswordGetter) authentication.JWTGenerator {
	return &authUseCase{
		l: l,
		ar: ar,
		cfg: cfg,
	}
}

func (au *authUseCase)GenerateJWT(creds model.Credentials) (string, error) {
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

	return token.SignedString([]byte(au.cfg.JWTSecret))

}
