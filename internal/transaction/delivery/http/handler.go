package http

import (
	"encoding/hex"
	"encoding/json"
	"eth_fetcher/infrastructure/database"
	"eth_fetcher/internal/model"
	"eth_fetcher/internal/transaction"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
	"log"
	"math/big"
	"net/http"
	"strings"
)

type transactionHandler struct {
	txUseCase transaction.Fetcher
}

func NewTransactionHandler(txUseCase transaction.Fetcher) *transactionHandler {
	return &transactionHandler{
		txUseCase: txUseCase,
	}
}

func (th *transactionHandler) ListTransactionsByRLP( c *gin.Context)  {
	encodedRLP := c.Param("rlphex")

	//rlpEncodedData := "0xf90110b842307839623266366133633265316165643263636366393262613636366332326430353361643064386135646137616131666435343737646364363537376234353234b842307835613537653330353163623932653264343832353135623037653762336431383531373232613734363534363537626436346131346333396361336639636632b842307837316239653262343464343034393863303861363239383866616337373664306561633062356239363133633337663966366639613462383838613862303537b842307863356639366266316235346433333134343235643233373962643737643765643465363434663763366538343961373438333230323862333238643464373938"
	encodedRLP = strings.TrimPrefix(encodedRLP, "0x")
	decoded, err := hex.DecodeString(encodedRLP)
	if err != nil {
		log.Fatal(err)
	}

	var transactionHashes []string

	err = rlp.DecodeBytes(decoded, &transactionHashes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Hashes: %+v\n", transactionHashes)

	txs := th.txUseCase.FetchBlockchainTransactionsByHashes(transactionHashes)
	res := model.Transactions{
		Transactions: txs,
	}



	c.IndentedJSON(http.StatusOK, res)

}

func (th *transactionHandler) ListAllTransactions(c *gin.Context) {
	db := database.Init()
	// Ping the database to check the connection
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to ping the database:", err)
		return
	}

	fmt.Println("Connected to the PostgreSQL database")



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
			log.Fatal(err)
		}

		// Append the transaction to the slice
		transactions = append(transactions, transaction)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Print or process the retrieved transactions
	res := model.Transactions{
		Transactions: transactions,
	}

	r, _ := json.Marshal(res)
	fmt.Printf("\n\n\n%v\n", string(r))
}

func DecodeTransactionValue(valueHex string) (*big.Int, error) {
	value := new(big.Int)
	//todo use trim prefix
	value, success := value.SetString(valueHex[2:], 16) // Remove the "0x" prefix
	if !success {
		return nil, fmt.Errorf("failed to decode value")
	}

	return value, nil
}

