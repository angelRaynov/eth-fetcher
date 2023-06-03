package usecase

import (
	"eth_fetcher/helper"
	"eth_fetcher/infrastructure/api"
	"eth_fetcher/infrastructure/logger"
	"eth_fetcher/internal/model"
	"eth_fetcher/internal/transaction"
	"fmt"
)

type transactionUseCase struct {
	alchemy api.TransactionFetcher
	txRepo  transaction.StoreFinder
	l       logger.ILogger
}

func NewTransactionUseCase(alchemyAPI api.TransactionFetcher, txRepo transaction.StoreFinder, l logger.ILogger) *transactionUseCase {
	return &transactionUseCase{
		alchemy: alchemyAPI,
		txRepo:  txRepo,
		l:       l,
	}
}

func (tuc *transactionUseCase) FetchBlockchainTransactionsByHashes(transactionHashes []string) []model.Transaction {
	var transactions []model.Transaction

	for _, hash := range transactionHashes {
		tuc.l.Debugw("fetching single transaction", "transaction_hash", hash)
		tx, err := tuc.fetchSingleTransaction(hash)
		if err != nil {
			tuc.l.Warnw("fetching single transaction", "error", err, "transaction_hash", hash)
			continue
		}

		transactions = append(transactions, tx)

	}

	tuc.l.Infow("all transactions fetched successfully")

	return transactions
}

func (tuc *transactionUseCase) ListRequestedTransactions() ([]model.Transaction, error) {
	res, err := tuc.txRepo.FindAll()
	if err != nil {
		tuc.l.Infow("finding all records", "error", err)
		return nil, fmt.Errorf("finding all records:%w", err)
	}

	tuc.l.Infow("all requested transactions listed successfully")

	return res, nil

}

func (tuc *transactionUseCase) fetchSingleTransaction(hash string) (model.Transaction, error) {
	tx, err := tuc.txRepo.FindByHash(hash)
	if err == nil {
		tuc.l.Debugw("transaction found in db, skip request", "transaction_hash", hash)
		return tx, nil
	}

	txReceipt, err := tuc.alchemy.GetTransactionReceiptByHash(hash)
	if err != nil {
		tuc.l.Warnw("getting transaction receipt by hash", "transaction_hash", hash, "error", err)
		return model.Transaction{}, err
	}

	status, err := helper.DecodeHexBigInt(txReceipt.Result.TransactionStatus)
	if err != nil {
		tuc.l.Warnw("decoding hex status", "transaction_hash", hash, "error", err, "hex_status", txReceipt.Result.TransactionStatus)
		return model.Transaction{}, err
	}

	txByHash, err := tuc.alchemy.GetTransactionByHash(hash)
	if err != nil {
		tuc.l.Warnw("getting transaction by hash", "transaction_hash", hash, "error", err)
		return model.Transaction{}, err
	}

	value, err := helper.DecodeHexBigInt(txByHash.Result.Value)
	if err != nil {
		tuc.l.Warnw("decoding hex value", "transaction_hash", hash, "error", err, "hex_value", txByHash.Result.Value)
		return model.Transaction{}, err
	}

	blockNumber, err := helper.DecodeHexBigInt(txReceipt.Result.BlockNumber)
	if err != nil {
		tuc.l.Warnw("decoding hex block number", "transaction_hash", hash, "error", err, "hex_value", txReceipt.Result.BlockNumber)
		return model.Transaction{}, err
	}

	// Sample data for the insert
	tx = model.Transaction{
		TransactionHash:   txReceipt.Result.TransactionHash,
		TransactionStatus: status.String(),
		BlockHash:         txReceipt.Result.BlockHash,
		BlockNumber:       blockNumber.String(),
		From:              txReceipt.Result.From,
		To:                txReceipt.Result.To,
		ContractAddress:   txReceipt.Result.ContractAddress,
		LogsCount:         len(txReceipt.Result.Logs),
		Input:             txByHash.Result.Input,
		Value:             value.String(),
	}

	err = tuc.txRepo.Store(tx)
	if err != nil {
		tuc.l.Warnw("storing transaction", "error", err, "transaction_hash", tx.TransactionHash)
	}

	return tx, nil
}
