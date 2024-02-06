package domain_test

import (
	"context"
	"errors"
	"unterlagen/pkg/domain"
)

var ErrNotFound = errors.New("entity not found")

type InMemoryRepository[T domain.Entity] struct {
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

func GenerateTestUserContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, domain.ContextKeyUserId, "userId")
	ctx = context.WithValue(ctx, domain.ContextKeyUserRole, domain.UserRoleUser)
	return ctx
}
