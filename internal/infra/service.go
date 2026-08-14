// Package infra holds the GPU VM lifecycle logic on top of CB-Tumblebug.
package infra

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/tryoo0607/gpu-vm-api/internal/config"
	"github.com/tryoo0607/gpu-vm-api/internal/tumblebug"
)

// Service drives the GPU VM lifecycle through CB-Tumblebug.
type Service struct {
	client   *tumblebug.Client
	template config.TemplateConfig
}

// NewService builds a Service around a CB-Tumblebug client.
func NewService(client *tumblebug.Client, template config.TemplateConfig) *Service {
	return &Service{client: client, template: template}
}

// Delete terminates every node of an Infra and removes its records.
//
// Shared resources are not touched here: CB-Tumblebug has no cascade, and the
// Infra must be gone before ReleaseSharedResources can succeed.
func (s *Service) Delete(ctx context.Context, nsID, infraID string) ([]string, error) {
	log.Info().Str("nsId", nsID).Str("infraId", infraID).Msg("Deleting infra")

	result, err := s.client.DeleteInfra(ctx, nsID, infraID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete infra: %w", err)
	}

	log.Info().Str("nsId", nsID).Str("infraId", infraID).Int("deleted", len(result.Output)).
		Msg("Deleted infra")
	return result.Output, nil
}

// ReleaseSharedResources deletes the shared vNet, security group and SSH key of a namespace.
// Run it after every Infra in the namespace is deleted.
func (s *Service) ReleaseSharedResources(ctx context.Context, nsID string) ([]tumblebug.ResourceDeleteResult, error) {
	log.Info().Str("nsId", nsID).Msg("Releasing shared resources")

	result, err := s.client.ReleaseSharedResources(ctx, nsID)
	if err != nil {
		return nil, fmt.Errorf("failed to release shared resources: %w", err)
	}

	log.Info().Str("nsId", nsID).Int("released", len(result.Output)).Msg("Released shared resources")
	return result.Output, nil
}
