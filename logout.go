package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Logout(c *gin.Context) {

	c.SetCookie("sparkles_auth", "", 3600*24*7, "/", C.Domain, true, true)
	c.Redirect(http.StatusTemporaryRedirect, "/login")
}
