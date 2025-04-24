package extractor

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/iovanom/asana_extractor/internal/models"
	"golang.org/x/sync/errgroup"
)

const (
	extractTimeout = 60 // in seconds
)

type client interface {
	Users() ([]*models.User, error)
	Projects() ([]*models.Project, error)
}

type storage interface {
	SaveFile(filename string, body io.Reader) error
}

type Extractor struct {
	client  client
	storage storage
}

func NewExtractor(c client, s storage) *Extractor {
	return &Extractor{c, s}
}

func (e *Extractor) ExtractData() error {
	slog.Debug("start extracting data")
	ctx, cancel := context.WithTimeout(context.Background(), extractTimeout*time.Second)
	defer cancel()
	g, _ := errgroup.WithContext(ctx)
	g.Go(e.ExtractUsers)
	g.Go(e.ExtractProjects)
	return g.Wait()
}

func (e *Extractor) ExtractUsers() error {
	slog.Debug("start extracting users")
	users, err := e.client.Users()
	if err != nil {
		return err
	}
	slog.Debug("users length", "length", len(users))
	for _, user := range users {
		b, err := json.Marshal(user)
		slog.Debug("user", "user", b)
		if err != nil {
			return err
		}
		err = e.storage.SaveFile("user_"+user.ID+".json", bytes.NewReader(b))
		if err != nil {
			return err
		}

	}
	return nil
}

func (e *Extractor) ExtractProjects() error {
	slog.Debug("start extracting projects")
	projects, err := e.client.Projects()
	if err != nil {
		return err
	}
	slog.Debug("projects length", "length", len(projects))
	for _, project := range projects {
		b, err := json.Marshal(project)
		slog.Debug("project", "project", b)
		if err != nil {
			return err
		}
		err = e.storage.SaveFile("project_"+project.ID+".json", bytes.NewReader(b))
		if err != nil {
			return err
		}
	}
	return nil
}
