package repository

import (
	"database/sql"
	"eth_fetcher/infrastructure/logger"
	"eth_fetcher/internal/model"
	"eth_fetcher/internal/transaction"
	"eth_fetcher/internal/transaction/delivery/http"
	"fmt"
)

const (
	QueryFindByHash = `SELECT id, transaction_hash, transaction_status, block_hash, block_number, sender, recipient, contract_address, logs_count, input, value FROM transactions WHERE transaction_hash = $1 LIMIT 1`
	QueryFindAll    = `SELECT id, transaction_hash, transaction_status, block_hash, block_number, sender, recipient, contract_address, logs_count, input, value  FROM transactions ORDER BY id OFFSET $1 LIMIT $2`
	QueryInsert     = `INSERT INTO transactions (transaction_hash, transaction_status, block_hash, block_number, sender, recipient, contract_address, logs_count, input,value) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	QueryCount      = `SELECT COUNT(*) FROM transactions`
)

type transactionRepository struct {
	db *sql.DB
	l  logger.ILogger
}

func NewTransactionRepository(db *sql.DB, l logger.ILogger) transaction.StoreFinder {
	return &transactionRepository{
		db: db,
		l:  l,
	}
}

func (tr *transactionRepository) Store(transaction *model.Transaction) error {
	err := tr.db.Ping()
	if err != nil {
		return fmt.Errorf("pinging database:%w", err)
	}

	stmt, err := tr.db.Prepare(QueryInsert)
	if err != nil {
		return fmt.Errorf("preparing insert statement:%w", err)
	}

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

	defer stmt.Close()

	if err != nil {
		return fmt.Errorf("executing insert statement:%w", err)
	}

	return nil
}

func (tr *transactionRepository) FindAll() ([]*model.Transaction, error) {
	err := tr.db.Ping()
	if err != nil {
		tr.l.Errorw("pinging database", "error", err)
		return nil, fmt.Errorf("pinging database:%w", err)
	}

	var transactions []*model.Transaction

	batchSize := 10
	var totalCount int

	err = tr.db.QueryRow(QueryCount).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("counting records:%w", err)
	}

	if totalCount == 0 {
		return nil, http.ErrNoRecords
	}

	for offset := 0; offset < totalCount; offset += batchSize {
		// Calculate the remaining count to fetch
		remainingCount := totalCount - offset
		if remainingCount < batchSize {
			batchSize = remainingCount
		}

		// Fetch transactions in batches
		rows, err := tr.db.Query(QueryFindAll, offset, batchSize)
		if err != nil {
			return nil, fmt.Errorf("finding records:%w", err)
		}

		for rows.Next() {
			var tx model.Transaction
			err = rows.Scan(
				&tx.ID,
				&tx.TransactionHash,
				&tx.TransactionStatus,
				&tx.BlockHash,
				&tx.BlockNumber,
				&tx.From,
				&tx.To,
				&tx.ContractAddress,
				&tx.LogsCount,
				&tx.Input,
				&tx.Value,
			)
			if err != nil {
				return nil, err
			}
			transactions = append(transactions, &tx)
		}

		rows.Close()

		if len(transactions) < batchSize {
			// Break the loop if the number of fetched transactions is less than the batch size
			break
		}
	}

	return transactions, nil
}

func (tr *transactionRepository) FindByHash(hash string) (*model.Transaction, error) {
	row := tr.db.QueryRow(QueryFindByHash, hash)

	var tx model.Transaction
	err := row.Scan(
		&tx.ID,
		&tx.TransactionHash,
		&tx.TransactionStatus,
		&tx.BlockHash,
		&tx.BlockNumber,
		&tx.From,
		&tx.To,
		&tx.ContractAddress,
		&tx.LogsCount,
		&tx.Input,
		&tx.Value,
	)
	if err != nil {
		return &model.Transaction{}, err
	}

	return &tx, nil
}
