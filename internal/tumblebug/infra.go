package tumblebug

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// DeleteOptionTerminate terminates CSP nodes before deleting CB-Tumblebug records.
//
// The "force" option is deliberately not exposed: it drops records without
// confirming CSP termination, which leaves billed orphan instances behind and
// blocks VNet/SecurityGroup cleanup.
const DeleteOptionTerminate = "terminate"

// IDList is the CB-Tumblebug response carrying affected resource IDs.
type IDList struct {
	Output []string `json:"output"`
}

// ResourceDeleteResults is the CB-Tumblebug response for shared resource release.
type ResourceDeleteResults struct {
	Output []ResourceDeleteResult `json:"output"`
}

// ResourceDeleteResult is the per-resource outcome of a shared resource release.
type ResourceDeleteResult struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Stderr  string `json:"stderr"`
}

// DeleteInfra terminates every node of an Infra and removes its records.
func (c *Client) DeleteInfra(ctx context.Context, nsID, infraID string) (*IDList, error) {
	query := url.Values{"option": []string{DeleteOptionTerminate}}
	path := fmt.Sprintf("/ns/%s/infra/%s", url.PathEscape(nsID), url.PathEscape(infraID))

	result := &IDList{}
	if err := c.Do(ctx, http.MethodDelete, path, query, nil, result); err != nil {
		return nil, fmt.Errorf("failed to delete infra %q: %w", infraID, err)
	}
	return result, nil
}

// ReleaseSharedResources deletes the shared vNet, security group and SSH key of a namespace.
//
// CB-Tumblebug does not cascade: an Infra must be deleted before its shared
// resources can be released, otherwise the CSP reports a dependency violation.
func (c *Client) ReleaseSharedResources(ctx context.Context, nsID string) (*ResourceDeleteResults, error) {
	path := fmt.Sprintf("/ns/%s/sharedResources", url.PathEscape(nsID))

	result := &ResourceDeleteResults{}
	if err := c.Do(ctx, http.MethodDelete, path, nil, nil, result); err != nil {
		return nil, fmt.Errorf("failed to release shared resources in namespace %q: %w", nsID, err)
	}
	return result, nil
}
