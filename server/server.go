package server

import (
	"context"
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
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Run() {
	cfg := config.New()
	l := logger.Init(cfg.AppMode)
	db := database.Init(cfg, l)

	alchemy := api.NewAlchemyAPI(cfg, l)

	//transactions
	tr := repository.NewTransactionRepository(db)
	tuc := usecase.NewTransactionUseCase(alchemy, tr, l)
	h := httpHandler.NewTransactionHandler(tuc, l)

	//auth
	ar := repository2.NewAuthRepository(db)
	auc := usecase2.NewAuthUseCase(l,ar)
	ah := http2.NewAuthHandler(l, auc)

	router := gin.Default()
	limeAPI := router.Group("/lime")
	limeAPI.GET("/all", h.ExploreAllTransactions)
	limeAPI.GET("/eth/:rlphex",AuthMiddleware(), h.ExploreTransactionsByRLP)
	limeAPI.GET("/my",AuthMiddleware(), h.ShowTransactionHistory)
	limeAPI.POST("/authenticate",ah.Authenticate)


	port := fmt.Sprintf(":%s", cfg.APIPort)

	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatalf("Server failed to start: %v", err)
		}
		l.Infow("listening on port", "port", cfg.APIPort)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		l.Info("here")
		l.Fatal("Server shutdown failed: %v", err)
	}

	l.Info("Server stopped gracefully")
}

