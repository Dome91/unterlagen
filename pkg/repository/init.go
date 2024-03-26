package repository

import (
	"path"
	"sync"
	"unterlagen/pkg/config"
)

var initialize = sync.OnceFunc(func() {
	if config.Get().Development {
		for index, file := range files {
			files[index] = path.Join(".ws", file)
		}
	}
})
