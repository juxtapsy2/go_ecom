package main

import (
	"log/slog"
	"os"
)

func main() {
	cfg := config{
		addr: ":9093",
		db:   dbConfig{},
	}

	api := &application{
		config: cfg,
	}

	// Structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	api.run(api.mount())

	if err := api.run(api.mount()); err != nil {
		slog.Error("server has failed to start", "error", err)
		os.Exit(1)
	}

}
