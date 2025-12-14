package ginrouter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wernsiet/morchy/pkg/controlplane/implementation/jsonformatter"
)

// putLease godoc
//
//	@Summary		Extend a lease
//	@Description	Extend a lease for a specific workload on a given node
//	@Tags			leases
//	@Accept			json
//	@Produce		json
//	@Param			workload_id	path		string	true	"Workload ID" minlength(1)
//	@Param			node_id		query		string	true	"Node ID"
//	@Success		200			{object}	jsonformatter.LeaseResponse
//	@Failure		400			{object}	map[string]string	"Invalid request parameters"
//	@Failure		500			{object}	map[string]string	"Internal server error"
//	@Router			/api/v1//workloads/{workload_id}/lease [put]
func (rh *RouterHandler) putLease(c *gin.Context) {
	wokrloadID := c.Param("workload_id")
	nodeID := c.Query("node_id")
	lease, err := rh.ucHandler.CreateOrExtendLease(c, nodeID, wokrloadID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, jsonformatter.NewLeaseResponseFromDomain(lease))
}
