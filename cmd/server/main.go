package main

import (
	"context"
	"log"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"exchange/internal/adapters/config"
	"exchange/internal/adapters/database"
	"exchange/internal/domain/transaction"
	"exchange/internal/domain/wallet"
	"exchange/internal/ports/http"
	"exchange/internal/ports/persistence"
	"exchange/internal/usecase"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.NewPostgresDB(database.PostgresConfig{
		Host:     cfg.Postgre.Host,
		Port:     cfg.Postgre.Port,
		User:     cfg.Postgre.User,
		Password: cfg.Postgre.Password,
		DBName:   cfg.Postgre.DBName,
		SSLMode:  cfg.Postgre.SSLMode,
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	walletRepo := persistence.NewPostgresWalletRepository(db)
	transactionRepo := persistence.NewPostgresTransactionRepository(db)

	walletService := wallet.NewWalletService(walletRepo)
	transactionService := transaction.NewTransactionService(transactionRepo)

	txManager := persistence.NewPostgresTransactionManager(db)

	walletUC := usecase.NewWalletUseCase(walletService, transactionService, txManager)

	handler := http.NewHandler(walletUC)
	router := http.NewRouter(handler)

	srv := &nethttp.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		cancel()
	}()

	go func() {
		log.Printf("Starting server on %s", cfg.Server.Address)
		if err := srv.ListenAndServe(); err != nil && err != nethttp.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Server exited properly.")
}
