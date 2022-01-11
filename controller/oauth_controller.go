package controller

import (
	"net/http"
	"time"

	"github.com/danangkonang/oauth2/helper"
	"github.com/danangkonang/oauth2/service"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-session/session"
	"golang.org/x/crypto/bcrypt"
)

type oauthController struct {
	server *server.Server
	user   service.UserService
}

func NewOauthController(manager *manage.Manager, user service.UserService) *oauthController {
	return &oauthController{
		server: server.NewServer(server.NewConfig(), manager),
		user:   user,
	}
}

func (s *oauthController) Login(w http.ResponseWriter, r *http.Request) {
	s.server.SetUserAuthorizationHandler(userAuthorizeHandler)
	res, err := s.user.Login(r.FormValue("username"))
	if err != nil {
		helper.MakeRespon(w, 500, err.Error(), nil)
		return
	}
	err_pass := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(r.FormValue("password")))
	if err_pass != nil {
		helper.MakeRespon(w, 400, "invalid username or password", nil)
		return
	}
	s.server.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		return res.Id, nil
	})
	if err := s.server.HandleTokenRequest(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *oauthController) Secure(w http.ResponseWriter, r *http.Request) {
	token, err := s.server.ValidationBearerToken(r)
	if err != nil {
		helper.MakeRespon(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	ut := time.Now()
	data := map[string]interface{}{
		"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(ut).Seconds()),
		"client_id":  token.GetClientID(),
		"user_id":    token.GetUserID(),
	}
	helper.MakeRespon(w, 200, "success", data)
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return
	}
	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}
		store.Set("ReturnUri", r.Form)
		store.Save()
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	userID = uid.(string)
	store.Delete("LoggedInUserID")
	store.Save()
	return
}
