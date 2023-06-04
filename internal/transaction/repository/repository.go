package repository

import (
	"database/sql"
	"eth_fetcher/internal/model"
	"eth_fetcher/internal/transaction"
	"eth_fetcher/internal/transaction/delivery/http"
	"fmt"
)

const (
	QueryFindByHash            = `SELECT id, transaction_hash, transaction_status, block_hash, block_number, sender, recipient, contract_address, logs_count, input, value FROM transactions WHERE transaction_hash = $1 LIMIT 1`
	QueryFindAll               = `SELECT id, transaction_hash, transaction_status, block_hash, block_number, sender, recipient, contract_address, logs_count, input, value  FROM transactions ORDER BY id OFFSET $1 LIMIT $2`
	QueryInsert                = `INSERT INTO transactions (transaction_hash, transaction_status, block_hash, block_number, sender, recipient, contract_address, logs_count, input,value) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	QueryInsertHash            = `INSERT INTO transaction_history (username, transaction_hash) VALUES ($1, $2)`
	QueryCount                 = `SELECT COUNT(*) FROM transactions`
	QueryRequestedTransactions = `SELECT transaction_hash FROM transaction_history WHERE username = $1`
	QueryTransactionExist = `SELECT COUNT(*) FROM transaction_history WHERE username = $1 AND transaction_hash = $2`
)

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) transaction.StoreFinder {
	return &transactionRepository{
		db: db,
	}
}

func (tr *transactionRepository) Store(transaction *model.Transaction) error {
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

func (tr *transactionRepository) StoreHashesPerUser(user, hash string) error {
	exist := tr.checkTransactionExists(user,hash)
	if exist {
		return nil
	}

	stmt, err := tr.db.Prepare(QueryInsertHash)
	if err != nil {
		return fmt.Errorf("preparing insert statement:%w", err)
	}

	_, err = stmt.Exec(
		user,
		hash,
	)

	defer stmt.Close()

	if err != nil {
		return fmt.Errorf("executing insert statement:%w", err)
	}

	return nil
}
func (tr *transactionRepository) FindAll() ([]*model.Transaction, error) {
	var transactions []*model.Transaction

	batchSize := 10
	var totalCount int

	err := tr.db.QueryRow(QueryCount).Scan(&totalCount)
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

func (tr *transactionRepository) FindTransactionHashesPerUser(user string) ([]string, error) {
	var transactionHashes []string

	rows, err := tr.db.Query(QueryRequestedTransactions,user)
	if err != nil {
		return nil, fmt.Errorf("finding records:%w", err)
	}

	var hash string
	for rows.Next() {
		err = rows.Scan(&hash)
		if err != nil {
			return nil, err
		}
		transactionHashes = append(transactionHashes, hash)
	}

	defer rows.Close()

	return transactionHashes, nil
}

func (tr *transactionRepository) checkTransactionExists(user, hash string) bool {
	var count int

	err := tr.db.QueryRow(QueryTransactionExist, user, hash).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}
