package routes

import (
	"github.com/getz-devs/librakeeper-server/internal/server/handlers"
	"github.com/getz-devs/librakeeper-server/internal/server/middlewares"
	"github.com/gin-gonic/gin"
)

// Handlers is a struct that groups all handler functions.
type Handlers struct {
	Users       *handlers.UserHandlers
	Bookshelves *handlers.BookshelfHandlers
	Books       *handlers.BookHandlers
}

// SetupRoutes sets up the API routes for the server.
func SetupRoutes(router *gin.Engine, h *Handlers) {
	api := router.Group("/api")

	api.GET("/health", handlers.HealthCheck)

	// User routes
	userGroup := api.Group("/users")
	{
		userGroup.POST("/", middlewares.AuthMiddleware(), h.Users.Create)
		userGroup.GET("/:id", middlewares.AuthMiddleware(), h.Users.GetByID)
		userGroup.PUT("/:id", middlewares.AuthMiddleware(), h.Users.Update)
		userGroup.DELETE("/:id", middlewares.AuthMiddleware(), h.Users.Delete)
	}

	// Bookshelf routes
	bookshelvesGroup := api.Group("/bookshelves")
	{
		bookshelvesGroup.POST("/", middlewares.AuthMiddleware(), h.Bookshelves.Create)
		bookshelvesGroup.GET("/:id", middlewares.AuthMiddleware(), h.Bookshelves.GetByID)
		bookshelvesGroup.GET("/user", middlewares.AuthMiddleware(), h.Bookshelves.GetByUser)
		bookshelvesGroup.PUT("/:id", middlewares.AuthMiddleware(), h.Bookshelves.Update)
		bookshelvesGroup.DELETE("/:id", middlewares.AuthMiddleware(), h.Bookshelves.Delete)
	}

	// Book routes
	booksGroup := api.Group("/books")
	{
		booksGroup.POST("/", middlewares.AuthMiddleware(), h.Books.Create)
		booksGroup.GET("/:id", h.Books.GetByID)
		booksGroup.GET("/user", middlewares.AuthMiddleware(), h.Books.GetByUser)
		booksGroup.GET("/bookshelf/:bookshelfId", middlewares.AuthMiddleware(), h.Books.GetByBookshelfID)
		booksGroup.PUT("/:id", middlewares.AuthMiddleware(), h.Books.Update)
		booksGroup.DELETE("/:id", middlewares.AuthMiddleware(), h.Books.Delete)
	}
}
