package transaction

import "eth_fetcher/internal/model"

type Fetcher interface {
	FetchBlockchainTransactionsByHashes(transactionHashes []string) ([]*model.Transaction, error)
	ListRequestedTransactions() ([]*model.Transaction, error)
	CreateTransactionHistory(user string, transactionHashes []string)
	GetTransactionHistory(user string) ([]string, error)
}
