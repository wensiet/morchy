package ginrouter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wernsiet/morchy/internal/implementation/jsonformatter"
)

// createLease godoc
//
//	@Summary		Create a lease
//	@Description	Create a lease for a specific workload on a given node
//	@Tags			leases
//	@Accept			json
//	@Produce		json
//	@Param			workload_id	path		string	true	"Workload ID" minlength(1)
//	@Param			node_id		query		string	true	"Node ID"
//	@Success		200			{object}	jsonformatter.LeaseResponse
//	@Failure		400			{object}	map[string]string	"Invalid request parameters"
//	@Failure		500			{object}	map[string]string	"Internal server error"
//	@Router			/api/v1//workloads/{workload_id}/lease [post]
func (rh *RouterHandler) createLease(c *gin.Context) {
	wokrloadID := c.Param("workload_id")
	nodeID := c.Query("node_id")
	lease, err := rh.ucHandler.CreateLease(c, nodeID, wokrloadID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, jsonformatter.NewLeaseResponseFromDomain(lease))
}

// extendLease godoc
//
//	@Summary		Extend a lease
//	@Description	Extend a lease for a specific workload on a given node
//	@Tags			leases
//	@Accept			json
//	@Produce		json
//	@Param			workload_id	path		string	true	"Workload ID" minlength(1)
//	@Param			node_id		query		string	true	"Node ID"
//	@Success		200			{object}	nil
//	@Failure		400			{object}	map[string]string	"Invalid request parameters"
//	@Failure		500			{object}	map[string]string	"Internal server error"
//	@Router			/api/v1//workloads/{workload_id}/lease [put]
func (rh *RouterHandler) extendLease(c *gin.Context) {
	wokrloadID := c.Param("workload_id")
	nodeID := c.Query("node_id")
	err := rh.ucHandler.ExtendLease(c, nodeID, wokrloadID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}
