package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/tryoo0607/gpu-vm-api/internal/model"
	"github.com/tryoo0607/gpu-vm-api/internal/tumblebug"
)

// statusCodeFor maps an upstream failure to the status this API returns.
func statusCodeFor(err error) int {
	switch tumblebug.StatusCode(err) {
	case http.StatusNotFound:
		return http.StatusNotFound
	case http.StatusBadRequest, http.StatusUnprocessableEntity:
		return http.StatusBadRequest
	case http.StatusConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// RestDeleteInfra godoc
// @ID DeleteInfra
// @Summary Delete a GPU Infra
// @Description Terminate every node of the Infra and remove its records. Shared resources are kept; release them separately once the namespace holds no Infra
// @Tags [Infra] GPU VM Lifecycle
// @Produce json
// @Param nsId path string true "Namespace ID" default(default)
// @Param infraId path string true "Infra ID"
// @Success 200 {object} model.DeleteInfraResp "Infra deleted"
// @Failure 404 {object} model.SimpleMsg "Infra not found"
// @Failure 500 {object} model.SimpleMsg "Infra deletion failed"
// @Router /ns/{nsId}/infra/{infraId} [delete]
func (h *Handler) RestDeleteInfra(c echo.Context) error {
	ctx := c.Request().Context()
	nsID := c.Param("nsId")
	infraID := c.Param("infraId")

	deleted, err := h.infra.Delete(ctx, nsID, infraID)
	if err != nil {
		log.Error().Err(err).Str("nsId", nsID).Str("infraId", infraID).Msg("Failed to delete infra")
		return c.JSON(statusCodeFor(err), model.SimpleMsg{Message: "Infra deletion failed"})
	}
	return c.JSON(http.StatusOK, model.DeleteInfraResp{Deleted: deleted})
}

// RestDeleteSharedResources godoc
// @ID DeleteSharedResources
// @Summary Release shared resources
// @Description Delete the shared vNet, security group and SSH key of the namespace. Delete every Infra first: CB-Tumblebug does not cascade and the CSP rejects dependent resources
// @Tags [Infra] GPU VM Lifecycle
// @Produce json
// @Param nsId path string true "Namespace ID" default(default)
// @Success 200 {object} model.ReleaseSharedResourcesResp "Shared resources released"
// @Failure 500 {object} model.SimpleMsg "Shared resource release failed"
// @Router /ns/{nsId}/shared-resources [delete]
func (h *Handler) RestDeleteSharedResources(c echo.Context) error {
	ctx := c.Request().Context()
	nsID := c.Param("nsId")

	released, err := h.infra.ReleaseSharedResources(ctx, nsID)
	if err != nil {
		log.Error().Err(err).Str("nsId", nsID).Msg("Failed to release shared resources")
		return c.JSON(statusCodeFor(err), model.SimpleMsg{Message: "Shared resource release failed"})
	}

	results := make([]model.ResourceResult, 0, len(released))
	for _, item := range released {
		results = append(results, model.ResourceResult{ID: item.ID, Message: item.Message})
	}
	return c.JSON(http.StatusOK, model.ReleaseSharedResourcesResp{Released: results})
}
