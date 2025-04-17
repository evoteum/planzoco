package github.com/evoteum/planzoco/go/planzoco/handlers

import (
	"fmt"
	"net/http"
	"github.com/evoteum/planzoco/go/planzoco/databases"
	"github.com/evoteum/planzoco/go/planzoco/models"
	"github.com/evoteum/planzoco/go/planzoco/utils"

	"github.com/gin-gonic/gin"
)

func ListEvents(c *gin.Context) {
	events, err := databases.ListEvents()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title":   "Error",
			"message": "Failed to fetch events",
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"events": events,
	})
}

func NewEventForm(c *gin.Context) {
	c.HTML(http.StatusOK, "new_event.html", nil)
}

func CreateEvent(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBind(&event); err != nil {
		c.HTML(http.StatusBadRequest, "new_event.html", gin.H{"error": err.Error()})
		return
	}

	id, err := utils.GenerateID()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "new_event.html", gin.H{"error": "Failed to generate ID"})
		return
	}

	event.ID = id

	if err := databases.CreateEvent(event); err != nil {
		c.HTML(http.StatusInternalServerError, "new_event.html", gin.H{"error": "Failed to save event"})
		return
	}

	c.Redirect(http.StatusFound, "/events/"+event.ID)
}

func GetEvent(c *gin.Context) {
	eventID := c.Param("id")

	event, err := databases.GetEvent(eventID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title":   "Error",
			"message": "Failed to fetch event",
		})
		return
	}

	if event == nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title":   "Not Found",
			"message": "Event not found",
		})
		return
	}

	scheme := getScheme(c)
	baseURL := fmt.Sprintf("%s://%s", scheme, c.Request.Host)
	c.HTML(http.StatusOK, "event.html", gin.H{
		"event":   event,
		"baseURL": baseURL,
	})
}

func getScheme(c *gin.Context) string {
	if forwardedProto := c.Request.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
		return forwardedProto
	}

	if urlScheme := c.Request.URL.Scheme; urlScheme != "" {
		return urlScheme
	}

	return "https"
}
