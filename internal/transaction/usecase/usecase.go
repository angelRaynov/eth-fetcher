package usecase

import (
	"eth_fetcher/infrastructure/api"
	"eth_fetcher/internal/model"
	http2 "eth_fetcher/internal/transaction/delivery/http"
	"fmt"
)

type transactionUseCase struct {
	alchemy api.TransactionFetcher
}

func NewTransactionUseCase(alchemyAPI api.TransactionFetcher) *transactionUseCase {
	return &transactionUseCase{
		alchemy: alchemyAPI,
	}
}

func (tuc *transactionUseCase) FetchBlockchainTransactionsByHashes(transactionHashes []string) []model.Transaction {

	var transactions []model.Transaction
	for _, hash := range transactionHashes {


		txReceipt := tuc.alchemy.GetTransactionReceiptByHash(hash)
		
		status, err := http2.DecodeTransactionValue(txReceipt.Result.TransactionStatus)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		txByHash := tuc.alchemy.GetTransactionByHash(hash)

		value, err := http2.DecodeTransactionValue(txByHash.Result.Value)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		blockNumber, err := http2.DecodeTransactionValue(txReceipt.Result.BlockNumber)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		// Sample data for the insert
		transaction := model.Transaction{
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
		transactions = append(transactions, transaction)

		//TODO TRANSACTION HASH MUST BE UNIQUE TO AVOID DUPLICATES!!!


	}
	return transactions
}
