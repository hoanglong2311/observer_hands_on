package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/observer/app/internal/handler"
	"github.com/observer/app/internal/migrate"
	"github.com/observer/app/internal/repository"
	"github.com/observer/app/internal/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))
	slog.SetDefault(logger)

	if err := run(logger); err != nil {
		logger.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	dbURL := getenv("DATABASE_URL", "postgres://observer:observer@localhost:5432/observer?sslmode=disable")
	addr := getenv("ADDR", ":8080")

	// ---- Database ----
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("create db pool: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}
	logger.Info("database connected")

	if err := migrate.Run(ctx, pool); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	// ---- Dependency wiring ----
	taskRepo := repository.NewTaskRepository(pool)
	taskSvc := service.NewTaskService(taskRepo)
	taskH := handler.NewTaskHandler(taskSvc)
	healthH := handler.NewHealthHandler(pool)

	// ---- Router (Go 1.22+ pattern matching) ----
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthH.Healthz)
	mux.HandleFunc("GET /readyz", healthH.Readyz)
	mux.HandleFunc("GET /api/tasks", taskH.List)
	mux.HandleFunc("POST /api/tasks", taskH.Create)
	mux.HandleFunc("GET /api/tasks/{id}", taskH.GetByID)
	mux.HandleFunc("PATCH /api/tasks/{id}", taskH.UpdateStatus)

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ---- Graceful shutdown ----
	shutdownCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("server starting", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "error", err)
			stop()
		}
	}()

	<-shutdownCtx.Done()
	logger.Info("shutdown signal received")

	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		return fmt.Errorf("graceful shutdown: %w", err)
	}
	logger.Info("server stopped gracefully")
	return nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
