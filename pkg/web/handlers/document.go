package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"
	"unterlagen/pkg/domain"
)

func UploadDocument(documents *domain.Documents) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		folderId := request.FormValue("folderId")
		if folderId == "" {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		file, info, err := request.FormFile("document")
		if err != nil {
			log.Err(err).Msg("failed to resolve form file")
			return
		}
		defer file.Close()
		err = documents.Upload(request.Context(), file, info.Filename, folderId)
		if err != nil {
			log.Err(err).Msg("failed to upload document")
			return
		}

		writer.WriteHeader(http.StatusCreated)
	}
}

func DownloadDocument(documents *domain.Documents) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id := chi.URLParam(request, "id")
		document, err := documents.Get(request.Context(), id)
		if err != nil {
			log.Err(err).Msg("")
			return
		}

		writer.Header().Set("Content-Type", "application/octet-stream")
		writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", document.Filename))
		writer.Header().Set("Content-Length", strconv.Itoa(document.Filesize))
		consumer := func(r io.Reader) error {
			_, err := io.Copy(writer, r)
			return err
		}

		err = documents.Download(request.Context(), id, consumer)
		if err != nil {
			log.Err(err)
		}
	}
}
