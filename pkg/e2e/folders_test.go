package e2e

import (
	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestCreateFolderSuccessfully(t *testing.T) {
	withServer(t, func(page playwright.Page, pwAssert playwright.PlaywrightAssertions) {
		assert.Nil(t, page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Create folder"}).Click())
		assert.Nil(t, page.GetByRole("textbox", playwright.PageGetByRoleOptions{Name: "Name of folder"}).Fill("New Folder"))
		assert.Nil(t, page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Create", Exact: playwright.Bool(true)}).Click())
		assert.Nil(t, pwAssert.Locator(page.GetByText("New Folder")).ToBeVisible())
	})
}

func TestUploadAndDownloadDocumentSuccessfully(t *testing.T) {
	withServer(t, func(page playwright.Page, pwAssert playwright.PlaywrightAssertions) {
		assert.Nil(t, page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Upload document"}).Click())
		chooser, err := page.ExpectFileChooser(func() error {
			err := page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Choose Document"}).Click()
			return err
		})
		assert.Nil(t, err)

		wd, err := os.Getwd()
		assert.Nil(t, err)
		assert.Nil(t, chooser.SetFiles(path.Join(wd, "files", "dummy.pdf")))

		assert.Nil(t, page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Upload", Exact: playwright.Bool(true)}).Click())
		uploadedFile := page.GetByText("dummy")
		assert.Nil(t, pwAssert.Locator(uploadedFile).ToBeVisible())
		download, err := page.ExpectDownload(func() error {
			return uploadedFile.Click()
		})
		assert.Nil(t, err)
		assert.Nil(t, download.Failure())
	})
}
