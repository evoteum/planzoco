package routes

import (
	"github.com/evoteum/planzoco/go/planzoco/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// Serve static files from the static directory
	r.Static("/static", "./static")

	// Basic routes for MVP
	r.GET("/", handlers.ListEvents)
	r.GET("/events/new", handlers.NewEventForm)
	r.POST("/events", handlers.CreateEvent)
	r.GET("/events/:id", handlers.GetEvent)
	r.POST("/events/:id/questions", handlers.CreateQuestion)
	r.POST("/questions/:id/options", handlers.CreateOption)
	r.POST("/options/:id/vote", handlers.VoteOption)
	r.GET("/questions/:id", handlers.GetQuestion)

	return r
}
