package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/danangkonang/oauth2/config"
	"github.com/danangkonang/oauth2/helper"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-session/session"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.MustTokenStorage(helper.NewMysqlTokenStore(config.Connection()))
	manager.MapAccessGenerate(helper.NewAccessGenerate())

	clientStore := helper.NewClientStore()
	clientStore.Set("client_id", &models.Client{
		ID:     "client_id",
		Secret: "client_secret",
		Domain: "http://localhost:9000",
	})
	manager.MapClientStorage(clientStore)

	srv := server.NewServer(server.NewConfig(), manager)

	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/auth", authHandler)

	http.HandleFunc("/oauth/authorize", func(w http.ResponseWriter, r *http.Request) {
		store, err := session.Start(r.Context(), w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var form url.Values
		if v, ok := store.Get("ReturnUri"); ok {
			form = v.(url.Values)
		}
		r.Form = form
		store.Delete("ReturnUri")
		store.Save()
		err = srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		token, err := srv.ValidationBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := map[string]interface{}{
			// "expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
			"client_id": token.GetClientID(),
			"user_id":   token.GetUserID(),
		}
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(data)
	})

	log.Printf("Server is running at %d port.\n", 9096)
	log.Printf("Point your OAuth client Auth endpoint to %s:%d%s", "http://localhost", 9096, "/oauth/authorize")
	log.Printf("Point your OAuth client Token endpoint to %s:%d%s", "http://localhost", 9096, "/oauth/token")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 9096), nil))
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Method == "POST" {
		if r.Form == nil {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		store.Set("LoggedInUserID", r.Form.Get("username"))
		store.Save()

		w.Header().Set("Location", "/auth")
		w.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(w, r, "static/login.html")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, ok := store.Get("LoggedInUserID"); !ok {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(w, r, "static/auth.html")
}

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}
