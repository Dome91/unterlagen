package domain_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
	"time"
	"unterlagen/pkg/domain"
)

func TestDocumentService(t *testing.T) {
	repository := NewDocumentRepository()
	storage := NewDocumentStorage()
	documents := domain.NewDocuments(repository, storage)
	t.Cleanup(func() {
		require.Nil(t, repository.DeleteAll())
		require.Nil(t, storage.Clear())
	})
	ctx := GenerateTestUserContext()

	t.Run("uploads document successfully", func(t *testing.T) {
		err := documents.Upload(ctx, bytes.NewReader([]byte{1, 2, 3, 4}), "doc.pdf", "root")
		assert.Nil(t, err)

		documents, err := repository.FindAll()
		assert.Nil(t, err)
		assert.Len(t, documents, 1)

		document := documents[0]
		assert.Equal(t, "doc.pdf", document.Filename)
		assert.Equal(t, 4, document.Filesize)
		assert.WithinDuration(t, document.CreatedAt, time.Now(), time.Second)
	})

	t.Run("downloads document successfully", func(t *testing.T) {
		err := documents.Upload(ctx, bytes.NewReader([]byte{1, 2, 3, 4}), "doc.pdf", "root")
		assert.Nil(t, err)

		allDocuments, err := repository.FindAll()
		assert.Nil(t, err)

		err = documents.Download(ctx, allDocuments[0].ID(), func(r io.Reader) error {
			buf := make([]byte, 8)
			n, err := r.Read(buf)
			assert.Nil(t, err)
			assert.Equal(t, n, 4)
			return nil
		})
		assert.Nil(t, err)
	})
}

type InMemoryDocumentRepository struct {
	InMemoryRepository[domain.Document]
}

func NewDocumentRepository() *InMemoryDocumentRepository {
	repository := InMemoryRepository[domain.Document]{
		Store: make(map[string]domain.Document),
	}
	return &InMemoryDocumentRepository{repository}
}

func (r *InMemoryDocumentRepository) FindAllByFolderIdAndUserId(folderId string, userId string) ([]domain.Document, error) {
	var documents []domain.Document
	for _, document := range r.Store {
		if document.FolderId == folderId && document.UserId == userId {
			documents = append(documents, document)
		}
	}

	return documents, nil
}

type InMemoryDocumentStorage struct {
	Storage map[string][]byte
}

func NewDocumentStorage() *InMemoryDocumentStorage {
	storage := make(map[string][]byte)
	return &InMemoryDocumentStorage{Storage: storage}
}

func (s *InMemoryDocumentStorage) Store(reader io.Reader, ID string) error {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, reader)
	if err != nil {
		return err
	}

	s.Storage[ID] = buf.Bytes()
	return nil
}

func (s *InMemoryDocumentStorage) Retrieve(id string, consumer func(reader io.Reader) error) error {
	file, ok := s.Storage[id]
	if !ok {
		return domain.ErrDocumentNotFound
	}

	reader := bytes.NewReader(file)
	return consumer(reader)
}

func (s *InMemoryDocumentStorage) Size(id string) (int, error) {
	fileBytes, ok := s.Storage[id]
	if !ok {
		return 0, domain.ErrDocumentFileNotFound
	}

	return len(fileBytes), nil
}

func (s *InMemoryDocumentStorage) Clear() error {
	clear(s.Storage)
	return nil
}
