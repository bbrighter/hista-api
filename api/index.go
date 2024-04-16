// api-users.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Creating new gin engine
	g := gin.New()
	// Your handler logic
	g.GET("/api/users", func(c *gin.Context) {
		// ... your handler logic here ...
		c.JSON(200, gin.H{
			"message": "Hello from /api/users",
		})
	})

	// running gin engine
	g.ServeHTTP(w, r)
}
