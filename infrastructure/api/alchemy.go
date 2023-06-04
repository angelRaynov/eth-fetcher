package api

import (
	"encoding/json"
	"eth_fetcher/infrastructure/config"
	"eth_fetcher/infrastructure/logger"
	"eth_fetcher/internal/model"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type TransactionFetcher interface {
	GetTransactionReceiptByHash(txHash string) (*model.TransactionReceipt, error)
	GetTransactionByHash(txHash string) (*model.TransactionByHash, error)
}

type alchemyAPI struct {
	l      logger.ILogger
	client *http.Client
	cfg    *config.Application
}

func NewAlchemyAPI(cfg *config.Application, l logger.ILogger) TransactionFetcher {
	return &alchemyAPI{
		l:   l,
		cfg: cfg,
	}
}

func (a *alchemyAPI) GetTransactionReceiptByHash(txHash string) (*model.TransactionReceipt, error) {
	url := fmt.Sprintf("%s/%s", a.cfg.EthNodeURL, a.cfg.APIKey)
	p := fmt.Sprintf("{\"id\":1,\"jsonrpc\":\"2.0\",\"params\":[\"%s\"],\"method\":\"eth_getTransactionReceipt\"}", txHash)

	payload := strings.NewReader(p)
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		a.l.Warnw("assembling transaction receipt request", "url", url, "payload", p, "transaction_hash", txHash, "error", err)
		return nil, fmt.Errorf("assembling transaction receipt request:%w", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		a.l.Warnw("requesting transaction receipt", "url", url, "payload", p, "transaction_hash", txHash, "error", err)
		return nil, fmt.Errorf("requesting transaction receipt:%w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		a.l.Warnw("reading transaction receipt response", "url", url, "payload", p, "transaction_hash", txHash, "error", err)
		return nil, fmt.Errorf("reading transaction receipt response:%w", err)
	}

	defer response.Body.Close()

	var t model.TransactionReceipt

	err = json.Unmarshal(body, &t)
	if err != nil {
		a.l.Warnw("unmarshalling transaction receipt result", "url", url, "payload", p, "transaction_hash", txHash, "error", err)
		return nil, fmt.Errorf("unmarshaling transaction receipt result:%w", err)
	}

	if t.Error.Message != "" {
		a.l.Warnw("fetching transaction receipt error", "url", url, "payload", p, "transaction_hash", txHash, "error", t.Error.Message)
		return nil, fmt.Errorf("fetching transaction receipt error:%s", t.Error.Message)
	}

	t.Result.LogsCount = len(t.Result.Logs)

	return &t, nil
}

func (a *alchemyAPI) GetTransactionByHash(txHash string) (*model.TransactionByHash, error) {
	url := fmt.Sprintf("%s/%s", a.cfg.EthNodeURL, a.cfg.APIKey)

	p := fmt.Sprintf("{\"id\":1,\"jsonrpc\":\"2.0\",\"params\":[\"%s\"],\"method\":\"eth_getTransactionByHash\"}", txHash)
	payload := strings.NewReader(p)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		a.l.Warnw("assembling transaction request", "url", url, "payload", p, "transaction_hash", txHash, "error", err)
		return nil, fmt.Errorf("assembling transaction request:%w", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		a.l.Warnw("requesting transaction", "url", url, "payload", p, "transaction_hash", txHash, "error", err)
		return nil, fmt.Errorf("requesting transaction:%w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		a.l.Warnw("reading transaction response", "url", url, "payload", p, "transaction_hash", txHash, "error", err)
		return nil, fmt.Errorf("reading transaction response:%w", err)
	}

	defer response.Body.Close()

	var t model.TransactionByHash

	err = json.Unmarshal(body, &t)
	if err != nil {
		a.l.Warnw("unmarshalling transaction result", "url", url, "payload", p, "transaction_hash", txHash, "error", err)
		return nil, fmt.Errorf("unmarshalling transaction result:%w", err)
	}

	if t.Error.Message != "" {
		a.l.Warnw("fetching transaction error", "url", url, "payload", p, "transaction_hash", txHash, "error", t.Error.Message)
		return nil, fmt.Errorf("fetching transaction error:%s", t.Error.Message)
	}

	return &t, nil
}
