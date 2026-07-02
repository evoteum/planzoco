package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/evoteum/planzoco/models"

	"github.com/gin-gonic/gin"
)

const voterNameCookie = "planzoco_voter_name"

// voterName returns the name this browser is currently voting as, or "" if unset.
func voterName(c *gin.Context) string {
	name, err := c.Cookie(voterNameCookie)
	if err != nil {
		return ""
	}
	return name
}

// requireVoterName redirects to the name-gate form if the browser has no voter name set,
// passing next as the URL to return to once a name is given. It returns true if a redirect happened.
func requireVoterName(c *gin.Context, next string) bool {
	if voterName(c) != "" {
		return false
	}

	c.Redirect(http.StatusFound, "/whoami?next="+url.QueryEscape(next))
	return true
}

func WhoAmIForm(c *gin.Context) {
	next := c.Query("next")
	if next == "" {
		next = "/"
	}

	c.HTML(http.StatusOK, "whoami.html", gin.H{"next": next})
}

func SetVoterName(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	next := c.PostForm("next")
	if next == "" {
		next = "/"
	}

	if name == "" || len(name) > models.MaxTextLength {
		c.HTML(http.StatusBadRequest, "whoami.html", gin.H{
			"next":  next,
			"error": "Please enter a name, up to 255 characters.",
		})
		return
	}

	c.SetCookie(voterNameCookie, name, 365*24*3600, "/", "", false, true)
	c.Redirect(http.StatusFound, next)
}
