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

func withServer(t *testing.T, block func(page playwright.Page, pwAssert playwright.PlaywrightAssertions)) {
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

	err := playwright.Install()
	require.Nil(t, err)
	pw, err := playwright.Run()
	defer pw.Stop()

	assert.Nil(t, err)
	browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	assert.Nil(t, err)
	baseUrl := "http://localhost:8080"
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		BaseURL: &baseUrl,
	})
	assert.Nil(t, err)
	page, err := context.NewPage()
	assert.Nil(t, err)
	_, err = page.Goto("/")
	assert.Nil(t, err)

	assert.Nil(t, page.GetByLabel("username").Fill("admin"))
	assert.Nil(t, page.GetByLabel("password").Fill("e2e"))
	assert.Nil(t, page.GetByText("Login").Click())
	assert.Nil(t, page.WaitForURL("**/folders"))

	pwAssert := playwright.NewPlaywrightAssertions()
	block(page, pwAssert)
}
