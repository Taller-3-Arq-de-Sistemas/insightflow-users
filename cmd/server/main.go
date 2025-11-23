package main

import (
	"log/slog"
	"os"
)

func main() {
	cfg := config{
		addr: ":8080",
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	api := application{
		config: cfg,
	}

	h := api.mount()
	if err := api.run(h); err != nil {
		logger.Error("Unable to start server", "error", err)
		os.Exit(1)
	}
}
