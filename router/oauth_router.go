package router

import (
	"github.com/danangkonang/oauth2-golang/config"
	"github.com/danangkonang/oauth2-golang/controller"
	"github.com/danangkonang/oauth2-golang/middleware"
	"github.com/danangkonang/oauth2-golang/service"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/gorilla/mux"
)

func OauthRouter(router *mux.Router, manager *manage.Manager, db *config.DB) {
	c := controller.NewOauthController(
		manager,
		service.NewServiceUser(db),
	)
	m := middleware.NewOauthMiddleware(manager)
	v1 := router.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/oauth/token", c.Login).Methods("POST")
	v1.HandleFunc("/secure", m.Auth(c.Secure)).Methods("POST")
}
