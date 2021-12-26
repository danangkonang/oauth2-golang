package controller

import (
	"net/http"

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
	err := s.server.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
