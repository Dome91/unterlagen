package domain

import (
	"context"
	"crypto/rand"
	"errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"unterlagen/pkg/config"
)

const (
	UserRoleUser       UserRole = "USER"
	UserRoleAdmin      UserRole = "ADMIN"
	ContextKeyUserId            = "userId"
	ContextKeyUserRole          = "role"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserUnauthorized  = errors.New("user unauthorized")
)

type CreatedUserEvent struct {
	UserId string
}

func (e CreatedUserEvent) Topic() string {
	return "created-user-event"
}

type UserRole string

type User struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Role     UserRole `json:"role"`
}

func (u User) ID() string {
	return u.Username
}

func (u User) IsValid(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

type UserContext interface {
	Current() (User, error)
}

type UserRepository interface {
	Save(user User) error
	FindById(id string) (User, error)
	ExistsById(id string) (bool, error)
	ExistsByRole(role UserRole) (bool, error)
	DeleteAll() error
}

type Users struct {
	repository UserRepository
	eventBus   EventBus
}

func NewUsers(repository UserRepository, eventBus EventBus) *Users {
	return &Users{repository: repository, eventBus: eventBus}
}

func (u *Users) Create(username string, password string, role UserRole) error {
	exists, err := u.repository.ExistsById(username)
	if err != nil || exists {
		return ErrUserAlreadyExists
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	user := User{
		Username: username,
		Password: hashedPassword,
		Role:     role,
	}

	err = u.repository.Save(user)
	if err != nil {
		return err
	}

	event := CreatedUserEvent{user.Username}
	return u.eventBus.Publish(event)
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (u *Users) ExistsByRole(role UserRole) (bool, error) {
	exists, err := u.repository.ExistsByRole(role)
	return exists, err
}

func (u *Users) Get(username string) (User, error) {
	return u.repository.FindById(username)
}

func (u *Users) CreateAdmin() error {
	generatePassword := func() string {
		lowercaseChars := "abcdefghijklmnopqrstuvwxyz"
		uppercaseChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numberChars := "0123456789"
		passwordLength := 10

		var password string
		charset := lowercaseChars + uppercaseChars + numberChars
		for i := 0; i < passwordLength; i++ {
			randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			password += string(charset[randomIndex.Int64()])
		}
		return password
	}

	adminExists, err := u.repository.ExistsByRole(UserRoleAdmin)
	if err != nil {
		return err
	}

	if !adminExists {
		var adminPassword string
		if config.Get().Development {
			adminPassword = "admin"
		} else if config.Get().E2E {
			adminPassword = "e2e"
		} else {
			adminPassword = generatePassword()
		}

		err := u.Create("admin", adminPassword, UserRoleAdmin)
		if err != nil {
			return err
		}

		log.Info().Str("password", adminPassword).Str("username", "admin").Msg("Generated credentials")
	}

	return nil
}

func CurrentUser(ctx context.Context) string {
	if ctx.Value(ContextKeyUserId) == nil {
		panic(ErrUserUnauthorized)
	}

	userId := ctx.Value(ContextKeyUserId).(string)
	if userId == "" {
		panic(ErrUserUnauthorized)
	}

	return userId
}
