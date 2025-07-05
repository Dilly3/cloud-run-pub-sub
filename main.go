package main

import (
	"log/slog"
	"os"

	"github.com/dilly3/cloud-run-pub-sub/server"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	s := server.NewServer()
	srv := s.SetupServer(s.Config.Port)

	go server.GracefulShutdown(srv, logger)
	err := srv.ListenAndServe()
	if err != nil {
		logger.Error("Failed to listen and serve", "error", err)
	}

}
