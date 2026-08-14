package tumblebug

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// Infra query options accepted by CB-Tumblebug.
const (
	OptionStatus      = "status"
	OptionAccessInfo  = "accessinfo"
	accessInfoShowKey = "showSshKey"
)

// GetInfraTemplate reads a stored Infra template.
//
// Templates are looked up per namespace with no fallback, so nsID must be the
// namespace that actually holds the template (the built-in ones live in "system").
func (c *Client) GetInfraTemplate(ctx context.Context, nsID, templateID string) (*InfraDynamicTemplateInfo, error) {
	path := fmt.Sprintf("/ns/%s/template/infra/%s", url.PathEscape(nsID), url.PathEscape(templateID))

	result := &InfraDynamicTemplateInfo{}
	if err := c.Do(ctx, http.MethodGet, path, nil, nil, result); err != nil {
		return nil, fmt.Errorf("failed to get infra template %q in namespace %q: %w", templateID, nsID, err)
	}
	return result, nil
}

// ReviewInfraDynamic validates a provisioning request without creating anything.
func (c *Client) ReviewInfraDynamic(ctx context.Context, nsID string, req *InfraDynamicReq) (*ReviewResult, error) {
	path := fmt.Sprintf("/ns/%s/infraDynamicReview", url.PathEscape(nsID))

	result := &ReviewResult{}
	if err := c.Do(ctx, http.MethodPost, path, nil, req, result); err != nil {
		return nil, fmt.Errorf("failed to review infra request: %w", err)
	}
	return result, nil
}

// CreateInfraDynamic provisions an Infra and blocks until provisioning finishes.
//
// CB-Tumblebug performs this synchronously; cancelling the context aborts
// provisioning and can leave partially created, billable resources behind.
func (c *Client) CreateInfraDynamic(ctx context.Context, nsID string, req *InfraDynamicReq) (*InfraInfo, error) {
	path := fmt.Sprintf("/ns/%s/infraDynamic", url.PathEscape(nsID))

	result := &InfraInfo{}
	if err := c.Do(ctx, http.MethodPost, path, nil, req, result); err != nil {
		return nil, fmt.Errorf("failed to create infra %q: %w", req.Name, err)
	}
	return result, nil
}

// GetInfra reads one Infra. Pass option "status" or "accessinfo" to change the view.
func (c *Client) GetInfra(ctx context.Context, nsID, infraID, option string) (*InfraInfo, error) {
	query := url.Values{}
	if option != "" {
		query.Set("option", option)
	}
	if option == OptionAccessInfo {
		query.Set("accessInfoOption", accessInfoShowKey)
	}
	path := fmt.Sprintf("/ns/%s/infra/%s", url.PathEscape(nsID), url.PathEscape(infraID))

	result := &InfraInfo{}
	if err := c.Do(ctx, http.MethodGet, path, query, nil, result); err != nil {
		return nil, fmt.Errorf("failed to get infra %q: %w", infraID, err)
	}
	return result, nil
}

// ListInfra reads every Infra of a namespace with its status.
func (c *Client) ListInfra(ctx context.Context, nsID string) ([]InfraInfo, error) {
	query := url.Values{"option": []string{OptionStatus}}
	path := fmt.Sprintf("/ns/%s/infra", url.PathEscape(nsID))

	var result struct {
		Infra []InfraInfo `json:"infra"`
	}
	if err := c.Do(ctx, http.MethodGet, path, query, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to list infra in namespace %q: %w", nsID, err)
	}
	return result.Infra, nil
}

// ControlInfra applies a lifecycle action to every node of an Infra.
//
// CB-Tumblebug exposes this as GET; this service takes POST inbound because the
// call changes state, and translates it here.
func (c *Client) ControlInfra(ctx context.Context, nsID, infraID, action string) (*IDList, error) {
	query := url.Values{"action": []string{action}}
	path := fmt.Sprintf("/ns/%s/control/infra/%s", url.PathEscape(nsID), url.PathEscape(infraID))

	result := &IDList{}
	if err := c.Do(ctx, http.MethodGet, path, query, nil, result); err != nil {
		return nil, fmt.Errorf("failed to run action %q on infra %q: %w", action, infraID, err)
	}
	return result, nil
}

// RunCommand executes shell commands on every node of an Infra.
func (c *Client) RunCommand(ctx context.Context, nsID, infraID string, req *CommandReq) (*CommandResults, error) {
	path := fmt.Sprintf("/ns/%s/cmd/infra/%s", url.PathEscape(nsID), url.PathEscape(infraID))

	result := &CommandResults{}
	if err := c.Do(ctx, http.MethodPost, path, nil, req, result); err != nil {
		return nil, fmt.Errorf("failed to run command on infra %q: %w", infraID, err)
	}
	return result, nil
}
