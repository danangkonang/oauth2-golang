package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/danangkonang/oauth2-golang/helper"
	"github.com/danangkonang/oauth2-golang/model"
	"github.com/danangkonang/oauth2-golang/service"
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
	var user *model.UserRegister
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helper.MakeRespon(w, 400, err.Error(), nil)
		return
	}
	defer r.Body.Close()
	if err := c.Service.IsUserAlrady(user.UserName); err != nil {
		helper.MakeRespon(w, 500, err.Error(), nil)
		return
	}
	hashPass := helper.HashPassword(user.Password)
	user.Password = hashPass
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	if err := c.Service.Register(user); err != nil {
		helper.MakeRespon(w, 500, err.Error(), nil)
		return
	}
	helper.MakeRespon(w, 200, "success", nil)
}
