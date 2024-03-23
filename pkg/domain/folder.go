package domain

import (
	"context"
	"errors"
)

var ErrParentFolderDoesNotExist = errors.New("parent folder does not exist")
var ErrFolderNotFound = errors.New("folder not found")

type Folder struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ParentId string `json:"parentId"`
	UserId   string `json:"userId"`
}

func (f Folder) ID() string {
	return f.Id
}

func (f Folder) IsRoot() bool {
	return f.ParentId == ""
}

type FolderRepository interface {
	Save(folder Folder) error
	FindById(ID string) (Folder, error)
	FindAllByParentId(parentId string) ([]Folder, error)
	FindByUserIdAndParentIdEmpty(userId string) (Folder, error)
	DeleteAll() error
}

type Folders struct {
	repository FolderRepository
}

func NewFolders(repository FolderRepository, eventBus EventBus) *Folders {
	folders := &Folders{repository: repository}
	err := eventBus.Subscribe(CreatedUserEvent{}.Topic(), folders)
	if err != nil {
		panic(err)
	}

	return folders
}

func (f *Folders) Handle(event Event) error {
	createdUserEvent, ok := event.(CreatedUserEvent)
	if !ok {
		return nil
	}

	return f.createRoot(createdUserEvent.UserId)
}

func (f *Folders) createRoot(userId string) error {
	folder := Folder{
		Id:       generateID(),
		Name:     "Home",
		ParentId: "",
		UserId:   userId,
	}

	return f.repository.Save(folder)
}

func (f *Folders) CreateChild(ctx context.Context, name string, parentId string) error {
	userId := CurrentUser(ctx)
	parent, err := f.repository.FindById(parentId)
	if err != nil {
		return ErrParentFolderDoesNotExist
	}
	if parent.UserId != userId {
		return ErrUserUnauthorized
	}

	folder := Folder{
		Id:       generateID(),
		Name:     name,
		ParentId: parent.ID(),
		UserId:   userId,
	}

	return f.repository.Save(folder)
}

func (f *Folders) Get(ctx context.Context, id string) (Folder, error) {
	folder, err := f.repository.FindById(id)
	if err != nil {
		return Folder{}, err
	}

	if folder.UserId != CurrentUser(ctx) {
		return Folder{}, ErrUserUnauthorized
	}

	return folder, nil
}

func (f *Folders) GetRoot(ctx context.Context) (Folder, error) {
	userId := CurrentUser(ctx)
	return f.repository.FindByUserIdAndParentIdEmpty(userId)
}

func (f *Folders) GetChildren(ctx context.Context, parentId string) ([]Folder, error) {
	userId := CurrentUser(ctx)
	folder, err := f.repository.FindById(parentId)
	if err != nil {
		return nil, err
	}

	if folder.UserId != userId {
		return nil, ErrUserUnauthorized
	}

	return f.repository.FindAllByParentId(parentId)
}
