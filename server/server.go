package server

import (
	"eth_fetcher/infrastructure/api"
	"eth_fetcher/infrastructure/database"
	"eth_fetcher/internal/transaction/delivery/http"
	"eth_fetcher/internal/transaction/repository"
	"eth_fetcher/internal/transaction/usecase"
	"github.com/gin-gonic/gin"
	"log"
)

func Run() {
	db := database.Init()

	alchemy := api.NewAlchemyAPI()
	tuc := usecase.NewTransactionUseCase(alchemy)
	tr := repository.NewTransactionRepository(db)
	h := http.NewTransactionHandler(tuc, tr)
	router := gin.Default()

	//todo use routing groups
	router.GET("/lime/all", h.ListAllTransactions)
	router.GET("/lime/eth/:rlphex",h.ListTransactionsByRLP)

	//todo extract port in env
	log.Println("listening on port :%s", ":8080")

	log.Fatal(router.Run(":8080"))
}


