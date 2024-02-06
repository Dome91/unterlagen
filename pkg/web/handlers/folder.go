package handlers

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"slices"
	"unterlagen/pkg/domain"
	"unterlagen/pkg/web/httpx"
)

func CreateFolder(folders *domain.Folders) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		name := request.FormValue("name")
		parentId := request.FormValue("parentId")

		err := folders.CreateChild(request.Context(), name, parentId)
		if err != nil {
			log.Err(err).Msg("failed to create folder")
			return
		}

		httpx.Redirect(writer, fmt.Sprintf("/folders?folderId=%s", parentId), http.StatusCreated)
	}
}

func GetFolder(documents *domain.Documents, folders *domain.Folders, executor TemplateExecutor) http.HandlerFunc {
	type Breadcrumb struct {
		Id   string
		Name string
	}

	buildBreadcrumbs := func(ctx context.Context, folderId string) []Breadcrumb {
		var breadcrumbs []Breadcrumb
		for {
			parent, err := folders.Get(ctx, folderId)
			if err != nil {
				log.Err(err).Msg("failed to get parent")
				break
			}

			breadcrumbs = append(breadcrumbs, Breadcrumb{
				Id:   parent.Id,
				Name: parent.Name,
			})

			if parent.ParentId == "" {
				break
			}
			folderId = parent.ParentId
		}

		slices.Reverse(breadcrumbs)
		return breadcrumbs
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		folderIds := request.URL.Query()["folderId"]
		var folderId string
		if len(folderIds) == 0 {
			rootFolderOfUser, err := folders.GetRoot(ctx)
			if err != nil {
				return
			}
			url := fmt.Sprintf("/folders?folderId=%s", rootFolderOfUser.Id)
			http.Redirect(writer, request, url, http.StatusPermanentRedirect)
			return
		}

		folderId = folderIds[0]
		documentsInFolder, err := documents.GetInFolder(ctx, folderId)
		if err != nil {
			log.Err(err).Msg("failed to get documents")
			return
		}

		foldersInFolder, err := folders.GetChildren(ctx, folderId)
		if err != nil {
			log.Err(err).Msg("failed to get folders")
			return
		}

		err = executor.ExecuteTemplate(writer, "folders.gohtml", map[string]any{
			"folderId":    folderId,
			"breadcrumbs": buildBreadcrumbs(ctx, folderId),
			"documents":   documentsInFolder,
			"folders":     foldersInFolder,
		})
		if err != nil {
			log.Err(err).Msg("failed to template folder.gohtml")
		}
	}
}
