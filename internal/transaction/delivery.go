package transaction

import "github.com/gin-gonic/gin"

type Explorer interface {
	ExploreTransactionsByRLP(c *gin.Context)
	ExploreAllTransactions(c *gin.Context)
	ShowTransactionHistory(c *gin.Context)
}
