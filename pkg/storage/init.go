package storage

import (
	"path"
	"sync"
	"unterlagen/pkg/config"
)

var initialize = sync.OnceFunc(func() {
	if config.Get().Development {
		documentsFolder = path.Join(".ws", documentsFolder)
	} else {
		documentsFolder = path.Join(config.Get().DataDirectory, documentsFolder)
	}
})
