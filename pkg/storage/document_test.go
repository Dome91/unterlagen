package storage_test

import (
	"bytes"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	"unterlagen/pkg/storage"
)

func TestDocumentStorage(t *testing.T) {
	documentStorage := storage.NewDocumentStorage(storage.FileDocumentStorageOptions{FS: afero.NewMemMapFs()})
	cleanup := func() {
		err := documentStorage.Clear()
		assert.Nil(t, err)
	}

	t.Run("stores and retrieves document", func(t *testing.T) {
		t.Cleanup(cleanup)
		reader := bytes.NewReader([]byte{1, 2, 3, 4})
		err := documentStorage.Store(reader, "id")
		assert.Nil(t, err)

		err = documentStorage.Retrieve("id", func(reader io.Reader) error {
			documentData, err := io.ReadAll(reader)
			if err != nil {
				return err
			}
			assert.Len(t, documentData, 4)
			return nil
		})
		assert.Nil(t, err)
	})

	t.Run("returns correct size of document", func(t *testing.T) {
		t.Cleanup(cleanup)
		reader := bytes.NewReader([]byte{1, 2, 3, 4})
		err := documentStorage.Store(reader, "id")
		assert.Nil(t, err)

		size, err := documentStorage.Size("id")
		assert.Nil(t, err)
		assert.Equal(t, 4, size)
	})
}
