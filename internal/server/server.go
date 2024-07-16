package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/getz-devs/librakeeper-server/internal/server/auth"
	"github.com/getz-devs/librakeeper-server/internal/server/config"
	"github.com/getz-devs/librakeeper-server/internal/server/handlers"
	"github.com/getz-devs/librakeeper-server/internal/server/routes"
	"github.com/getz-devs/librakeeper-server/internal/server/services/book"
	"github.com/getz-devs/librakeeper-server/internal/server/services/bookshelf"
	"github.com/getz-devs/librakeeper-server/internal/server/services/search"
	"github.com/getz-devs/librakeeper-server/internal/server/services/storage"
	"github.com/getz-devs/librakeeper-server/internal/server/storage/mongo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	err := auth.InitializeFirebase(s.config.Auth.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to initialize Firebase: %w", err)
	}

	db, err := storage.Initialize(s.config, s.log)
	if err != nil {
		return fmt.Errorf("failed to initialize Database: %w", err)
	}

	conn, err := grpc.NewClient(
		s.config.GRPC.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	bookRepo := mongo.NewBookRepo(db, s.log, "user_books")
	allBooksRepo := mongo.NewBookRepo(db, s.log, "all_books")
	bookshelfRepo := mongo.NewBookshelfRepo(db, s.log)
	searcherClient := search.NewSearcherClient(conn, s.log)

	bookService := book.NewBookService(bookRepo, allBooksRepo, bookshelfRepo, s.log)
	bookshelfService := bookshelf.NewBookshelfService(bookshelfRepo, s.log)
	searchService := search.NewSearchService(searcherClient, allBooksRepo, s.log)

	h := &routes.Handlers{
		Books:       handlers.NewBookHandlers(bookService, s.log),
		Bookshelves: handlers.NewBookshelfHandlers(bookshelfService, s.log),
		Search:      handlers.NewSearchHandlers(searchService, s.log),
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
