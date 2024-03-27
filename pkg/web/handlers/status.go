package handlers

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

func NotFound(executor TemplateExecutor) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := executor.ExecuteTemplate(writer, "notFound.gohtml", nil)
		if err != nil {
			log.Err(err).Msg("failed to template notFound.gohtml")
		}
	}
}
