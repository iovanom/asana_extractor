package storage

import (
	"fmt"
	"io"
	"os"
	"path"
)

type LocalStorage struct {
	dir string
}

func NewLocalStorage(dir string) (*LocalStorage, error) {
	stat, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("directory %s does not exist", dir)
		}
		return nil, err
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("this is not a directory %s", dir)
	}
	return &LocalStorage{dir}, nil
}

func (s *LocalStorage) filePath(filename string) string {
	return path.Join(s.dir, filename)
}

func (s *LocalStorage) SaveFile(filename string, body io.Reader) error {
	f, err := os.OpenFile(s.filePath(filename), os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, body)
	return err
}
