package server

import (
	"eth_fetcher/infrastructure/api"
	"eth_fetcher/infrastructure/config"
	"eth_fetcher/infrastructure/database"
	"eth_fetcher/infrastructure/logger"
	"eth_fetcher/internal/transaction/delivery/http"
	"eth_fetcher/internal/transaction/repository"
	"eth_fetcher/internal/transaction/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func Run() {
	cfg := config.New()
	l := logger.Init(cfg.AppMode)
	db := database.Init(cfg, l)

	alchemy := api.NewAlchemyAPI(cfg, l)
	tr := repository.NewTransactionRepository(db, l)
	tuc := usecase.NewTransactionUseCase(alchemy, tr, l)
	h := http.NewTransactionHandler(tuc, l)
	router := gin.Default()

	router.GET("/lime/all", h.ExploreAllTransactions)
	router.GET("/lime/eth/:rlphex", h.ExploreTransactionsByRLP)

	l.Infow("listening on port", "port", cfg.APIPort)

	port := fmt.Sprintf(":%s", cfg.APIPort)
	log.Fatal(router.Run(port))
}
