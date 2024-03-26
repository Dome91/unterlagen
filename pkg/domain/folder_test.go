package domain_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"unterlagen/pkg/domain"
	"unterlagen/pkg/eventing"
)

func TestFolders(t *testing.T) {
	eventBus := eventing.NewEventBus()
	repository := NewFolderRepository()
	folders := domain.NewFolders(repository, eventBus)
	ctx := GenerateTestUserContext()
	cleanup := func() {
		require.Nil(t, repository.DeleteAll())
	}

	t.Run("creates root folder for new user", func(t *testing.T) {
		t.Cleanup(cleanup)
		event := domain.CreatedUserEvent{UserId: domain.GenerateTestUser().ID()}
		err := eventBus.Publish(event)
		assert.Nil(t, err)

		folder, err := repository.FindByUserIdAndParentIdEmpty(event.UserId)
		assert.Nil(t, err)
		assert.Equal(t, folder.Name, "Home")
		assert.Equal(t, folder.ParentId, "")
	})

	t.Run("creates child folder successfully", func(t *testing.T) {
		t.Cleanup(cleanup)
		folder := domain.GenerateTestFolder()
		err := repository.Save(folder)
		assert.Nil(t, err)

		err = folders.CreateChild(ctx, "child", folder.Id)
		assert.Nil(t, err)

		folders, err := repository.FindAllByParentId(folder.Id)
		assert.Nil(t, err)
		assert.Equal(t, folders[0].Name, "child")
	})

	t.Run("returns error if parent does not exist", func(t *testing.T) {
		t.Cleanup(cleanup)
		err := folders.CreateChild(ctx, "child", "unknown")
		assert.EqualError(t, err, domain.ErrParentFolderDoesNotExist.Error())
	})

	t.Run("returns children of folder", func(t *testing.T) {
		t.Cleanup(cleanup)
		folder := domain.GenerateTestFolder()
		child1 := domain.GenerateTestFolder()
		child2 := domain.GenerateTestFolder()
		child1.ParentId = folder.Id
		child2.ParentId = folder.Id
		assert.Nil(t, repository.Save(folder))
		assert.Nil(t, repository.Save(child1))
		assert.Nil(t, repository.Save(child2))

		children, err := folders.GetChildren(ctx, folder.Id)
		assert.Nil(t, err)
		assert.Contains(t, children, child1)
		assert.Contains(t, children, child2)
	})

}

type InMemoryFolderRepository struct {
	InMemoryRepository[domain.Folder]
}

func NewFolderRepository() domain.FolderRepository {
	repository := InMemoryRepository[domain.Folder]{
		Store: make(map[string]domain.Folder),
	}
	return &InMemoryFolderRepository{
		repository,
	}
}

func (r *InMemoryFolderRepository) FindAllByParentId(parentId string) ([]domain.Folder, error) {
	var folders []domain.Folder
	for _, folder := range r.Store {
		if folder.ParentId == parentId {
			folders = append(folders, folder)
		}
	}

	return folders, nil
}

func (r *InMemoryFolderRepository) FindByUserIdAndParentIdEmpty(userId string) (domain.Folder, error) {
	for _, folder := range r.Store {
		if folder.UserId == userId && folder.ParentId == "" {
			return folder, nil
		}
	}

	return domain.Folder{}, domain.ErrFolderNotFound
}
