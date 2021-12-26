package router

import (
	"github.com/danangkonang/oauth2/config"
	"github.com/danangkonang/oauth2/controller"
	"github.com/danangkonang/oauth2/service"
	"github.com/gorilla/mux"
)

func UserRouter(router *mux.Router, db *config.DB) {
	rest := controller.NewUserController(
		service.NewServiceUser(db),
	)
	v1 := router.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/user/register", rest.Register).Methods("GET")
}