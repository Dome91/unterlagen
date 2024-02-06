package repository

import (
	"path"
	"unterlagen/pkg/config"
)

func init() {
	if config.Development {
		for index, file := range files {
			files[index] = path.Join(".ws", file)
		}
	}
}
