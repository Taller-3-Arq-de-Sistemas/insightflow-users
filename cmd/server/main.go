package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/config"
	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/adapters/sqlite"
	repository "github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/adapters/sqlite/sqlc"
	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/auth"
	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/users"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db, err := sqlite.New(cfg.DBUrl)
	if err != nil {
		logger.Error("Unable to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	repo := repository.New(db)
	usersService := users.NewService(repo)
	usersHandler := users.NewHandler(usersService)

	authService := auth.NewService(repo)
	authHandler := auth.NewHandler(authService)

	api := application{
		config:       cfg,
		usersHandler: usersHandler,
		authHandler:  authHandler,
	}

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      api.mount(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  time.Minute,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logger.Error("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			logger.Error("server shutdown error", "error", err)
		}
		serverStopCtx()
	}()

	// Run the server
	logger.Info("server starting", "port", cfg.Port)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Error("server start error", "error", err)
		os.Exit(1)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	logger.Info("server stopped")
}
