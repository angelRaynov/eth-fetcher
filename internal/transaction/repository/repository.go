package repository

import (
	"database/sql"
	"eth_fetcher/internal/model"
	"fmt"
)

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *transactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (tr *transactionRepository) Store(transactions []model.Transaction) {
	err := tr.db.Ping()
	if err != nil {
		fmt.Println("Failed to ping the database:", err)
		return
	}
		//TODO TRANSACTION HASH MUST BE UNIQUE TO AVOID DUPLICATES!!!
	for _, transaction := range transactions {

		stmt, err := tr.db.Prepare("INSERT INTO transactions (transaction_hash, transaction_status, block_hash, block_number, sender, recipient, contract_address, logs_count, input,value) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)")
		if err != nil {
			fmt.Println("Failed to prepare SQL statement:", err)
			return
		}
		// Execute the insert statement with the provided data
		_, err = stmt.Exec(
			transaction.TransactionHash,
			transaction.TransactionStatus,
			transaction.BlockHash,
			transaction.BlockNumber,
			transaction.From,
			transaction.To,
			transaction.ContractAddress,
			transaction.LogsCount,
			transaction.Input,
			transaction.Value,
		)
		stmt.Close()

		if err != nil {
			fmt.Println("Failed to execute insert statement:", err)
		} else {
			fmt.Println("Insert successful")
		}

	}

}
