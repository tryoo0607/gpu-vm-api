package rest

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/tryoo0607/gpu-vm-api/internal/model"
)

// allowedActions are the lifecycle actions this API forwards.
//
// Recovery actions (refine, reconcile, abort) are intentionally excluded: they
// change resource bookkeeping in ways that need a deliberate decision.
var allowedActions = map[string]bool{
	"suspend":   true,
	"resume":    true,
	"reboot":    true,
	"terminate": true,
}

// RestPostInfra godoc
// @ID PostInfra
// @Summary Create a GPU Infra
// @Description Provision an Ubuntu GPU VM from the configured template. The request is validated first and blocks until provisioning finishes, which takes several minutes
// @Tags [Infra] GPU VM Lifecycle
// @Accept json
// @Produce json
// @Param nsId path string true "Namespace ID" default(default)
// @Param infraReq body model.CreateInfraReq true "Infra creation request"
// @Success 200 {object} tumblebug.InfraInfo "Infra created"
// @Failure 400 {object} model.SimpleMsg "Malformed request body"
// @Failure 500 {object} model.SimpleMsg "Infra creation failed"
// @Router /ns/{nsId}/infra [post]
func (h *Handler) RestPostInfra(c echo.Context) error {
	ctx := c.Request().Context()
	nsID := c.Param("nsId")

	var req model.CreateInfraReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.SimpleMsg{Message: "Malformed request body: check JSON syntax"})
	}
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, model.SimpleMsg{Message: "Infra name required"})
	}

	created, err := h.infra.Create(ctx, nsID, req.Name)
	if err != nil {
		log.Error().Err(err).Str("nsId", nsID).Str("name", req.Name).Msg("Failed to create infra")
		return c.JSON(statusCodeFor(err), model.SimpleMsg{Message: "Infra creation failed"})
	}
	return c.JSON(http.StatusOK, created)
}

// RestGetAllInfra godoc
// @ID GetAllInfra
// @Summary List GPU Infras
// @Description List every Infra of the namespace with its node status
// @Tags [Infra] GPU VM Lifecycle
// @Produce json
// @Param nsId path string true "Namespace ID" default(default)
// @Success 200 {array} tumblebug.InfraInfo "Infra list"
// @Failure 500 {object} model.SimpleMsg "Infra lookup failed"
// @Router /ns/{nsId}/infra [get]
func (h *Handler) RestGetAllInfra(c echo.Context) error {
	ctx := c.Request().Context()
	nsID := c.Param("nsId")

	list, err := h.infra.List(ctx, nsID)
	if err != nil {
		log.Error().Err(err).Str("nsId", nsID).Msg("Failed to list infra")
		return c.JSON(statusCodeFor(err), model.SimpleMsg{Message: "Infra lookup failed"})
	}
	return c.JSON(http.StatusOK, list)
}

// RestGetInfra godoc
// @ID GetInfra
// @Summary Get a GPU Infra
// @Description Read one Infra with its nodes
// @Tags [Infra] GPU VM Lifecycle
// @Produce json
// @Param nsId path string true "Namespace ID" default(default)
// @Param infraId path string true "Infra ID"
// @Success 200 {object} tumblebug.InfraInfo "Infra information"
// @Failure 404 {object} model.SimpleMsg "Infra not found"
// @Failure 500 {object} model.SimpleMsg "Infra lookup failed"
// @Router /ns/{nsId}/infra/{infraId} [get]
func (h *Handler) RestGetInfra(c echo.Context) error {
	ctx := c.Request().Context()
	nsID := c.Param("nsId")
	infraID := c.Param("infraId")

	result, err := h.infra.Get(ctx, nsID, infraID)
	if err != nil {
		log.Error().Err(err).Str("nsId", nsID).Str("infraId", infraID).Msg("Failed to get infra")
		return c.JSON(statusCodeFor(err), model.SimpleMsg{Message: "Infra lookup failed"})
	}
	return c.JSON(http.StatusOK, result)
}

// RestGetInfraStatus godoc
// @ID GetInfraStatus
// @Summary Get the status of a GPU Infra
// @Description Report the node status summary, suitable for polling while an Infra is being created
// @Tags [Infra] GPU VM Lifecycle
// @Produce json
// @Param nsId path string true "Namespace ID" default(default)
// @Param infraId path string true "Infra ID"
// @Success 200 {object} tumblebug.InfraStatusView "Infra status"
// @Failure 404 {object} model.SimpleMsg "Infra not found"
// @Failure 500 {object} model.SimpleMsg "Infra status lookup failed"
// @Router /ns/{nsId}/infra/{infraId}/status [get]
func (h *Handler) RestGetInfraStatus(c echo.Context) error {
	ctx := c.Request().Context()
	nsID := c.Param("nsId")
	infraID := c.Param("infraId")

	result, err := h.infra.Status(ctx, nsID, infraID)
	if err != nil {
		log.Error().Err(err).Str("nsId", nsID).Str("infraId", infraID).Msg("Failed to get infra status")
		return c.JSON(statusCodeFor(err), model.SimpleMsg{Message: "Infra status lookup failed"})
	}
	return c.JSON(http.StatusOK, result)
}

