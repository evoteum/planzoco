package routes

import (
	"planzoco/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// Serve static files from the static directory
	r.Static("/static", "./static")

	// Event routes
	r.GET("/", handlers.ListEvents)
	r.GET("/events/new", handlers.NewEventForm)
	r.POST("/events", handlers.CreateEvent)
	r.GET("/events/:id", handlers.GetEvent)
	r.GET("/events/:id/edit", handlers.UpdateEventForm)
	r.POST("/events/:id", handlers.UpdateEvent)
	r.POST("/events/:id/delete", handlers.DeleteEvent)

	// Question routes
	r.POST("/events/:id/questions", handlers.CreateQuestion)
	r.GET("/questions/:id", handlers.GetQuestion)
	r.GET("/questions/:id/edit", handlers.UpdateQuestionForm)
	r.POST("/questions/:id", handlers.UpdateQuestion)
	r.POST("/questions/:id/delete", handlers.DeleteQuestion)

	// Option routes
	r.POST("/questions/:id/options", handlers.CreateOption)
	r.GET("/options/:id/edit", handlers.UpdateOptionForm)
	r.POST("/options/:id", handlers.UpdateOption)
	r.POST("/options/:id/delete", handlers.DeleteOption)
	r.POST("/options/:id/vote", handlers.VoteOption)

	r.GET("/health", handlers.HealthCheck)


	return r
}
