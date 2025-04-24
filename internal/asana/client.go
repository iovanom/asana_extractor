package asana

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/iovanom/asana_extractor/internal/models"
)

const (
	defaultTimeout    = 1 * time.Second
	baseURL           = "https://app.asana.com/api/1.0/"
	maxRetries        = 25
	defaultRetryDelay = 60 // in seconds
)

var projectOptFields = []string{
	"name",
	"archived",
	"color",
	"completed",
	"completed_at",
	"completed_by",
}

var userOptFields = []string{
	"email",
	"name",
}

type bodyResponse[T any] struct {
	Data     []T `json:"data"`
	NextPage struct {
		Offset string `json:"offset"`
	} `json:"next_page"`
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient httpClient
	token      string
	workspace  string
}

func NewClient(token, workspace string) (*Client, error) {
	httpClient := &http.Client{
		Timeout: defaultTimeout,
	}
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}
	return &Client{httpClient, token, workspace}, nil
}

func (c *Client) prepareRequest(method, uri string, q *url.Values, body io.Reader) (*http.Request, error) {
	if q == nil {
		q = &url.Values{}
	}
	q.Add("workspace", c.workspace)
	url := baseURL + uri
	url = url + "?" + q.Encode()
	req, err := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", "Bearer "+c.token)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	// TODO: here we need to implement the retry logic and rate limit handling logic
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	retryCount := 0
	if resp.StatusCode == http.StatusTooManyRequests && retryCount < maxRetries {
		retryAfter, err := strconv.Atoi(resp.Header.Get("Retry-After"))
		if err != nil {
			slog.Debug(`invalid Retry-After header`, "retry-after", resp.Header.Get("Retry-After"))
			retryAfter = defaultRetryDelay
		}
		time.Sleep(time.Duration(retryAfter) * time.Second)
		return c.do(req)
	}
	return resp, err
}

func (c *Client) Users() ([]*models.User, error) {
	var users []*models.User

	var offset string
	for {
		query := url.Values{}
		query.Set("limit", "100")
		query.Set("opt_fields", strings.Join(userOptFields, ","))
		if offset != "" {
			query.Set("offset", offset)
		}
		req, err := c.prepareRequest("GET", "users", &query, nil)
		if err != nil {
			return nil, err
		}
		resp, err := c.do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			slog.Debug("wrong response on get users", "status", resp.Status, "body", b)
			return nil, fmt.Errorf("response status not ok on get users %s", resp.Status)
		}

		var body bodyResponse[*models.User]

		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			return nil, err
		}
		users = append(users, body.Data...)
		offset = body.NextPage.Offset
		if offset == "" {
			break
		}
	}
	return users, nil
}

func (c *Client) Projects() ([]*models.Project, error) {
	var projects []*models.Project

	var offset string
	for {
		query := url.Values{}
		query.Set("limit", "100")
		query.Set("opt_fields", strings.Join(projectOptFields, ","))
		if offset != "" {
			query.Set("offset", offset)
		}
		req, err := c.prepareRequest("GET", "projects", &query, nil)
		if err != nil {
			return nil, err
		}
		resp, err := c.do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			slog.Debug("wrong response on get projects", "status", resp.Status, "body", b)
			return nil, fmt.Errorf("response status not ok on get projects %s", resp.Status)
		}

		var body bodyResponse[*models.Project]

		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			return nil, err
		}
		projects = append(projects, body.Data...)
		offset = body.NextPage.Offset
		if offset == "" {
			break
		}
	}
	return projects, nil
}
