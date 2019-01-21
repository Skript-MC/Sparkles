package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authorize(c *gin.Context) {

	cookie, _ := c.Cookie("sparkles_auth")
	fmt.Println(cookie)
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return signKey, nil
	})

	if err != nil || !token.Valid {
		c.String(http.StatusUnauthorized, "")
		return
	}

	c.String(http.StatusOK, "")
}
