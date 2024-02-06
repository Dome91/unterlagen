package repository

import (
	"encoding/json"
	"errors"
	"github.com/spf13/afero"
	"io"
	"maps"
	"os"
	"sync"
	"unterlagen/pkg/domain"
)

const (
	DOCUMENT int = iota
	FOLDER
	USER
)

var files = map[int]string{
	DOCUMENT: "documents.json",
	FOLDER:   "folders.json",
	USER:     "users.json",
}

var ErrEmptyFile = errors.New("unexpected end of JSON input")

type FileRepository[T domain.Entity] struct {
	fs       afero.Fs
	filename string
	store    map[string]T
	mutex    sync.Mutex
}

func (fr *FileRepository[T]) Save(entity T) error {
	fr.store[entity.ID()] = entity
	err := fr.persist()
	if err != nil {
		delete(fr.store, entity.ID())
		return err
	}
	return nil
}

func (fr *FileRepository[T]) FindById(id string) (T, error) {
	entity, ok := fr.store[id]
	if !ok {
		return entity, domain.ErrDocumentNotFound
	}

	return entity, nil
}

func (fr *FileRepository[T]) FindAll() ([]T, error) {
	var entities []T
	for _, entity := range fr.store {
		entities = append(entities, entity)
	}

	return entities, nil
}

func (fr *FileRepository[T]) ExistsById(id string) (bool, error) {
	_, ok := fr.store[id]
	return ok, nil
}

func (fr *FileRepository[T]) DeleteAll() error {
	tmp := make(map[string]T, len(fr.store))
	maps.Copy(tmp, fr.store)
	clear(fr.store)

	err := fr.persist()
	if err != nil {
		maps.Copy(fr.store, tmp)
		return err
	}

	return nil
}

func (fr *FileRepository[T]) persist() error {
	var values []T
	for _, value := range fr.store {
		values = append(values, value)
	}

	data, err := json.Marshal(values)
	if err != nil {
		return err
	}

	fr.mutex.Lock()
	defer fr.mutex.Unlock()

	file, err := fr.fs.OpenFile(fr.filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, string(data))
	return err
}

func (fr *FileRepository[T]) load() {
	file, err := fr.fs.OpenFile(fr.filename, os.O_APPEND|os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	info, err := fr.fs.Stat(fr.filename)
	if err != nil {
		return
	}

	var entities []T
	if info.Size() > 0 {
		err = json.Unmarshal(data, &entities)
		if err != nil && !errors.Is(err, ErrEmptyFile) {
			panic(err)
		}

		for _, value := range entities {
			fr.store[value.ID()] = value
		}
	}
}
