package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vasyukov1/ai-doc-gen/internal/config"
	"github.com/vasyukov1/ai-doc-gen/internal/documents"
	"github.com/vasyukov1/ai-doc-gen/internal/pgqueue"
	"github.com/vasyukov1/ai-doc-gen/internal/platform"
	"github.com/vasyukov1/ai-doc-gen/internal/storage"
	"github.com/vasyukov1/ai-doc-gen/internal/web"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	db, err := storage.Open(connectCtx, cfg.DatabaseURL)
	cancel()
	if err != nil {
		return err
	}
	defer db.Close()

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		migrateCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		return storage.Migrate(migrateCtx, db)
	}

	jobs := storage.NewRepository(db)
	queue := pgqueue.New(db)

	if len(os.Args) == 1 || os.Args[1] == "serve" {
		return serve(ctx, cfg, jobs)
	}
	if os.Args[1] == "cleanup" {
		return documents.NewCleaner(jobs, cfg.DataDir, cfg.UploadDir, cfg.SourceDir).Once(ctx)
	}
	if os.Args[1] == "worker" {
		go documents.NewCleaner(jobs, cfg.DataDir, cfg.UploadDir, cfg.SourceDir).Run(ctx)
		owner, err := platform.NewID()
		if err != nil {
			return err
		}

		log.Printf("worker started: %s", owner)

		var fixture documents.Provider
		if os.Getenv("LLM_FIXTURE_MODE") == "1" {
			fixture = &documents.ChatProvider{Endpoint: os.Getenv("LLM_ENDPOINT"), Model: os.Getenv("LLM_MODEL"), Key: os.Getenv("LLM_API_KEY")}
		}
		return documents.NewWorker(queue, jobs, owner, cfg.SourceDir, cfg.GitHubToken, fixture, cfg.LLMEncryptionKey).Run(ctx)
	}

	return fmt.Errorf("unknown command: %s", os.Args[1])
}

func serve(ctx context.Context, cfg config.Config, jobs *storage.Repository) error {
	handler, err := web.NewHandler(jobs, cfg.UploadDir, cfg.LLMEncryptionKey, cfg.Retention)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown: %v", err)
		}
	}()

	log.Printf("listening on %s", cfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve HTTP: %w", err)
	}

	return nil
}
