package repository_test

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
	"unterlagen/pkg/domain"
	"unterlagen/pkg/repository"
)

func TestDocumentRepository(t *testing.T) {
	documentRepository := repository.NewDocumentRepository(repository.FileDocumentRepositoryOptions{
		FS: afero.NewMemMapFs(),
	})

	cleanup := func() {
		assert.Nil(t, documentRepository.DeleteAll())
	}

	t.Run("stores and retrieves document", func(t *testing.T) {
		t.Cleanup(cleanup)
		document := domain.GenerateTestDocument()

		err := documentRepository.Save(document)
		assert.Nil(t, err)

		foundDocument, err := documentRepository.FindById(document.Id)
		assert.Nil(t, err)
		assert.Equal(t, document, foundDocument)
	})

	t.Run("returns all documents in folder", func(t *testing.T) {
		t.Cleanup(cleanup)
		document1 := domain.GenerateTestDocument()
		document2 := domain.GenerateTestDocument()
		document3 := domain.GenerateTestDocument()
		document3.FolderId = "otherId"

		assert.Nil(t, documentRepository.Save(document1))
		assert.Nil(t, documentRepository.Save(document2))
		assert.Nil(t, documentRepository.Save(document3))

		documents, err := documentRepository.FindAllByFolderIdAndUserId(document1.FolderId, document1.UserId)
		assert.Nil(t, err)
		assert.Len(t, documents, 2)
		assert.Contains(t, documents, document1)
		assert.Contains(t, documents, document2)
	})
}
