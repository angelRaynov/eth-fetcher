package api

import (
	"encoding/json"
	"eth_fetcher/internal/model"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type TransactionFetcher interface {
	GetTransactionReceiptByHash(txHash string) *model.TransactionReceipt
	GetTransactionByHash(txHash string) *model.TransactionByHash
}
type alchemyAPI struct {

}

func NewAlchemyAPI() *alchemyAPI {
	return &alchemyAPI{}
}

func (a *alchemyAPI) GetTransactionReceiptByHash(txHash string) *model.TransactionReceipt {
	url := "https://eth-goerli.g.alchemy.com/v2/jEvj-KdZ92ZUmX01Jpegiu52fpgEpE8_"

	s := fmt.Sprintf("{\"id\":1,\"jsonrpc\":\"2.0\",\"params\":[\"%s\"],\"method\":\"eth_getTransactionReceipt\"}", txHash)
	payload := strings.NewReader(s)
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	var resp model.TransactionReceipt

	err := json.Unmarshal(body, &resp)
	if err != nil {
		log.Fatal(err)
	}
	resp.Result.LogsCount = len(resp.Result.Logs)

	return &resp
}

func (a *alchemyAPI) GetTransactionByHash(txHash string) *model.TransactionByHash {
	url := "https://eth-goerli.g.alchemy.com/v2/jEvj-KdZ92ZUmX01Jpegiu52fpgEpE8_"

	s := fmt.Sprintf("{\"id\":1,\"jsonrpc\":\"2.0\",\"params\":[\"%s\"],\"method\":\"eth_getTransactionByHash\"}", txHash)
	payload := strings.NewReader(s)
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var resp2 model.TransactionByHash

	err = json.Unmarshal(body, &resp2)
	if err != nil {
		log.Fatal(err)
	}
	return &resp2
}