package tumblebug

// InfraDynamicReq is the CB-Tumblebug dynamic provisioning request.
// Only the fields this service sets or forwards are modelled.
type InfraDynamicReq struct {
	Name                   string                `json:"name"`
	Description            string                `json:"description,omitempty"`
	PolicyOnPartialFailure string                `json:"policyOnPartialFailure,omitempty"`
	InstallMonAgent        string                `json:"installMonAgent,omitempty"`
	NodeGroups             []NodeGroupDynamicReq `json:"nodeGroups"`
	Label                  map[string]string     `json:"label,omitempty"`
}

// NodeGroupDynamicReq describes one NodeGroup of a dynamic provisioning request.
type NodeGroupDynamicReq struct {
	Name          string            `json:"name,omitempty"`
	NodeGroupSize int               `json:"nodeGroupSize,omitempty"`
	SpecID        string            `json:"specId"`
	ImageID       string            `json:"imageId"`
	RootDiskType  string            `json:"rootDiskType,omitempty"`
	RootDiskSize  int               `json:"rootDiskSize,omitempty"`
	Zone          string            `json:"zone,omitempty"`
	Label         map[string]string `json:"label,omitempty"`
}

// InfraDynamicTemplateInfo is a stored Infra template.
type InfraDynamicTemplateInfo struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	InfraDynamicReq InfraDynamicReq `json:"infraDynamicReq"`
}

// ReviewResult is the pre-flight validation of a dynamic provisioning request.
type ReviewResult struct {
	OverallStatus  string   `json:"overallStatus"`
	OverallMessage string   `json:"overallMessage"`
	CreationViable bool     `json:"creationViable"`
	EstimatedCost  string   `json:"estimatedCost,omitempty"`
	InfraName      string   `json:"infraName"`
	TotalNodeCount int      `json:"totalNodeCount"`
	Recommendation []string `json:"recommendations,omitempty"`
}

// InfraInfo is the CB-Tumblebug view of a provisioned Infra.
type InfraInfo struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Status       string          `json:"status"`
	TargetStatus string          `json:"targetStatus"`
	StatusCount  StatusCountInfo `json:"statusCount"`
	Description  string          `json:"description"`
	Node         []NodeInfo      `json:"node"`
}

// StatusCountInfo summarizes node states of an Infra.
type StatusCountInfo struct {
	CountTotal    int `json:"countTotal"`
	CountCreating int `json:"countCreating"`
	CountRunning  int `json:"countRunning"`
}

// NodeInfo is a single provisioned node.
type NodeInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	NodeGroupID string `json:"nodeGroupId"`
	Status      string `json:"status"`
	PublicIP    string `json:"publicIP"`
	PrivateIP   string `json:"privateIP"`
	SSHPort     int    `json:"sshPort"`
	SpecID      string `json:"specId"`
	ImageID     string `json:"imageId"`
}

// CommandReq runs shell commands on the nodes of an Infra.
type CommandReq struct {
	Command        []string `json:"command"`
	UserName       string   `json:"userName,omitempty"`
	TimeoutMinutes int      `json:"timeoutMinutes,omitempty"`
}

// CommandResults is the response of a remote command execution.
type CommandResults struct {
	Results []CommandResult `json:"results"`
}

// CommandResult is the outcome of one command on one node.
type CommandResult struct {
	NodeID  string            `json:"nodeId"`
	Command map[string]string `json:"command"`
	Stdout  map[string]string `json:"stdout"`
	Stderr  map[string]string `json:"stderr"`
	Err     string            `json:"err"`
}
