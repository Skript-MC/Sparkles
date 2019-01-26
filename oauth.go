package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

type User struct {
	Email string `json:"email"`
	Group group  `json:"primaryGroup"`
}

type group struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func contains(a []int, x int) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func Login(c *gin.Context) {

	state := randToken()
	session := sessions.Default(c)
	session.Set("state", state)
	session.Save()

	url := conf.AuthCodeURL(state)

	c.HTML(http.StatusOK, "login.html", gin.H{
		"url": url,
	})
}

func Callback(c *gin.Context) {
	session := sessions.Default(c)

	code := c.Query("code")
	state := session.Get("state")

	if state != c.Query("state") {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{
			"error": "Invalid state",
		})
		return
	}

	tok, err := conf.Exchange(ctx, code)

	if err != nil {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{
			"error": err,
		})
		return
	}

	client := conf.Client(ctx, tok)
	response, err := client.Get("https://skript-mc.fr/forum/api/core/me")

	if err != nil {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{
			"error": err,
		})
		return
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	var user User
	err = json.Unmarshal(data, &user)

	if err != nil || !contains(C.Groups, user.Group.Id) {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{
			"error": "access denied",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		Issuer:    "Sparkles",
	})
	signedToken, _ := token.SignedString(signKey)
	c.SetCookie("sparkles_auth", signedToken, 3600*24*7, "/", C.Domain, true, true)

	c.Redirect(http.StatusTemporaryRedirect, C.Redirect)
}
