package scheduler

import (
	"log/slog"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	c *cron.Cron
}

type Job func()

func NewScheduler() *Scheduler {
	c := cron.New()
	return &Scheduler{c: c}
}

func (s *Scheduler) AddJob(cronSpec string, job Job) error {
	_, err := s.c.AddFunc(cronSpec, job)
	return err
}

func (s *Scheduler) Start() {
	slog.Info("Scheduler starting")
	s.c.Start()
}

func (s *Scheduler) Stop() {
	slog.Info("Scheduler stop")
	s.c.Stop()
}
