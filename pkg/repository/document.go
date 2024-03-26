package repository

import (
	"github.com/spf13/afero"
	"sync"
	"unterlagen/pkg/domain"
)

type FileDocumentRepositoryOptions struct {
	FS afero.Fs
}

type FileDocumentRepository struct {
	FileRepository[domain.Document]
}

func NewDocumentRepository(options ...FileDocumentRepositoryOptions) *FileDocumentRepository {
	initialize()
	var _options FileDocumentRepositoryOptions
	if len(options) == 0 {
		_options = FileDocumentRepositoryOptions{FS: afero.NewOsFs()}
	} else {
		_options = options[0]
	}

	repository := &FileDocumentRepository{
		FileRepository[domain.Document]{
			fs:       _options.FS,
			filename: files[DOCUMENT],
			store:    make(map[string]domain.Document),
			mutex:    sync.Mutex{},
		},
	}

	repository.load()
	return repository
}

func (fdr *FileDocumentRepository) FindAllByFolderIdAndUserId(folderId string, userId string) ([]domain.Document, error) {
	var documents []domain.Document
	for _, document := range fdr.store {
		if document.FolderId == folderId && document.UserId == userId {
			documents = append(documents, document)
		}
	}

	return documents, nil
}
