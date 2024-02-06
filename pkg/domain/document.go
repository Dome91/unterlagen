package domain

import (
	"context"
	"errors"
	"io"
	"time"
)

var ErrDocumentNotFound = errors.New("document not found")
var ErrDocumentFileNotFound = errors.New("document file not found")

type Document struct {
	Id        string    `json:"Id"`
	Filename  string    `json:"filename"`
	Filesize  int       `json:"filesize"`
	CreatedAt time.Time `json:"createdAt"`
	FolderId  string    `json:"folderId"`
	UserId    string    `json:"userId"`
}

func (d Document) ID() string {
	return d.Id
}

type DocumentRepository interface {
	Save(document Document) error
	FindById(id string) (Document, error)
	FindAll() ([]Document, error)
	FindAllByFolderIdAndUserId(folderId string, userId string) ([]Document, error)
	DeleteAll() error
}

type DocumentStorage interface {
	Store(reader io.Reader, ID string) error
	Retrieve(id string, consumer func(reader io.Reader) error) error
	Size(id string) (int, error)
	Clear() error
}

type Documents struct {
	repository DocumentRepository
	storage    DocumentStorage
}

func NewDocuments(repository DocumentRepository, storage DocumentStorage) *Documents {
	return &Documents{repository: repository, storage: storage}
}

func (d *Documents) Upload(ctx context.Context, reader io.Reader, filename string, folderId string) error {
	id := generateID()
	err := d.storage.Store(reader, id)
	if err != nil {
		return err
	}

	filesize, err := d.storage.Size(id)
	if err != nil {
		return err
	}

	document := Document{
		Id:        id,
		Filename:  filename,
		Filesize:  filesize,
		CreatedAt: time.Now(),
		FolderId:  folderId,
		UserId:    CurrentUser(ctx),
	}

	return d.repository.Save(document)
}

func (d *Documents) GetInFolder(ctx context.Context, folderId string) ([]Document, error) {
	return d.repository.FindAllByFolderIdAndUserId(folderId, CurrentUser(ctx))
}

func (d *Documents) Get(ctx context.Context, id string) (Document, error) {
	document, err := d.repository.FindById(id)
	if err != nil {
		return Document{}, err
	}

	if document.UserId == CurrentUser(ctx) {
		return Document{}, ErrUserUnauthorized
	}

	return document, nil
}

func (d *Documents) Download(ctx context.Context, id string, consumer func(r io.Reader) error) error {
	document, err := d.repository.FindById(id)
	if err != nil {
		return err
	}

	if document.UserId != CurrentUser(ctx) {
		return ErrUserUnauthorized
	}

	return d.storage.Retrieve(id, consumer)
}
