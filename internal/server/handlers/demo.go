package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DemoHandler(c *gin.Context) {
	foo := c.Query("foo")
	if foo != "bar" {
		log.Printf("foo: %v", foo)
		c.JSON(http.StatusTeapot, gin.H{"error": "Not Happy :("})
		return
	}

	c.JSON(http.StatusOK, gin.H{"info": "Happy :)"})
}

/*
func CreateTodo(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var body struct {
		Todo string `json:"todo"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Save the todo for the user with UID 'uid'
	c.JSON(http.StatusOK, gin.H{"todo": body.Todo})
}
*/
