package handlers

import (
	"net/http"
	"planzoco/databases"
	"planzoco/models"
	"planzoco/utils"

	"github.com/gin-gonic/gin"
)

func CreateOption(c *gin.Context) {
	questionID := c.Param("id")

	var option models.Option
	if err := c.ShouldBind(&option); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := utils.GenerateID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ID"})
		return
	}
	option.ID = id
	option.QuestionID = questionID
	option.Votes = 0

	if err := databases.AddOption(questionID, option); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save option"})
		return
	}

	c.Redirect(http.StatusFound, "/questions/"+questionID)
}

func UpdateOptionForm(c *gin.Context) {
	optionID := c.Param("id")

	option, err := databases.GetOption(optionID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title":   "Error",
			"message": "Failed to fetch option",
		})
		return
	}

	if option == nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title":   "Not Found",
			"message": "Option not found",
		})
		return
	}

	// Get the question for context
	question, _, err := databases.GetQuestionWithEvent(option.QuestionID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title":   "Error",
			"message": "Failed to fetch question",
		})
		return
	}

	c.HTML(http.StatusOK, "edit_option.html", gin.H{
		"option":   option,
		"question": question,
	})
}

func UpdateOption(c *gin.Context) {
	optionID := c.Param("id")

	// Get existing option to preserve question ID and votes
	existingOption, err := databases.GetOption(optionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch option"})
		return
	}

	if existingOption == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Option not found"})
		return
	}

	var option models.Option
	if err := c.ShouldBind(&option); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Preserve ID, QuestionID, and votes
	option.ID = optionID
	option.QuestionID = existingOption.QuestionID
	option.Votes = existingOption.Votes

	if err := databases.UpdateOption(option); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update option"})
		return
	}

	c.Redirect(http.StatusFound, "/questions/"+option.QuestionID)
}

func DeleteOption(c *gin.Context) {
	optionID := c.Param("id")

	// Get the option first to know which question to redirect to
	option, err := databases.GetOption(optionID)
	if err != nil || option == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch option"})
		return
	}

	questionID := option.QuestionID

	if err := databases.DeleteOption(optionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete option"})
		return
	}

	c.Redirect(http.StatusFound, "/questions/"+questionID)
}

func VoteOption(c *gin.Context) {
	optionID := c.Param("id")

	// Get the option to find its question
	option, err := databases.GetOption(optionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch option"})
		return
	}

	if option == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Option not found"})
		return
	}

	questionID := option.QuestionID

	if err := databases.VoteOption(optionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record vote"})
		return
	}

	c.Redirect(http.StatusFound, "/questions/"+questionID)
}
