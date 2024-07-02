package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/server/config"
	"github.com/getz-devs/librakeeper-server/internal/server/handlers"
	"github.com/getz-devs/librakeeper-server/internal/server/routes"
	"github.com/getz-devs/librakeeper-server/internal/server/services"
	"github.com/getz-devs/librakeeper-server/internal/server/services/auth"
	"github.com/getz-devs/librakeeper-server/internal/server/services/storage"
	"github.com/getz-devs/librakeeper-server/lib/prettylog"
	"github.com/gin-contrib/cors"
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

	err := auth.InitializeFirebase(cfg.FirebaseConfigPath)
	if err != nil {
		log.Error("failed to initialize Firebase", slog.Any("error", err))
		// TODO: Handle the error appropriately (e.g., panic or graceful shutdown)
	}

	_, collections := storage.Initialize(cfg, log)

	userService := services.NewUserService(collections.UsersCollection, log)
	bookshelfService := services.NewBookshelfService(collections.BookshelvesCollection, log)
	bookService := services.NewBookService(collections.BooksCollection, log)

	h := &routes.Handlers{
		Users:       handlers.NewUserHandlers(userService, log),
		Bookshelves: handlers.NewBookshelfHandlers(bookshelfService, log),
		Books:       handlers.NewBookHandlers(bookService, log),
	}

	// Configure CORS
	corsConfig := cors.Config{
		AllowOrigins:     []string{"https://libra.potat.dev", "http://localhost:3000"}, // Allow specific origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},          // Allow methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},          // Allow headers including Authorization
		AllowCredentials: true,                                                         // Allow credentials
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), cors.New(corsConfig))

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
