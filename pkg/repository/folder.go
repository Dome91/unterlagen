package repository

import (
	"github.com/spf13/afero"
	"path"
	"sync"
	"unterlagen/pkg/domain"
)

type FileFolderRepositoryOptions struct {
	FS afero.Fs
}

type FileFolderRepository struct {
	FileRepository[domain.Folder]
}

func NewFolderRepository(options ...FileFolderRepositoryOptions) *FileFolderRepository {
	initialize()
	var _options FileFolderRepositoryOptions
	if len(options) == 0 {
		_options = FileFolderRepositoryOptions{FS: afero.NewOsFs()}
	} else {
		_options = options[0]
	}

	err := _options.FS.MkdirAll(path.Dir(files[FOLDER]), 0755)
	if err != nil {
		panic(err)
	}

	repository := &FileFolderRepository{
		FileRepository[domain.Folder]{
			fs:       _options.FS,
			filename: files[FOLDER],
			store:    make(map[string]domain.Folder),
			mutex:    sync.Mutex{},
		},
	}
	repository.load()
	return repository
}

func (ffr *FileFolderRepository) FindAllByParentId(parentId string) ([]domain.Folder, error) {
	var folders []domain.Folder

	for _, folder := range ffr.store {
		if folder.ParentId == parentId {
			folders = append(folders, folder)
		}
	}

	return folders, nil
}

func (ffr *FileFolderRepository) FindByUserIdAndParentIdEmpty(userId string) (domain.Folder, error) {
	for _, folder := range ffr.store {
		if folder.UserId == userId && folder.ParentId == "" {
			return folder, nil
		}
	}

	return domain.Folder{}, domain.ErrFolderNotFound
}
