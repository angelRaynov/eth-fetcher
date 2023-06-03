package transaction

import "eth_fetcher/internal/model"

type Fetcher interface {
	FetchBlockchainTransactionsByHashes(transactionHashes []string) []model.Transaction
	ListRequestedTransactions() ([]model.Transaction, error)
}
