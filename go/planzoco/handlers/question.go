package handlers

import (
	"net/http"
	"planzoco/databases"
	"planzoco/models"
	"planzoco/utils"

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

	question, event, err := databases.GetQuestionWithEvent(questionID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title":   "Error",
			"message": "Failed to fetch question",
		})
		return
	}

	if question == nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title":   "Not Found",
			"message": "Question not found",
		})
		return
	}

	c.HTML(http.StatusOK, "question.html", gin.H{
		"event":    event,
		"question": question,
	})
}

func UpdateQuestionForm(c *gin.Context) {
	questionID := c.Param("id")

	question, event, err := databases.GetQuestionWithEvent(questionID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title":   "Error",
			"message": "Failed to fetch question",
		})
		return
	}

	if question == nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title":   "Not Found",
			"message": "Question not found",
		})
		return
	}

	c.HTML(http.StatusOK, "edit_question.html", gin.H{
		"event":    event,
		"question": question,
	})
}

func UpdateQuestion(c *gin.Context) {
	questionID := c.Param("id")

	// Get existing question to preserve eventID
	existingQuestion, _, err := databases.GetQuestionWithEvent(questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch question"})
		return
	}

	if existingQuestion == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	var question models.Question
	if err := c.ShouldBind(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Preserve ID and EventID
	question.ID = questionID
	question.EventID = existingQuestion.EventID

	// Preserve existing options
	question.Options = existingQuestion.Options

	if err := databases.UpdateQuestion(question); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update question"})
		return
	}

	c.Redirect(http.StatusFound, "/questions/"+questionID)
}

func DeleteQuestion(c *gin.Context) {
	questionID := c.Param("id")

	// Get the question first to know which event to redirect to
	question, _, err := databases.GetQuestionWithEvent(questionID)
	if err != nil || question == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch question"})
		return
	}

	eventID := question.EventID

	if err := databases.DeleteQuestion(questionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete question"})
		return
	}

	c.Redirect(http.StatusFound, "/events/"+eventID)
}
