package e2e

import (
	"github.com/playwright-community/playwright-go"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"unterlagen/pkg/domain"
	"unterlagen/pkg/eventing"
	"unterlagen/pkg/repository"
	"unterlagen/pkg/storage"
	"unterlagen/pkg/web"
)

func withServer(t *testing.T, block func(context playwright.BrowserContext)) {
	viper.Set("e2e", true)
	fs := afero.NewMemMapFs()
	eventBus := eventing.NewEventBus()
	documentStorage := storage.NewDocumentStorage(storage.FileDocumentStorageOptions{FS: fs})

	folderRepository := repository.NewFolderRepository(repository.FileFolderRepositoryOptions{FS: fs})
	documentRepository := repository.NewDocumentRepository(repository.FileDocumentRepositoryOptions{FS: fs})
	userRepository := repository.NewUserRepository(repository.FileUserRepositoryOptions{FS: fs})

	folders := domain.NewFolders(folderRepository, eventBus)
	documents := domain.NewDocuments(documentRepository, documentStorage)
	users := domain.NewUsers(userRepository, eventBus)
	go web.StartServer(documents, folders, users)
	defer web.StopServer()

	baseUrl := "http://localhost:8080"

	err := playwright.Install()
	require.Nil(t, err)
	pw, err := playwright.Run()
	assert.Nil(t, err)
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	assert.Nil(t, err)
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		BaseURL: &baseUrl,
	})
	assert.Nil(t, err)
	block(context)
}
