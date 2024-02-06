package storage

import (
	"path"
	"unterlagen/pkg/config"
)

func init() {
	if config.Development {
		documentsFolder = path.Join(".ws", documentsFolder)
	}
}
