package usecase

import (
	"errors"
	"eth_fetcher/helper"
	"eth_fetcher/infrastructure/api"
	"eth_fetcher/infrastructure/logger"
	"eth_fetcher/internal/model"
	"eth_fetcher/internal/transaction"
	"fmt"
	"strconv"
)

type transactionUseCase struct {
	alchemy api.TransactionFetcher
	txRepo  transaction.StoreFinder
	l       logger.ILogger
}

func NewTransactionUseCase(alchemyAPI api.TransactionFetcher, txRepo transaction.StoreFinder, l logger.ILogger) transaction.Fetcher {
	return &transactionUseCase{
		alchemy: alchemyAPI,
		txRepo:  txRepo,
		l:       l,
	}
}

func (tuc *transactionUseCase) FetchBlockchainTransactionsByHashes(transactionHashes []string) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	var errs error
	for _, hash := range transactionHashes {
		tuc.l.Debugw("fetching single transaction", "transaction_hash", hash)
		tx, err := tuc.fetchSingleTransaction(hash)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("fetch with tx hash %s:%w", hash, err))
			tuc.l.Warnw("fetching single transaction", "error", err, "transaction_hash", hash)
			continue
		}

		transactions = append(transactions, tx)

	}

	return transactions, errs
}

func (tuc *transactionUseCase) ListRequestedTransactions() ([]*model.Transaction, error) {
	res, err := tuc.txRepo.FindAll()
	if err != nil {
		tuc.l.Infow("finding all records", "error", err)
		return nil, fmt.Errorf("finding all records:%w", err)
	}

	tuc.l.Infow("all requested transactions listed successfully")

	return res, nil

}

func (tuc *transactionUseCase) fetchSingleTransaction(hash string) (*model.Transaction, error) {
	tx, err := tuc.txRepo.FindByHash(hash)
	if err == nil {
		tuc.l.Debugw("transaction found in db, skip request", "transaction_hash", hash)
		return tx, nil
	}

	txReceipt, err := tuc.alchemy.GetTransactionReceiptByHash(hash)
	if err != nil {
		tuc.l.Warnw("getting transaction receipt by hash", "transaction_hash", hash, "error", err)
		return &model.Transaction{}, err
	}

	txByHash, err := tuc.alchemy.GetTransactionByHash(hash)
	if err != nil {
		tuc.l.Warnw("getting transaction by hash", "transaction_hash", hash, "error", err)
		return &model.Transaction{}, err
	}

	tx, err = prepareTxData(txReceipt, txByHash)
	if err != nil {
		tuc.l.Warnw("preparing data for insert", "transaction_hash", hash, "error", err)
		return &model.Transaction{}, err
	}

	err = tuc.txRepo.Store(tx)
	if err != nil {
		tuc.l.Warnw("storing transaction", "error", err, "transaction_hash", tx.TransactionHash)
	}

	return tx, nil
}

func prepareTxData(txReceipt *model.TransactionReceipt, txByHash *model.TransactionByHash) (*model.Transaction, error) {
	var errs error

	status, err := helper.DecodeHexBigInt(txReceipt.Result.TransactionStatus)
	if err != nil {
		errs = errors.Join(errs, fmt.Errorf("decoding hex status:%w", err))
	}

	intStatus, err := strconv.ParseInt(status.String(), 10, 8)
	if err != nil {
		errs = errors.Join(errs, fmt.Errorf("converting string status:%w", err))
	}

	value, err := helper.DecodeHexBigInt(txByHash.Result.Value)
	if err != nil {
		errs = errors.Join(errs, fmt.Errorf("decoding hex value:%w", err))
	}

	blockNumber, err := helper.DecodeHexBigInt(txReceipt.Result.BlockNumber)
	if err != nil {
		errs = errors.Join(errs, fmt.Errorf("decoding hex block number:%w", err))
	}

	intBlockNumber, err := strconv.ParseInt(blockNumber.String(), 10, 64)
	if err != nil {
		errs = errors.Join(errs, fmt.Errorf("converting string block number:%w", err))
	}

	return &model.Transaction{
		TransactionHash:   txReceipt.Result.TransactionHash,
		TransactionStatus: int8(intStatus),
		BlockHash:         txReceipt.Result.BlockHash,
		BlockNumber:       intBlockNumber,
		From:              txReceipt.Result.From,
		To:                txReceipt.Result.To,
		ContractAddress:   txReceipt.Result.ContractAddress,
		LogsCount:         len(txReceipt.Result.Logs),
		Input:             txByHash.Result.Input,
		Value:             value.String(),
	}, err
}

func (tuc *transactionUseCase) CreateTransactionHistory(user string, transactionHashes []string) {
	for _, hash := range transactionHashes {
		err := tuc.txRepo.StoreHashesPerUser(user, hash)
		if err != nil {
			tuc.l.Warnw("storing hash per user", "error", err, "transaction_hash", hash, "username", user)
		}
		tuc.l.Infow("storing hash per user",  "transaction_hash", hash, "username", user)
	}
}

func (tuc *transactionUseCase) GetTransactionHistory(user string) ([]string, error) {
		transactionHashes, err := tuc.txRepo.FindTransactionHashesPerUser(user)
		if err != nil {
			tuc.l.Warnw("finding transactions per user", "error", err,"username", user)
			return nil, fmt.Errorf("finding transactions per user:%w",err)
		}

		return transactionHashes,nil
}
