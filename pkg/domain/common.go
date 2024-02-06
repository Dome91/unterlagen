package domain

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type Event interface {
	Topic() string
}

type EventHandler interface {
	Handle(event Event) error
}
type EventBus interface {
	Publish(event Event) error
	Subscribe(topic string, handler EventHandler) error
}

type Entity interface {
	ID() string
}

func generateID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return id.String()
}

func GenerateTestDocument() Document {
	createdAt := time.Date(2024, 8, 28, 21, 12, 0, 0, time.UTC)
	return Document{
		Id:        uuid.NewString(),
		Filename:  "important_document",
		Filesize:  1024,
		CreatedAt: createdAt,
		FolderId:  "userId-root",
		UserId:    "userId",
	}
}

func GenerateTestFolder() Folder {
	return Folder{
		Id:       uuid.NewString(),
		Name:     "Folder",
		ParentId: "",
		UserId:   "userId",
	}
}

func GenerateTestUser() User {
	return User{
		Username: uuid.NewString(),
		Password: "password",
		Role:     UserRoleUser,
	}
}

var ErrNotFound = errors.New("entity not found")

type InMemoryRepository[T Entity] struct {
	Store map[string]T
}

func (r *InMemoryRepository[T]) Save(entity T) error {
	r.Store[entity.ID()] = entity
	return nil
}

func (r *InMemoryRepository[T]) FindById(id string) (T, error) {
	entity, ok := r.Store[id]
	if !ok {
		return entity, ErrNotFound
	}

	return entity, nil
}

func (r *InMemoryRepository[T]) FindAll() ([]T, error) {
	var entities []T
	for _, entity := range r.Store {
		entities = append(entities, entity)
	}

	return entities, nil
}

func (r *InMemoryRepository[T]) DeleteAll() error {
	clear(r.Store)
	return nil
}
