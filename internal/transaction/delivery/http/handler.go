package http

import (
	"encoding/hex"
	"errors"
	"eth_fetcher/infrastructure/logger"
	"eth_fetcher/internal/model"
	"eth_fetcher/internal/transaction"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

type transactionHandler struct {
	txUseCase transaction.Fetcher
	l         logger.ILogger
}

func NewTransactionHandler(txUseCase transaction.Fetcher, l logger.ILogger) transaction.Explorer {
	return &transactionHandler{
		txUseCase: txUseCase,
		l:         l,
	}
}

func (th *transactionHandler) ExploreTransactionsByRLP(c *gin.Context) {
	//test
	authToken := c.GetHeader("AUTH_TOKEN")
	var user string
	if authToken != "" {
		// Verify and parse the JWT token
		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			// Replace "secret" with your own secret key used for signing the token
			return []byte("secret"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Retrieve the username from the token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		user = claims["username"].(string)
	}

	//test
	encodedRLP := c.Param("rlphex")
	encodedRLP = strings.TrimPrefix(encodedRLP, "0x")
	decoded, err := hex.DecodeString(encodedRLP)
	if err != nil {
		th.l.Infow("invalid rlp input", "rlp", encodedRLP, "error", err)

		c.IndentedJSON(http.StatusBadRequest, ErrorResponse{
			Message: "invalid rlp",
		})
		return
	}

	var transactionHashes []string

	err = rlp.DecodeBytes(decoded, &transactionHashes)
	if err != nil {
		th.l.Infow("decoding rlp", "rlp_bytes", decoded, "error", err)

		c.IndentedJSON(http.StatusBadRequest, ErrorResponse{
			Message: "invalid rlp",
		})
		return
	}

	txs, err := th.txUseCase.FetchBlockchainTransactionsByHashes(transactionHashes)
	if err != nil {
		th.l.Infow("fetching transaction errors", "error", err)
	}

	res := model.Transactions{
		Transactions: txs,
	}

	th.l.Infow("transactions listed successfully", "rlp", encodedRLP)
	c.IndentedJSON(http.StatusOK, res)
	fmt.Printf("USER %s requested %v",user, transactionHashes)

	return
}

func (th *transactionHandler) ExploreAllTransactions(c *gin.Context) {
	txs, err := th.txUseCase.ListRequestedTransactions()
	res := model.Transactions{
		Transactions: txs,
	}

	if err != nil {
		th.l.Infow("listing all transactions", "error", err)
		if errors.Is(err, ErrNoRecords) {
			c.IndentedJSON(http.StatusNotFound, res)
			return
		}

		c.IndentedJSON(http.StatusInternalServerError, res)
		return
	}

	c.IndentedJSON(http.StatusOK, res)
	return
}
