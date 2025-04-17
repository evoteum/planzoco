package handlers

import (
	"net/http"
	"github.com/evoteum/planzoco/go/planzoco/databases"
	"github.com/evoteum/planzoco/go/planzoco/models"
	"github.com/evoteum/planzoco/go/planzoco/utils"

	"github.com/gin-gonic/gin"
)

func CreateQuestion(c *gin.Context) {
	eventID := c.Param("id")

	var question models.Question
	if err := c.ShouldBind(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := utils.GenerateID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ID"})
		return
	}
	question.ID = id
	question.EventID = eventID

	if err := databases.AddQuestion(eventID, question); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save question"})
		return
	}

	c.Redirect(http.StatusFound, "/events/"+eventID)
}

func GetQuestion(c *gin.Context) {
	questionID := c.Param("id")

	question, event, err := databases.GetQuestion(questionID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title":   "Error",
			"message": "Failed to fetch question",
		})
		return
	}

	c.HTML(http.StatusOK, "question.html", gin.H{
		"event":    event,
		"question": question,
	})
}
