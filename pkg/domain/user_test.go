package domain_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unterlagen/pkg/domain"
	"unterlagen/pkg/eventing"
)

type TestSubscriber func(event domain.Event) error

func (s TestSubscriber) Handle(event domain.Event) error {
	return s(event)
}

func TestUsers(t *testing.T) {
	eventBus := eventing.NewEventBus()
	repository := NewUserRepository()
	users := domain.NewUsers(repository, eventBus)
	cleanup := func() {
		assert.Nil(t, repository.DeleteAll())
	}

	t.Run("creates new user", func(t *testing.T) {
		t.Cleanup(cleanup)
		err := eventBus.Subscribe(domain.CreatedUserEvent{}.Topic(), TestSubscriber(func(event domain.Event) error {
			createdUserEvent, ok := event.(domain.CreatedUserEvent)
			assert.True(t, ok)
			assert.Equal(t, createdUserEvent.UserId, "username")
			return nil
		}))
		assert.Nil(t, err)

		err = users.Create("username", "password", domain.UserRoleUser)
		assert.Nil(t, err)

		foundUser, err := repository.FindById("username")
		assert.Equal(t, domain.UserRoleUser, foundUser.Role)
		assert.True(t, foundUser.IsValid("password"))
	})

	t.Run("does not create user if username already exists", func(t *testing.T) {
		t.Cleanup(cleanup)
		user := domain.GenerateTestUser()
		err := repository.Save(user)
		assert.Nil(t, err)

		err = users.Create(user.Username, "password", domain.UserRoleUser)
		assert.Equal(t, domain.ErrUserAlreadyExists, err)

	})
}

type InMemoryUserRepository struct {
	InMemoryRepository[domain.User]
}

func NewUserRepository() domain.UserRepository {
	repository := InMemoryRepository[domain.User]{
		Store: make(map[string]domain.User),
	}
	return &InMemoryUserRepository{
		repository,
	}
}

func (r *InMemoryUserRepository) ExistsById(id string) (bool, error) {
	_, ok := r.Store[id]
	return ok, nil
}

func (r *InMemoryUserRepository) ExistsByRole(role domain.UserRole) (bool, error) {
	for _, user := range r.Store {
		if user.Role == role {
			return true, nil
		}
	}

	return false, nil
}
