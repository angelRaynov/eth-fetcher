package model

type Transaction struct {
	ID                int    `json:"id,omitempty"`
	TransactionHash   string `json:"transaction_hash"`
	TransactionStatus int8   `json:"transaction_status"`
	BlockHash         string `json:"block_hash"`
	BlockNumber       int64  `json:"block_number"`
	From              string `json:"sender"`
	To                string `json:"recipient"`
	ContractAddress   string `json:"contract_address"`
	LogsCount         int    `json:"logs_count"`
	Input             string `json:"input"`
	Value             string `json:"value"`
}

type Transactions struct {
	Transactions []*Transaction `json:"transactions"`
}
