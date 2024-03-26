package auth

import (
	"context"
	"encoding/gob"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"net/http"
	"slices"
	"unterlagen/pkg/config"
	"unterlagen/pkg/domain"
)

const (
	sessionName = "unterlagen"
)

var (
	whitelist = []string{"/", "/register", "/unterlagen.css", "/unterlagen.js"}
	store     sessions.Store
)

func ConfigureSession(router chi.Router) {
	gob.Register(domain.UserRole(""))
	store = sessions.NewCookieStore(config.Get().CookieSecret)
	router.Use(validateSession)
}

func validateSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		path := request.URL.Path
		if slices.Contains(whitelist, path) {
			next.ServeHTTP(writer, request)
			return
		}

		session, err := store.Get(request, sessionName)
		if err != nil || session.IsNew {
			http.Redirect(writer, request, "/", http.StatusMovedPermanently)
			return
		}

		username := session.Values[domain.ContextKeyUserId]
		role := session.Values[domain.ContextKeyUserRole]
		ctx := context.WithValue(request.Context(), domain.ContextKeyUserId, username)
		ctx = context.WithValue(ctx, domain.ContextKeyUserRole, role)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func CreateSession(writer http.ResponseWriter, request *http.Request, user domain.User) error {
	session, err := store.Get(request, sessionName)
	if err != nil {
		return err
	}

	session.Values[domain.ContextKeyUserId] = user.ID()
	session.Values[domain.ContextKeyUserRole] = user.Role
	return session.Save(request, writer)
}

func DeleteSession(writer http.ResponseWriter, request *http.Request) error {
	session, err := store.Get(request, sessionName)
	if err != nil {
		return err
	}

	session.Options.MaxAge = -1
	return session.Save(request, writer)
}
