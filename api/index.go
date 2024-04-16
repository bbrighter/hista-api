// api-users.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bbrighter/hista-api/api/repository"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Request-Method", "*")
		c.Writer.Header().Set("Access-Control-Content-Type", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, accept, origin, Cache-Control, If-None-Match, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PATCH, DELETE, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	repo := repository.InitRepository()
	type Version struct {
		Version string
	}
	var version Version
	repo.Db.Raw("select version()").Scan(&version)

	// Creating new gin engine
	g := gin.New()
	g.Use(CORSMiddleware())
	// Your handler logic
	g.GET("/api/users", func(c *gin.Context) {
		// ... your handler logic here ...
		c.JSON(200, gin.H{
			"message": version.Version,
		})
	})

	// running gin engine
	g.ServeHTTP(w, r)
}
