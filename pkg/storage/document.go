package storage

import (
	"github.com/spf13/afero"
	"io"
	"io/fs"
	"os"
	"path"
)

var documentsFolder = "documents"

type FilesystemDocumentStorage struct {
	fs afero.Fs
}

type FileDocumentStorageOptions struct {
	FS afero.Fs
}

func NewDocumentStorage(options ...FileDocumentStorageOptions) *FilesystemDocumentStorage {
	initialize()
	var _options FileDocumentStorageOptions
	if len(options) == 0 {
		_options = FileDocumentStorageOptions{FS: afero.NewOsFs()}
	} else {
		_options = options[0]
	}

	err := _options.FS.MkdirAll(documentsFolder, 0755)
	if err != nil {
		panic(err)
	}

	return &FilesystemDocumentStorage{
		fs: _options.FS,
	}
}

func (fds *FilesystemDocumentStorage) Retrieve(id string, consumer func(reader io.Reader) error) error {
	file, err := fds.fs.OpenFile(path.Join(documentsFolder, id), os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return consumer(file)
}

func (fds *FilesystemDocumentStorage) Store(reader io.Reader, ID string) error {
	file, err := fds.fs.OpenFile(path.Join(documentsFolder, ID), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := file.Write(buf[:n]); err != nil {
			return err
		}
	}

	return nil
}

func (fds *FilesystemDocumentStorage) Size(id string) (int, error) {
	fpath := path.Join(documentsFolder, id)
	info, err := fds.fs.Stat(fpath)
	if err != nil {
		return 0, err
	}
	return int(info.Size()), nil
}

func (fds *FilesystemDocumentStorage) Clear() error {
	return afero.Walk(fds.fs, documentsFolder, func(path string, info fs.FileInfo, _ error) error {
		if !info.IsDir() {
			return fds.fs.Remove(path)
		}
		return nil
	})
}
