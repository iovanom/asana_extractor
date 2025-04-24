package main

import (
	"log/slog"
	"os"
	"os/signal"

	"github.com/iovanom/asana_extractor/internal/asana"
	"github.com/iovanom/asana_extractor/internal/extractor"
	"github.com/iovanom/asana_extractor/internal/scheduler"
	"github.com/iovanom/asana_extractor/internal/storage"
	"github.com/joho/godotenv"
)

func main() {
	// set logger to debug
	// TODO: Read logger level from env
	slog.SetLogLoggerLevel(slog.LevelDebug)
	godotenv.Load()

	// init local storage dir
	dir := os.Getenv("STORAGE_DIR")
	if dir == "" {
		slog.Error("LOCAL_STORAGE_DIR is not set")
		os.Exit(1)
	}
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		slog.Error("failed to create local storage dir", "error", err)
		os.Exit(1)
	}

	s, err := storage.NewLocalStorage(dir)
	if err != nil {
		slog.Error("failed to init local storage", "error", err)
		os.Exit(1)
	}

	asanaClient, err := asana.NewClient(os.Getenv("ASANA_TOKEN"), os.Getenv("ASANA_WORKSPACE_ID"))
	if err != nil {
		slog.Error("failed to init asana client", "error", err)
	}

	asanaExtractor := extractor.NewExtractor(asanaClient, s)
	scheduler := scheduler.NewScheduler()

	scheduler.AddJob("@every 30s", extractJob(asanaExtractor))
	scheduler.AddJob("@every 5m", extractJob(asanaExtractor))
	scheduler.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	scheduler.Stop()
}

func extractJob(asanaExtractor *extractor.Extractor) scheduler.Job {
	return func() {
		err := asanaExtractor.ExtractData()
		if err != nil {
			slog.Error("error on extract data", "error", err)
		} else {
			slog.Info("data was extracted in storage directory")
		}
	}
}
