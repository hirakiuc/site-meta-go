package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hirakiuc/site-meta-go/server/api"
	"github.com/hirakiuc/site-meta-go/server/health"
)

func main() {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	apis := r.Group("/api")
	// TBD configure auth
	// api.Use(AuthRequired())
	{
		apis.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong, api",
			})
		})
		apis.GET("/meta", api.MetaHandler)
	}

	r.GET("/health", health.Handler)
	_ = r.Run()
}
