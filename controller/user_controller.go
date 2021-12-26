package controller

import (
	"net/http"

	"github.com/danangkonang/oauth2/model"
	"github.com/danangkonang/oauth2/service"
)

type userController struct {
	Service service.UserService
}

func NewUserController(user service.UserService) *userController {
	return &userController{
		Service: user,
	}
}
func (c *userController) Register(w http.ResponseWriter, r *http.Request) {
	var u model.UserRegister
	c.Service.Register(&u)
}
