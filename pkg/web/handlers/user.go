package handlers

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"unterlagen/pkg/domain"
	"unterlagen/pkg/web/auth"
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
			http.Redirect(writer, request, "/", http.StatusMovedPermanently)
			return
		}

		if !user.IsValid(password) {
			log.Err(err).Msg("password is invalid")
			http.Redirect(writer, request, "/", http.StatusMovedPermanently)
			return
		}

		err = auth.CreateSession(writer, request, user)
		if err != nil {
			log.Err(err).Msg("session could not be created")
			http.Redirect(writer, request, "/", http.StatusMovedPermanently)
			return
		}

		http.Redirect(writer, request, "/folders", http.StatusMovedPermanently)
	}
}

func LogoutUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := auth.DeleteSession(writer, request)
		if err != nil {
			log.Err(err).Msg("session could not be deleted")
		}
		writer.Header().Set("HX-Redirect", "/")
		//http.Redirect(writer, request, "/", http.StatusMovedPermanently)
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
