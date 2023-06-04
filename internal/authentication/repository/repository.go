package repository

import (
	"database/sql"
	"eth_fetcher/internal/authentication"
	"fmt"
)

const QueryGetPassword = `SELECT password FROM users WHERE username = $1`

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) authentication.PasswordGetter {
	return &authRepository{
		db: db,
	}
}

func (ar *authRepository) GetUserPassword(username string) (string, error) {

	row := ar.db.QueryRow(QueryGetPassword, username)

	var password string
	err := row.Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("failed to retrieve password: %v", err)
	}

	return password, nil

}
