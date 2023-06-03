package repository

type transactionRepository struct {
	
}

func NewTransactionRepository() *transactionRepository {
	return &transactionRepository{}
}

func (tr *transactionRepository) Store() {

		//TODO TRANSACTION HASH MUST BE UNIQUE TO AVOID DUPLICATES!!!

		//stmt, err := db.Prepare("INSERT INTO transactions (transaction_hash, transaction_status, block_hash, block_number, sender, recipient, contract_address, logs_count, input,value) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)")
		//if err != nil {
		//	fmt.Println("Failed to prepare SQL statement:", err)
		//	return
		//}
		//defer stmt.Close()
		//// Execute the insert statement with the provided data
		//_, err = stmt.Exec(
		//	transaction.TransactionHash,
		//	transaction.TransactionStatus,
		//	transaction.BlockHash,
		//	transaction.BlockNumber,
		//	transaction.From,
		//	transaction.To,
		//	transaction.ContractAddress,
		//	transaction.LogsCount,
		//	transaction.Input,
		//	transaction.Value,
		//)
		//if err != nil {
		//	fmt.Println("Failed to execute insert statement:", err)
		//	return
		//}
		//
		//fmt.Println("Insert successful")

}
