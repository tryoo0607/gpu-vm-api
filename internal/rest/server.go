// Package rest wires the Echo HTTP server and its routes.
package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/tryoo0607/gpu-vm-api/internal/config"
	"github.com/tryoo0607/gpu-vm-api/internal/infra"
	"github.com/tryoo0607/gpu-vm-api/internal/tumblebug"
)

// BasePath is the mount point for every route of this service.
const BasePath = "/gpuvm"

// Server owns the Echo instance and its lifecycle.
type Server struct {
	echo *echo.Echo
	port int
}

// NewServer builds the Echo instance and registers all routes.
func NewServer(cfg *config.Config) *Server {
	handler := NewHandler(infra.NewService(tumblebug.NewClient(cfg.Tumblebug)))

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	g := e.Group(BasePath)
	g.GET("/readyz", RestGetReadyz)
	g.DELETE("/ns/:nsId/infra/:infraId", handler.RestDeleteInfra)
	g.DELETE("/ns/:nsId/shared-resources", handler.RestDeleteSharedResources)

	return &Server{echo: e, port: cfg.Server.Port}
}

// Start blocks serving HTTP until the server is shut down.
func (s *Server) Start() error {
	if err := s.echo.Start(fmt.Sprintf(":%d", s.port)); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http server stopped: %w", err)
	}
	return nil
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.echo.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shut down http server: %w", err)
	}
	return nil
}

// Port reports the port the server listens on.
func (s *Server) Port() int {
	return s.port
}