// RestGetInfraAccess godoc
// @ID GetInfraAccess
// @Summary Get access information of a GPU Infra
// @Description Report public IP, SSH port and login user per node. The SSH private key is returned only when showSshKey is true
// @Tags [Infra] GPU VM Lifecycle
// @Produce json
// @Param nsId path string true "Namespace ID" default(default)
// @Param infraId path string true "Infra ID"
// @Param showSshKey query bool false "Include the SSH private key" default(false)
// @Success 200 {object} tumblebug.InfraAccessInfo "Access information"
// @Failure 404 {object} model.SimpleMsg "Infra not found"
// @Failure 500 {object} model.SimpleMsg "Access information lookup failed"
// @Router /ns/{nsId}/infra/{infraId}/access [get]
func (h *Handler) RestGetInfraAccess(c echo.Context) error {
	ctx := c.Request().Context()
	nsID := c.Param("nsId")
	infraID := c.Param("infraId")
	showSSHKey := strings.EqualFold(c.QueryParam("showSshKey"), "true")

	result, err := h.infra.AccessInfo(ctx, nsID, infraID, showSSHKey)
	if err != nil {
		log.Error().Err(err).Str("nsId", nsID).Str("infraId", infraID).Msg("Failed to get access info")
		return c.JSON(statusCodeFor(err), model.SimpleMsg{Message: "Access information lookup failed"})
	}
	return c.JSON(http.StatusOK, result)
}

// RestGetInfraGPU godoc
// @ID GetInfraGPU
// @Summary Report GPUs of an Infra
// @Description Run nvidia-smi on every node and return the reported GPU model, memory and driver version
// @Tags [Infra] GPU VM Lifecycle
// @Produce json
// @Param nsId path string true "Namespace ID" default(default)
// @Param infraId path string true "Infra ID"
// @Success 200 {object} tumblebug.CommandResults "GPU report per node"
// @Failure 404 {object} model.SimpleMsg "Infra not found"
// @Failure 500 {object} model.SimpleMsg "GPU report failed"
// @Router /ns/{nsId}/infra/{infraId}/gpu [get]
func (h *Handler) RestGetInfraGPU(c echo.Context) error {
	ctx := c.Request().Context()
	nsID := c.Param("nsId")
	infraID := c.Param("infraId")

	result, err := h.infra.GPUInfo(ctx, nsID, infraID)
	if err != nil {
		log.Error().Err(err).Str("nsId", nsID).Str("infraId", infraID).Msg("Failed to probe gpu")
		return c.JSON(statusCodeFor(err), model.SimpleMsg{Message: "GPU report failed"})
	}
	return c.JSON(http.StatusOK, result)
}

// RestPostInfraControl godoc
// @ID PostInfraControl
// @Summary Control a GPU Infra
// @Description Apply suspend, resume, reboot or terminate to every node. Terminate stops the nodes but keeps the records, so delete the Infra afterwards to stop billing
// @Tags [Infra] GPU VM Lifecycle
// @Accept json
// @Produce json
// @Param nsId path string true "Namespace ID" default(default)
// @Param infraId path string true "Infra ID"
// @Param controlReq body model.ControlInfraReq true "Lifecycle action"
// @Success 200 {object} model.ControlInfraResp "Action applied"
// @Failure 400 {object} model.SimpleMsg "Unsupported action"
// @Failure 404 {object} model.SimpleMsg "Infra not found"
// @Failure 500 {object} model.SimpleMsg "Infra control failed"
// @Router /ns/{nsId}/infra/{infraId}/control [post]
func (h *Handler) RestPostInfraControl(c echo.Context) error {
	ctx := c.Request().Context()
	nsID := c.Param("nsId")
	infraID := c.Param("infraId")

	var req model.ControlInfraReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.SimpleMsg{Message: "Malformed request body: check JSON syntax"})
	}
	if !allowedActions[req.Action] {
		return c.JSON(http.StatusBadRequest, model.SimpleMsg{Message: "Action must be suspend, resume, reboot or terminate"})
	}

	affected, err := h.infra.Control(ctx, nsID, infraID, req.Action)
	if err != nil {
		log.Error().Err(err).Str("nsId", nsID).Str("infraId", infraID).Str("action", req.Action).
			Msg("Failed to control infra")
		return c.JSON(statusCodeFor(err), model.SimpleMsg{Message: "Infra control failed"})
	}
	return c.JSON(http.StatusOK, model.ControlInfraResp{Affected: affected})
}
