package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/danangkonang/oauth2/helper"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-session/session"
)

type oauthController struct {
	server *server.Server
}

func NewOauthController(manager *manage.Manager) *oauthController {
	return &oauthController{
		server: server.NewServer(server.NewConfig(), manager),
	}
}

func (s *oauthController) Login(w http.ResponseWriter, r *http.Request) {
	s.server.SetUserAuthorizationHandler(userAuthorizeHandler)
	// var pass string
	// 	var id int
	// config.Connection().QueryRow("SELECT id, password from users WHERE user_name = ?", r.FormValue("username")).Scan(&id, &pass)
	// err_pass := bcrypt.CompareHashAndPassword([]byte(pass), []byte(r.FormValue("password")))
	// if err_pass != nil {
	// 	helper.MakeRespon(w, 400, "invalid username or password", nil)
	// 	return
	// }
	// srv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
	// 	return strconv.Itoa(id), nil
	// })
	s.server.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		return "1", nil
	})
	err := s.server.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *oauthController) Secure(w http.ResponseWriter, r *http.Request) {
	token, err := s.server.ValidationBearerToken(r)
	if err != nil {
		fmt.Println(err.Error())
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
