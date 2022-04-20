package middleware

import (
	"net/http"

	"github.com/danangkonang/oauth2-golang/helper"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
)

type oauthMiddleware struct {
	server *server.Server
}

func NewOauthMiddleware(m *manage.Manager) *oauthMiddleware {
	return &oauthMiddleware{
		server: server.NewServer(server.NewConfig(), m),
	}
}

func (o *oauthMiddleware) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := o.server.ValidationBearerToken(r)
		if err != nil {
			helper.MakeRespon(w, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}
		next(w, r)
	}
}
