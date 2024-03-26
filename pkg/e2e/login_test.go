package e2e

import (
	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserSuccessfullyLogsIn(t *testing.T) {
	withServer(t, func(context playwright.BrowserContext) {
		page, err := context.NewPage()
		_, err = page.Goto("/")
		assert.Nil(t, err)

		assert.Nil(t, page.GetByLabel("username").Fill("admin"))
		assert.Nil(t, page.GetByLabel("password").Fill("e2e"))
		assert.Nil(t, page.GetByText("Login").Click())
		assert.Nil(t, page.WaitForURL("**/folders"))
	})
}
