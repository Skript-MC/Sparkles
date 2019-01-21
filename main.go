package main

import (
	"context"
	"encoding/json"
	"github.com/coreos/go-systemd/daemon"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Domain    string `json:"domain"`
	Address   string `json:"address"`
	Redirect  string `json:"redirect"`
	Groups    []int  `json:"groups"`
	StoreKey  string `json:"store_key"`
	CookieKey string `json:"cookie_key"`
	Oauth     Oauth  `json:"oauth"`
}

type Oauth struct {
	AuthUrl      string   `json:"auth_url"`
	TokenUrl     string   `json:"token_url"`
	RedirectUrl  string   `json:"redirect_url"`
	ClientId     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Scopes       []string `json:"scopes"`
}

var C Config
var ctx = context.Background()
var conf = &oauth2.Config{}
var signKey = []byte("")

func loadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		panic(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func main() {
	C = loadConfiguration("./config/config.json")
	conf = &oauth2.Config{
		ClientID:     C.Oauth.ClientId,
		ClientSecret: C.Oauth.ClientSecret,
		Scopes:       C.Oauth.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  C.Oauth.AuthUrl,
			TokenURL: C.Oauth.TokenUrl,
		},
		RedirectURL: C.Oauth.RedirectUrl,
	}

	signKey = []byte(C.CookieKey)
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	store := memstore.NewStore([]byte(C.StoreKey))
	router.Use(sessions.Sessions("sparkles_store", store))
	router.LoadHTMLGlob("templates/*")

	router.GET("/", Index)
	router.GET("/validate", Authorize)
	router.GET("/login", Login)
	router.GET("/logout", Logout)
	router.GET("/callback", Callback)

	daemon.SdNotify(false, "READY=1")
	router.Run(C.Address)

	go func() {
		interval, err := daemon.SdWatchdogEnabled(false)
		if err != nil || interval == 0 {
			return
		}
		for {
			_, err := http.Get(C.Address)
			if err == nil {
				daemon.SdNotify(false, "WATCHDOG=1")
			}
			time.Sleep(interval / 3)
		}
	}()

}
