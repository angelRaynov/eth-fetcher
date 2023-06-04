package model

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
	Error AlchemyError `json:"error"`
}

type TransactionByHash struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Input string `json:"input"`
		Value string `json:"value"`
	} `json:"result"`
	Error AlchemyError `json:"error"`
}

type AlchemyError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
