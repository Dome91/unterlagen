package repository

import (
	"path"
	"sync"
	"unterlagen/pkg/config"
)

var initialize = sync.OnceFunc(func() {
	for index, file := range files {
		if config.Get().Development {
			files[index] = path.Join(".ws", file)
		} else {
			files[index] = path.Join(config.Get().DataDirectory, file)
		}
	}
})
