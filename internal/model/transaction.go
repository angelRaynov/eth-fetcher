package model

type Transaction struct {
	ID                int    `json:"id,omitempty"`
	TransactionHash   string `json:"transaction_hash"`
	TransactionStatus string `json:"transaction_status"`
	BlockHash         string `json:"block_hash"`
	BlockNumber       string `json:"block_number"`
	From              string `json:"sender"`
	To                string `json:"recipient"`
	ContractAddress   string `json:"contract_address"`
	LogsCount         int    `json:"logs_count"`
	Input             string `json:"input"`
	Value             string `json:"value"`
}

type Transactions struct {
	Transactions []Transaction `json:"transactions"`
}

type TransactionReceipt struct {
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		TransactionHash   string `json:"transactionHash"`
		TransactionStatus string `json:"status"`
		BlockHash         string `json:"blockHash"`
		BlockNumber       string `json:"blockNumber"`
		From              string `json:"from"`
		To                string `json:"to"`
		ContractAddress   string `json:"contractAddress"`
		LogsCount         int
		Logs              []interface{} `json:"logs"`
	} `json:"result"`
}

type TransactionByHash struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Input string `json:"input"`
		Value string `json:"value"`
	} `json:"result"`
}
