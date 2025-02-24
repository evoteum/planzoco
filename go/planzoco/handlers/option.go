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

func VoteOption(c *gin.Context) {
	optionID := c.Param("id")

	if err := databases.VoteOption(optionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record vote"})
		return
	}

	c.Redirect(http.StatusFound, "/questions/"+optionID)
}
