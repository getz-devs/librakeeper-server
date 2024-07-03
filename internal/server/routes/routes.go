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

func SetupRoutes(router *gin.Engine, h *Handlers) {
	api := router.Group("/api")

	api.GET("/health", handlers.HealthCheck)

	userGroup := api.Group("/users")
	{
		userGroup.POST("/", middlewares.AuthMiddleware(), h.Users.CreateUser)
		userGroup.GET("/:id", middlewares.AuthMiddleware(), h.Users.GetUser)
		userGroup.PUT("/:id", middlewares.AuthMiddleware(), h.Users.UpdateUser)
		userGroup.DELETE("/:id", middlewares.AuthMiddleware(), h.Users.DeleteUser)
	}

	bookshelvesGroup := api.Group("/bookshelves")
	{
		bookshelvesGroup.POST("/", middlewares.AuthMiddleware(), h.Bookshelves.CreateBookshelf)
		bookshelvesGroup.GET("/:id", h.Bookshelves.GetBookshelf)
		bookshelvesGroup.GET("/user", middlewares.AuthMiddleware(), h.Bookshelves.GetBookshelvesByUserID)
		bookshelvesGroup.PUT("/:id", middlewares.AuthMiddleware(), h.Bookshelves.UpdateBookshelf)
		bookshelvesGroup.DELETE("/:id", middlewares.AuthMiddleware(), h.Bookshelves.DeleteBookshelf)
	}

	booksGroup := api.Group("/books")
	{
		booksGroup.POST("/", h.Books.CreateBook)
		booksGroup.GET("/:id", h.Books.GetBook)
		booksGroup.GET("/", h.Books.GetBooks)
		booksGroup.GET("/bookshelf/:bookshelfId", h.Books.GetBooksByBookshelfID)
		booksGroup.PUT("/:id", h.Books.UpdateBook)
		booksGroup.DELETE("/:id", h.Books.DeleteBook)
	}
}
