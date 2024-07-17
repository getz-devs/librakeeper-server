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
	Search      *handlers.SearchHandlers
}

// SetupRoutes sets up the API routes for the server.
func SetupRoutes(router *gin.Engine, h *Handlers) {
	api := router.Group("/api")

	api.GET("/health", handlers.HealthCheck)

	// Book routes
	booksGroup := api.Group("/books")
	{
		booksGroup.POST("/add", middlewares.AuthMiddleware(), h.Books.Create)
		booksGroup.GET("/", middlewares.AuthMiddleware(), h.Books.GetByUser)
		booksGroup.GET("/:id", middlewares.AuthMiddleware(), h.Books.GetByID)
		booksGroup.GET("/isbn/:isbn", middlewares.AuthMiddleware(), h.Books.GetByISBN)
		booksGroup.GET("/bookshelf/:id", middlewares.AuthMiddleware(), h.Books.GetByBookshelfID)
		booksGroup.PUT("/:id", middlewares.AuthMiddleware(), h.Books.Update)
		booksGroup.DELETE("/:id", middlewares.AuthMiddleware(), h.Books.Delete)
	}

	// Bookshelf routes
	bookshelvesGroup := api.Group("/bookshelves")
	{
		bookshelvesGroup.POST("/add", middlewares.AuthMiddleware(), h.Bookshelves.Create)
		bookshelvesGroup.GET("/", middlewares.AuthMiddleware(), h.Bookshelves.GetByUser)
		bookshelvesGroup.GET("/:id", middlewares.AuthMiddleware(), h.Bookshelves.GetByID)
		bookshelvesGroup.PUT("/:id", middlewares.AuthMiddleware(), h.Bookshelves.Update)
		bookshelvesGroup.DELETE("/:id", middlewares.AuthMiddleware(), h.Bookshelves.Delete)
	}

	searchGroup := api.Group("/search")
	{
		searchGroup.GET("/simple", middlewares.AuthMiddleware(), h.Search.Simple)
		searchGroup.GET("/advanced", middlewares.AuthMiddleware(), h.Search.Advanced)
	}
}
