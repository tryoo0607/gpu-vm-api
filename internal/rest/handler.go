package rest

import "github.com/tryoo0607/gpu-vm-api/internal/infra"

// Handler serves the GPU VM lifecycle endpoints.
type Handler struct {
	infra *infra.Service
}

// NewHandler builds a Handler around the lifecycle service.
func NewHandler(service *infra.Service) *Handler {
	return &Handler{infra: service}
}
