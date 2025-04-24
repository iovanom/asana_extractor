package storage

import (
	"io"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func initTestStorage(t *testing.T) (*LocalStorage, func()) {
	t.Helper()

	dir, err := os.MkdirTemp("", "*")
	if err != nil {
		t.Fatalf("error on create tmp dir %s", err)
	}

	s, err := NewLocalStorage(dir)
	if err != nil {
		t.Fatalf("error on initiate local storage %s", err)
	}

	return s, func() { os.RemoveAll(dir) }
}

func TestSaveFile(t *testing.T) {
	t.Parallel()

	t.Run("save file should work", func(t *testing.T) {
		s, tearDown := initTestStorage(t)
		defer tearDown()

		content := `{ "email": "test@test.te" }`
		fileName := "test_file.json"

		err := s.SaveFile(fileName, strings.NewReader(content))
		require.NoError(t, err)

		filePath := path.Join(s.dir, fileName)
		f, err := os.Open(filePath)
		require.NoError(t, err)
		fileContent, err := io.ReadAll(f)
		require.NoError(t, err)
		require.Equal(t, content, string(fileContent))
	})
}
