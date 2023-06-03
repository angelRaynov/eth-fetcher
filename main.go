package main

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"eth_fetcher/internal/model"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)



// TODO optimize table definitions
func main() {
	// Connection parameters
	host := "postgres"
	port := 5432
	user := "postgres"
	password := "postgres"
	dbname := "transaction_data"

	// Create the connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}
	defer db.Close()

	// Ping the database to check the connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to ping the database:", err)
		return
	}

	fmt.Println("Connected to the PostgreSQL database")

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

		var resp model.TransactionReceipt

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

		var resp2 model.TransactionByHash

		err = json.Unmarshal(body, &resp2)

		value, err := decodeTransactionValue(resp2.Result.Value)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		status, err := decodeTransactionValue(resp.Result.TransactionStatus)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		// Sample data for the insert
		transaction := model.Transaction{
			TransactionHash:   resp.Result.TransactionHash,
			TransactionStatus: status.String(),
			BlockHash:         resp.Result.BlockHash,
			BlockNumber:       resp.Result.BlockNumber,
			From:              resp.Result.From,
			To:                resp.Result.To,
			ContractAddress:   resp.Result.ContractAddress,
			LogsCount:         len(resp.Result.Logs),
			Input:             resp2.Result.Input,
			Value:             value.String(),
		}
		fmt.Printf("%#v\n", transaction)

		stmt, err := db.Prepare("INSERT INTO transactions (transaction_hash, transaction_status, block_hash, block_number, sender, recipient, contract_address, logs_count, input,value) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)")
		if err != nil {
			fmt.Println("Failed to prepare SQL statement:", err)
			return
		}
		defer stmt.Close()
		// Execute the insert statement with the provided data
		_, err = stmt.Exec(
			transaction.TransactionHash,
			transaction.TransactionStatus,
			transaction.BlockHash,
			transaction.BlockNumber,
			transaction.From,
			transaction.To,
			transaction.ContractAddress,
			transaction.LogsCount,
			transaction.Input,
			transaction.Value,
		)
		if err != nil {
			fmt.Println("Failed to execute insert statement:", err)
			return
		}

		fmt.Println("Insert successful")

	}

	rows, err := db.Query("SELECT * FROM transactions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Slice to hold the result structs
	var transactions []model.Transaction

	// Iterate over the rows and retrieve the column values
	for rows.Next() {
		var transaction model.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.TransactionHash,
			&transaction.TransactionStatus,
			&transaction.BlockHash,
			&transaction.BlockNumber,
			&transaction.From,
			&transaction.To,
			&transaction.ContractAddress,
			&transaction.LogsCount,
			&transaction.Input,
			&transaction.Value,
		)
		if err != nil {
			panic(err)
		}

		// Append the transaction to the slice
		transactions = append(transactions, transaction)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Print or process the retrieved transactions
	res := struct {
		Transactions []model.Transaction `json:"transactions"`
	}{
		Transactions: transactions,
	}

	r, _ := json.Marshal(res)
	fmt.Printf("\n\n\n%v\n", string(r))

}

func decodeTransactionValue(valueHex string) (*big.Int, error) {
	value := new(big.Int)
	value, success := value.SetString(valueHex[2:], 16) // Remove the "0x" prefix
	if !success {
		return nil, fmt.Errorf("failed to decode value")
	}

	return value, nil
}

