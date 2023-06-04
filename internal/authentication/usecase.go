package authentication

import "eth_fetcher/internal/model"

type JWTGenerator interface {
	GenerateJWT(creds model.Credentials) (string, error)
}
