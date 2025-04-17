package github.com/evoteum/planzoco/go/planzoco/handlers

import (
	"net/http"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	// TODO: Implement actual Auth0 login flow
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

func Callback(c *gin.Context) {
	// TODO: Implement Auth0 callback handling
	session := sessions.Default(c)
	session.Set("user_id", "temp-user-id") // Temporary for testing
	session.Save()

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusTemporaryRedirect, "/login")
}
