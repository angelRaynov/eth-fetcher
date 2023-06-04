package repository

import (
	"database/sql"
	"eth_fetcher/infrastructure/logger"
	"fmt"
)

const QueryGetPassword = `SELECT password FROM users WHERE username = $1`

type AuthRepository struct {
	db *sql.DB
	l  logger.ILogger
}

func NewAuthRepository(db *sql.DB, l logger.ILogger) *AuthRepository {
	return &AuthRepository{
		db: db,
		l:  l,
	}
}

func (ar *AuthRepository) GetUserPassword(username string) (string, error) {

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
