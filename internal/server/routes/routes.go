package routes

import (
	"github.com/getz-devs/librakeeper-server/internal/server/handlers"
	"github.com/getz-devs/librakeeper-server/internal/server/middlewares"
	"github.com/gin-gonic/gin"
)

// Handlers is a struct that groups all handler functions.
type Handlers struct {
	Bookshelves *handlers.BookshelfHandlers
	Books       *handlers.BookHandlers
}

// SetupRoutes sets up the API routes for the server.
func SetupRoutes(router *gin.Engine, h *Handlers) {
	api := router.Group("/api")

	api.GET("/health", handlers.HealthCheck)

	// Bookshelf routes
	bookshelvesGroup := api.Group("/bookshelves")
	{
		bookshelvesGroup.GET("/", middlewares.AuthMiddleware(), h.Bookshelves.GetByUser)
		bookshelvesGroup.GET("/:id", middlewares.AuthMiddleware(), h.Bookshelves.GetByID)
		bookshelvesGroup.POST("/add", middlewares.AuthMiddleware(), h.Bookshelves.Create)
		bookshelvesGroup.PUT("/:id", middlewares.AuthMiddleware(), h.Bookshelves.Update)
		bookshelvesGroup.DELETE("/:id", middlewares.AuthMiddleware(), h.Bookshelves.Delete)
	}

	// Book routes
	booksGroup := api.Group("/books")
	{
		booksGroup.GET("/", middlewares.AuthMiddleware(), h.Books.GetByUser)
		booksGroup.GET("/:id", h.Books.GetByID)
		booksGroup.GET("/bookshelf/:id", middlewares.AuthMiddleware(), h.Books.GetByBookshelfID)
		booksGroup.POST("/", middlewares.AuthMiddleware(), h.Books.Create)
		booksGroup.PUT("/:id", middlewares.AuthMiddleware(), h.Books.Update)
		booksGroup.DELETE("/:id", middlewares.AuthMiddleware(), h.Books.Delete)
	}
}
