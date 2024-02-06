package repository

import (
	"github.com/spf13/afero"
	"sync"
	"unterlagen/pkg/domain"
)

type FileUserRepositoryOptions struct {
	FS afero.Fs
}

type FileUserRepository struct {
	FileRepository[domain.User]
}

func (r *FileUserRepository) ExistsByRole(role domain.UserRole) (bool, error) {
	for _, user := range r.store {
		if user.Role == role {
			return true, nil
		}
	}

	return false, nil
}

func NewUserRepository(options ...FileUserRepositoryOptions) *FileUserRepository {
	var _options FileUserRepositoryOptions
	if len(options) == 0 {
		_options = FileUserRepositoryOptions{FS: afero.NewOsFs()}
	} else {
		_options = options[0]
	}

	repository := &FileUserRepository{
		FileRepository[domain.User]{
			fs:       _options.FS,
			filename: files[USER],
			store:    make(map[string]domain.User),
			mutex:    sync.Mutex{},
		},
	}

	repository.load()
	return repository
}
