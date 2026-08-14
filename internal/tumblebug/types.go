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

// InfraStatusView is the status view of an Infra.
//
// CB-Tumblebug wraps this under a "status" key when option=status is used, and
// the shape differs from InfraInfo, so it needs its own type.
type InfraStatusView struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Status       string          `json:"status"`
	StatusCount  StatusCountInfo `json:"statusCount"`
	TargetStatus string          `json:"targetStatus"`
	TargetAction string          `json:"targetAction"`
}

// InfraAccessInfo is the access view of an Infra.
//
// The upstream struct carries no JSON tags on its outer fields, so the wire
// format uses the exported Go names verbatim.
type InfraAccessInfo struct {
	InfraID    string                `json:"InfraId"`
	NodeGroups []NodeGroupAccessInfo `json:"InfraNodeGroupAccessInfo"`
}

// NodeGroupAccessInfo is the access view of one NodeGroup.
type NodeGroupAccessInfo struct {
	NodeGroupID   string           `json:"NodeGroupId"`
	BastionNodeID string           `json:"BastionNodeId,omitempty"`
	Nodes         []NodeAccessInfo `json:"NodeAccessInfo"`
}

// NodeAccessInfo is how to reach one node.
//
// PrivateKey is only populated when the caller asks for it, and must never be
// logged or persisted.
type NodeAccessInfo struct {
	NodeID       string `json:"nodeId"`
	PublicIP     string `json:"publicIP"`
	PrivateIP    string `json:"privateIP"`
	SSHPort      int    `json:"sshPort"`
	NodeUserName string `json:"nodeUserName,omitempty"`
	PrivateKey   string `json:"privateKey,omitempty"`
}
