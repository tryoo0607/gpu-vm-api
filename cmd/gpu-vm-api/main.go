// Command gpu-vm-api serves a REST API that drives the GPU VM lifecycle on CB-Tumblebug.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/tryoo0607/gpu-vm-api/internal/config"
	"github.com/tryoo0607/gpu-vm-api/internal/rest"
)

const shutdownTimeout = 10 * time.Second

// @title GPU VM API
// @version 0.1.0
// @description REST API for the Ubuntu GPU VM lifecycle on CB-Tumblebug.
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /gpuvm
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Error().Err(err).Msg("Failed to load configuration")
		os.Exit(1)
	}

	server := rest.NewServer(cfg)
	rest.MarkReady()

	go func() {
		log.Info().Int("port", server.Port()).Str("basePath", rest.BasePath).Msg("Starting server")
		if err := server.Start(); err != nil {
			log.Error().Err(err).Msg("Server stopped unexpectedly")
			os.Exit(1)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Graceful shutdown failed")
		os.Exit(1)
	}
	log.Info().Msg("Server stopped")
}
