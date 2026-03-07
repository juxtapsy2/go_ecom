package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/juxtapsy2/go_ecom/internal/env"
)

func main() {
	ctx := context.Background()

	cfg := config{
		addr: ":9093",
		db: dbConfig{
			dbConnectionString: env.GetString("GOOSE_DBSTRING", "host=localhost port=5435 user=postgres password=postgres dbname=goecom sslmode=disable"),
		},
	}

	// Structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Database
	conn, err := pgx.Connect(ctx, cfg.db.dbConnectionString)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	slog.Info("connected to database", "dsn", cfg.db.dbConnectionString)

	api := &application{
		config: cfg,
		db:     conn,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("server has failed to start", "error", err)
		os.Exit(1)
	}

}
