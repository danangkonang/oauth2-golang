package app

import (
	"log"
	"net/http"
	"time"

	"github.com/danangkonang/oauth2-golang/config"
	"github.com/danangkonang/oauth2-golang/helper"
	"github.com/danangkonang/oauth2-golang/router"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/gorilla/mux"
)

func Run() {
	manager := manage.NewDefaultManager()
	// manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.SetAuthorizeCodeTokenCfg(&manage.Config{
		AccessTokenExp:    time.Minute * 1,
		RefreshTokenExp:   time.Minute * 2,
		IsGenerateRefresh: true,
	})
	manager.MustTokenStorage(helper.NewMysqlTokenStore(config.Connection()))
	manager.MapAccessGenerate(helper.NewAccessGenerate())

	clientStore := helper.NewClientStore(config.Connection())
	clientStore.Set("client_id", &models.Client{
		ID:     "client_id",
		Secret: "client_secret",
		Domain: "http://localhost:9000",
	})
	manager.MapClientStorage(clientStore)

	r := mux.NewRouter()
	router.OauthRouter(r, manager, config.Connection())
	router.UserRouter(r, config.Connection())

	log.Printf("running on %s:%d%s", "http://localhost", 9096, "/oauth/token")
	server := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:9096",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
