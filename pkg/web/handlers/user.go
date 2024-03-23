package handlers

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"unterlagen/pkg/domain"
	"unterlagen/pkg/web/auth"
	"unterlagen/pkg/web/httpx"
)

func ShowLogin(executor TemplateExecutor) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := executor.ExecuteTemplate(writer, "login.gohtml", nil)
		if err != nil {
			log.Err(err).Msg("failed to template login.gohtml")
			return
		}
	}
}

func LoginUser(users *domain.Users) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		username := request.FormValue("username")
		password := request.FormValue("password")

		user, err := users.Get(username)
		if err != nil {
			log.Err(err).Msg("user does not exist")
			httpx.Redirect(writer, "/", http.StatusMovedPermanently)
			return
		}

		if !user.IsValid(password) {
			log.Err(err).Msg("password is invalid")
			httpx.Redirect(writer, "/", http.StatusMovedPermanently)
			return
		}

		err = auth.CreateSession(writer, request, user)
		if err != nil {
			log.Err(err).Msg("session could not be created")
			httpx.Redirect(writer, "/", http.StatusMovedPermanently)
			return
		}

		httpx.Redirect(writer, "/folders", http.StatusOK)
	}
}

func LogoutUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := auth.DeleteSession(writer, request)
		if err != nil {
			log.Err(err).Msg("session could not be deleted")
		}
		httpx.Redirect(writer, "/", http.StatusMovedPermanently)
	}
}

func ShowRegister(executor TemplateExecutor) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := executor.ExecuteTemplate(writer, "register.gohtml", nil)
		if err != nil {
			log.Err(err).Msg("failed to template register.gohtml")
			return
		}
	}
}

func RegisterUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}
