package server

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
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server represents the API server.
type Server struct {
	config     *config.Config
	log        *slog.Logger
	router     *gin.Engine
	httpServer *http.Server
}

// NewServer creates a new Server instance.
func NewServer(config *config.Config, log *slog.Logger) *Server {
	return &Server{
		config: config,
		log:    log,
		router: gin.New(),
	}
}

// Initialize initializes the server components.
func (s *Server) Initialize() error {
	// Initialize Firebase
	err := auth.InitializeFirebase(s.config.FirebaseConfigPath)
	if err != nil {
		return fmt.Errorf("failed to initialize Firebase: %w", err)
	}

	_, collections := storage.Initialize(s.config, s.log)

	userService := services.NewUserService(collections.UsersCollection, s.log)
	bookshelfService := services.NewBookshelfService(collections.BookshelvesCollection, s.log)
	bookService := services.NewBookService(collections.BooksCollection, s.log)

	h := &routes.Handlers{
		Users:       handlers.NewUserHandlers(userService, s.log),
		Bookshelves: handlers.NewBookshelfHandlers(bookshelfService, s.log),
		Books:       handlers.NewBookHandlers(bookService, s.log),
	}

	// Configure CORS
	corsConfig := cors.Config{
		AllowOrigins:     s.config.Server.AllowedOrigins, // Get origins from config
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}
	s.router.Use(gin.Logger(), gin.Recovery(), cors.New(corsConfig))
	routes.SetupRoutes(s.router, h)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Server.Port),
		Handler: s.router,
	}

	return nil
}

// Run starts the HTTP server and handles graceful shutdown.
func (s *Server) Run() error {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error("failed to start server", slog.Any("error", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}
