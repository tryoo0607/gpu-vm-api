package model

// DeleteInfraResp reports the resources removed with an Infra.
type DeleteInfraResp struct {
	Deleted []string `json:"deleted"`
}

// ResourceResult is the outcome of releasing one shared resource.
type ResourceResult struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// ReleaseSharedResourcesResp reports the shared resources released in a namespace.
type ReleaseSharedResourcesResp struct {
	Released []ResourceResult `json:"released"`
}
