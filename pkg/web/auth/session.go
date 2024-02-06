package auth

import (
	"context"
	"crypto/rand"
	"encoding/gob"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog/log"
	"math/big"
	"net/http"
	"slices"
	"unterlagen/pkg/config"
	"unterlagen/pkg/domain"
)

const (
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChars    = "0123456789"
	passwordLength = 10
	sessionName    = "unterlagen"
)

var (
	whitelist = []string{"/", "/register", "/unterlagen.css", "/unterlagen.js"}
	store     sessions.Store
)

func ConfigureSession(router chi.Router, users *domain.Users) {
	gob.Register(domain.UserRole(""))
	store = sessions.NewCookieStore(config.CookieSecret)

	adminExists, err := users.ExistsByRole(domain.UserRoleAdmin)
	if err != nil {
		panic(err)
	}

	if !adminExists {
		var adminPassword string
		if config.Development {
			adminPassword = "admin"
		} else {
			adminPassword = generatePassword()
		}

		err := users.Create("admin", adminPassword, domain.UserRoleAdmin)
		if err != nil {
			panic(err)
		}
		log.Info().Str("password", adminPassword).Str("username", "admin").Msg("Generated credentials")
	}

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

func generatePassword() string {
	var password string
	charset := lowercaseChars + uppercaseChars + numberChars

	for i := 0; i < passwordLength; i++ {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password += string(charset[randomIndex.Int64()])
	}

	return password
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
