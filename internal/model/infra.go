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

// CreateInfraReq asks for a GPU Infra built from the configured template.
type CreateInfraReq struct {
	// Name becomes the Infra ID; node names get a -N suffix.
	Name string `json:"name" validate:"required" example:"gpu-lab"`
}

// ControlInfraReq applies a lifecycle action to every node of an Infra.
type ControlInfraReq struct {
	// Action is one of suspend, resume, reboot, terminate.
	Action string `json:"action" validate:"required" example:"suspend" enums:"suspend,resume,reboot,terminate"`
}

// ControlInfraResp reports the nodes an action was applied to.
type ControlInfraResp struct {
	Affected []string `json:"affected"`
}
