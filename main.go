package main

import (
	"log/slog"
	"os"

	"github.com/dilly3/cloud-run-pub-sub/config"
	"github.com/dilly3/cloud-run-pub-sub/publisher"
	"github.com/dilly3/cloud-run-pub-sub/server"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	config, err := config.GetConfig(logger)
	if err != nil {
		logger.Error("Failed to get config", "error", err)
		os.Exit(1)
	}
	logger.Info("Config", "config", config)
	publisher := publisher.NewPublisher(config.ProjectID, config.TopicID, logger)
	s := server.NewServer(logger, publisher)
	srv := s.SetupServer(s.Config.Port)

	go server.GracefulShutdown(srv, logger)
	err = srv.ListenAndServe()
	if err != nil {
		logger.Error("Failed to listen and serve", "error", err)
	}

}
