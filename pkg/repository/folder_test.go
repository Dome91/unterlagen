package repository_test

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
	"unterlagen/pkg/domain"
	"unterlagen/pkg/repository"
)

func TestFolderRepository(t *testing.T) {
	folderRepository := repository.NewFolderRepository(repository.FileFolderRepositoryOptions{
		FS: afero.NewMemMapFs(),
	})

	cleanup := func() {
		assert.Nil(t, folderRepository.DeleteAll())
	}

	t.Run("stores and retrieves folder", func(t *testing.T) {
		t.Cleanup(cleanup)
		folder := domain.GenerateTestFolder()

		err := folderRepository.Save(folder)
		assert.Nil(t, err)

		foundDocument, err := folderRepository.FindById(folder.Id)
		assert.Nil(t, err)
		assert.Equal(t, folder, foundDocument)
	})

	t.Run("returns all children folder in folder", func(t *testing.T) {
		t.Cleanup(cleanup)
		folder1 := domain.GenerateTestFolder()
		folder2 := domain.GenerateTestFolder()
		folder3 := domain.GenerateTestFolder()
		folder3.ParentId = "otherId"

		assert.Nil(t, folderRepository.Save(folder1))
		assert.Nil(t, folderRepository.Save(folder2))
		assert.Nil(t, folderRepository.Save(folder3))

		folders, err := folderRepository.FindAllByParentId(folder1.ParentId)
		assert.Nil(t, err)
		assert.Len(t, folders, 2)
		assert.Contains(t, folders, folder1)
		assert.Contains(t, folders, folder2)
	})
}
