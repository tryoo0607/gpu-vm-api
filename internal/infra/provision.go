package infra

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/tryoo0607/gpu-vm-api/internal/tumblebug"
)

// gpuProbeCommand reports the GPUs visible inside a node.
var gpuProbeCommand = []string{"nvidia-smi --query-gpu=name,memory.total,driver_version --format=csv"}

const gpuProbeTimeoutMinutes = 5

// Create provisions a GPU Infra from the configured template.
//
// The template is the baseline only: its extra GPU NodeGroups are dropped so
// that a single NodeGroup is provisioned, mirroring the CB-MapUI flow of loading
// a template into the configuration and removing the NodeGroups that are not needed.
// The request is validated before anything is created.
func (s *Service) Create(ctx context.Context, nsID, name string) (*tumblebug.InfraInfo, error) {
	template, err := s.client.GetInfraTemplate(ctx, s.template.Namespace, s.template.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load provisioning template: %w", err)
	}

	req := template.InfraDynamicReq
	req.Name = name
	req.NodeGroups, err = selectNodeGroups(req.NodeGroups, s.template.SpecID)
	if err != nil {
		return nil, err
	}

	review, err := s.client.ReviewInfraDynamic(ctx, nsID, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate provisioning request: %w", err)
	}
	if !review.CreationViable {
		return nil, fmt.Errorf("provisioning request is not viable: %s (%s)", review.OverallMessage, review.OverallStatus)
	}

	log.Info().Str("nsId", nsID).Str("infraId", name).
		Str("estimatedCost", review.EstimatedCost).Int("nodeCount", review.TotalNodeCount).
		Msg("Creating infra")

	created, err := s.client.CreateInfraDynamic(ctx, nsID, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to create infra: %w", err)
	}

	log.Info().Str("nsId", nsID).Str("infraId", created.ID).Str("status", created.Status).
		Msg("Created infra")
	return created, nil
}

// selectNodeGroups keeps only the NodeGroup matching specID.
//
// Provisioning every NodeGroup of the GPU template would multiply the hourly
// cost, so an unmatched spec is an error rather than a silent full deployment.
func selectNodeGroups(groups []tumblebug.NodeGroupDynamicReq, specID string) ([]tumblebug.NodeGroupDynamicReq, error) {
	selected := make([]tumblebug.NodeGroupDynamicReq, 0, 1)
	for _, group := range groups {
		if group.SpecID == specID {
			selected = append(selected, group)
		}
	}
	if len(selected) == 0 {
		return nil, fmt.Errorf("template has no node group with spec %q", specID)
	}
	return selected, nil
}

// Get reads one Infra.
func (s *Service) Get(ctx context.Context, nsID, infraID string) (*tumblebug.InfraInfo, error) {
	result, err := s.client.GetInfra(ctx, nsID, infraID)
	if err != nil {
		return nil, fmt.Errorf("failed to get infra: %w", err)
	}
	return result, nil
}

// Status reads the node status summary of an Infra.
func (s *Service) Status(ctx context.Context, nsID, infraID string) (*tumblebug.InfraStatusView, error) {
	result, err := s.client.GetInfraStatus(ctx, nsID, infraID)
	if err != nil {
		return nil, fmt.Errorf("failed to get infra status: %w", err)
	}
	return result, nil
}

// AccessInfo reads how to reach the nodes of an Infra.
func (s *Service) AccessInfo(ctx context.Context, nsID, infraID string, showSSHKey bool) (*tumblebug.InfraAccessInfo, error) {
	result, err := s.client.GetInfraAccessInfo(ctx, nsID, infraID, showSSHKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get infra access info: %w", err)
	}
	return result, nil
}

// List reads every Infra of a namespace.
func (s *Service) List(ctx context.Context, nsID string) ([]tumblebug.InfraInfo, error) {
	result, err := s.client.ListInfra(ctx, nsID)
	if err != nil {
		return nil, fmt.Errorf("failed to list infra: %w", err)
	}
	return result, nil
}

// Control applies a lifecycle action to every node of an Infra.
func (s *Service) Control(ctx context.Context, nsID, infraID, action string) ([]string, error) {
	log.Info().Str("nsId", nsID).Str("infraId", infraID).Str("action", action).Msg("Controlling infra")

	result, err := s.client.ControlInfra(ctx, nsID, infraID, action)
	if err != nil {
		return nil, fmt.Errorf("failed to control infra: %w", err)
	}
	return result.Output, nil
}

// GPUInfo runs nvidia-smi on every node and returns the raw output per node.
func (s *Service) GPUInfo(ctx context.Context, nsID, infraID string) (*tumblebug.CommandResults, error) {
	req := &tumblebug.CommandReq{
		Command:        gpuProbeCommand,
		UserName:       s.template.NodeUserName,
		TimeoutMinutes: gpuProbeTimeoutMinutes,
	}

	result, err := s.client.RunCommand(ctx, nsID, infraID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to probe gpu: %w", err)
	}
	return result, nil
}
