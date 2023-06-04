package server

import (
	"eth_fetcher/infrastructure/api"
	"eth_fetcher/infrastructure/config"
	"eth_fetcher/infrastructure/database"
	"eth_fetcher/infrastructure/logger"
	http2 "eth_fetcher/internal/authentication/delivery/http"
	repository2 "eth_fetcher/internal/authentication/repository"
	usecase2 "eth_fetcher/internal/authentication/usecase"
	httpHandler "eth_fetcher/internal/transaction/delivery/http"
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
	tr := repository.NewTransactionRepository(db)
	tuc := usecase.NewTransactionUseCase(alchemy, tr, l)
	h := httpHandler.NewTransactionHandler(tuc, l)
	router := gin.Default()

	router.GET("/lime/all", h.ExploreAllTransactions)
	router.GET("/lime/eth/:rlphex", h.ExploreTransactionsByRLP)

	ar := repository2.NewAuthRepository(db)
	auc := usecase2.NewAuthUseCase(l,ar)
	ah := http2.NewAuthHandler(l, auc)
	router.POST("/lime/authenticate",ah.Authenticate)
	l.Infow("listening on port", "port", cfg.APIPort)

	port := fmt.Sprintf(":%s", cfg.APIPort)
	log.Fatal(router.Run(port))
}

