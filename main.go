package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"
)

type TransactionReceipt struct {
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		TransactionHash string      `json:"transactionHash"`
		TransactionStatus          string      `json:"status"`
		BlockHash       string      `json:"blockHash"`
		BlockNumber     string      `json:"blockNumber"`
		From            string      `json:"from"`
		To              string      `json:"to"`
		ContractAddress interface{} `json:"contractAddress"`
		LogsCount       int
		Logs            []interface{} `json:"logs"`

	} `json:"result"`
}

type TransactionByHash struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Input            string `json:"input"`
		Value            string `json:"value"`
	} `json:"result"`
}

func main() {
	rlpEncodedData := "0xf90110b842307839623266366133633265316165643263636366393262613636366332326430353361643064386135646137616131666435343737646364363537376234353234b842307835613537653330353163623932653264343832353135623037653762336431383531373232613734363534363537626436346131346333396361336639636632b842307837316239653262343464343034393863303861363239383866616337373664306561633062356239363133633337663966366639613462383838613862303537b842307863356639366266316235346433333134343235643233373962643737643765643465363434663763366538343961373438333230323862333238643464373938"
	rlpEncodedData = strings.TrimPrefix(rlpEncodedData, "0x")
	encodedData, err := hex.DecodeString(rlpEncodedData)
	if err != nil {
		log.Fatal(err)
	}

	var transactionHashes []string

	err = rlp.DecodeBytes(encodedData, &transactionHashes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Hashes: %+v\n", transactionHashes)

	url := "https://eth-goerli.g.alchemy.com/v2/jEvj-KdZ92ZUmX01Jpegiu52fpgEpE8_"

	for _, hash := range transactionHashes {
		s := fmt.Sprintf("{\"id\":1,\"jsonrpc\":\"2.0\",\"params\":[\"%s\"],\"method\":\"eth_getTransactionReceipt\"}", hash)
		payload := strings.NewReader(s)
		req, _ := http.NewRequest("POST", url, payload)

		req.Header.Add("accept", "application/json")
		req.Header.Add("content-type", "application/json")

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		var resp TransactionReceipt

		err = json.Unmarshal(body, &resp)
		resp.Result.LogsCount = len(resp.Result.Logs)

		s = fmt.Sprintf("{\"id\":1,\"jsonrpc\":\"2.0\",\"params\":[\"%s\"],\"method\":\"eth_getTransactionByHash\"}", hash)
		payload = strings.NewReader(s)
		req, _ = http.NewRequest("POST", url, payload)

		req.Header.Add("accept", "application/json")
		req.Header.Add("content-type", "application/json")

		res, _ = http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ = io.ReadAll(res.Body)

		var resp2 TransactionByHash

		err = json.Unmarshal(body, &resp2)

		fmt.Println("===========================================")
		fmt.Printf("%v\n",resp.Result)
		fmt.Printf("%v\n",resp2.Result)
		value, err := decodeTransactionValue(resp2.Result.Value)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Decoded Value: %s\n", value.String())
		}
		fmt.Println("===========================================")

	}

}

func decodeTransactionValue(valueHex string) (*big.Int, error) {
	value := new(big.Int)
	value, success := value.SetString(valueHex[2:], 16) // Remove the "0x" prefix
	if !success {
		return nil, fmt.Errorf("failed to decode value")
	}

	return value, nil
}