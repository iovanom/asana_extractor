package extractor

import (
	"io"
	"testing"

	"github.com/iovanom/asana_extractor/internal/models"
	"github.com/stretchr/testify/require"
)

type mockObj struct {
	usersFunc    func() ([]*models.User, error)
	projectsFunc func() ([]*models.Project, error)
	saveFunc     func(string, io.Reader) error
}

func (m *mockObj) Users() ([]*models.User, error) {
	if m.usersFunc != nil {
		return m.usersFunc()
	}
	return nil, nil
}

func (m *mockObj) Projects() ([]*models.Project, error) {
	if m.projectsFunc != nil {
		return m.projectsFunc()
	}
	return nil, nil
}

func (m *mockObj) SaveFile(filename string, body io.Reader) error {
	if m.saveFunc != nil {
		return m.saveFunc(filename, body)
	}
	return nil
}

func initExtractor(t *testing.T, mock *mockObj) *Extractor {
	t.Helper()
	return NewExtractor(mock, mock)
}

func TestExtractData(t *testing.T) {
	t.Parallel()

	t.Run("simple test", func(t *testing.T) {
		e := initExtractor(t, &mockObj{})
		err := e.ExtractData()
		require.NoError(t, err)
	})
}
