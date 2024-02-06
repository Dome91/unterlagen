package main

import (
	"unterlagen/pkg/domain"
	"unterlagen/pkg/eventing"
	"unterlagen/pkg/repository"
	"unterlagen/pkg/storage"
	"unterlagen/pkg/web"
)

func main() {
	eventBus := eventing.NewEventBus()
	folders := domain.NewFolders(repository.NewFolderRepository(), eventBus)
	documents := domain.NewDocuments(repository.NewDocumentRepository(), storage.NewDocumentStorage())
	users := domain.NewUsers(repository.NewUserRepository(), eventBus)

	err := web.StartServer(documents, folders, users)
	if err != nil {
		panic(err)
	}
}
