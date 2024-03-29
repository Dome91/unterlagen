package storage

import (
	"path"
	"sync"
	"unterlagen/pkg/config"
)

var initialize = sync.OnceFunc(func() {
	if config.Get().Development {
		documentsFolder = path.Join(".ws", config.Get().DataDirectory, documentsFolder)
	} else {
		documentsFolder = path.Join(config.Get().DataDirectory, documentsFolder)
	}
})
