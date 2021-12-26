package router

import (
	"github.com/danangkonang/oauth2/controller"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/gorilla/mux"
)

func OauthRouter(router *mux.Router, manager *manage.Manager) {
	c := controller.NewOauthController(manager)
	v1 := router.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/oauth/token", c.Login).Methods("POST")
	v1.HandleFunc("/secure", c.Secure).Methods("POST")
}