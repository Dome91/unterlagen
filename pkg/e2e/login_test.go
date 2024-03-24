package e2e

import (
	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var baseUrl = os.Getenv("BASE_URL")

func TestUserSuccessfullyLogsIn(t *testing.T) {
	if baseUrl == "" {
		baseUrl = "http://localhost:8080"
	}

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
	page, err := context.NewPage()
	_, err = page.Goto("/")
	assert.Nil(t, err)

	assert.Nil(t, page.GetByLabel("username").Fill("admin"))
	assert.Nil(t, page.GetByLabel("password").Fill("admin"))
	assert.Nil(t, page.GetByText("Login").Click())
	assert.Nil(t, page.WaitForURL("**/folders"))
}
