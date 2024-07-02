package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/server/config"
	"github.com/getz-devs/librakeeper-server/internal/server/handlers"
	"github.com/getz-devs/librakeeper-server/internal/server/routes"
	"github.com/getz-devs/librakeeper-server/internal/server/services"
	"github.com/getz-devs/librakeeper-server/internal/server/services/storage"
	"github.com/getz-devs/librakeeper-server/lib/prettylog"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting librakeeper server", slog.String("env", cfg.Env), slog.Int("port", cfg.Server.Port))

	_, collections := storage.Initialize(cfg, log)

	userService := services.NewUserService(collections.UsersCollection, log)
	bookshelfService := services.NewBookshelfService(collections.BookshelvesCollection, log)
	bookService := services.NewBookService(collections.BooksCollection, log)

	h := &routes.Handlers{
		Users:       handlers.NewUserHandlers(userService, log),
		Bookshelves: handlers.NewBookshelfHandlers(bookshelfService, log),
		Books:       handlers.NewBookHandlers(bookService, log),
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	routes.SetupRoutes(router, h)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start server", slog.Any("error", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", slog.Any("error", err))
	}

	log.Info("server exiting")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(prettylog.NewHandler(&slog.HandlerOptions{
			Level:       slog.LevelDebug,
			AddSource:   false,
			ReplaceAttr: nil,
		}))
	case "development":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "production":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
