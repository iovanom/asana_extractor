package asana

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockClient struct {
	resp *http.Response
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	return m.resp, nil
}

func initTestClient(t *testing.T, statusCode int, body string) *Client {
	t.Helper()
	resp := &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
	return &Client{httpClient: &mockClient{resp}}
}

func TestUsers(t *testing.T) {
	t.Parallel()

	t.Run("get empty list of users", func(t *testing.T) {
		c := initTestClient(t, http.StatusOK, `{ "data": [] }`)

		users, err := c.Users()
		require.NoError(t, err)
		require.Empty(t, users)
	})
}

func TestProjects(t *testing.T) {
	t.Parallel()

	t.Run("get empty list of projects", func(t *testing.T) {
		c := initTestClient(t, http.StatusOK, `{ "data": [] }`)

		projects, err := c.Projects()
		require.NoError(t, err)
		require.Empty(t, projects)
	})
}
