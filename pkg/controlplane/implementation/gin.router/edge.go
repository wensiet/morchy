package ginrouter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wernsiet/morchy/pkg/controlplane/implementation/jsonformatter"
)

// listEdges godoc
//
//	@Summary		List edges
//	@Description	Retrieve a list of all edges
//	@Tags			edges
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		jsonformatter.EdgeResponse
//	@Failure		400	{object}	map[string]string	"Invalid request parameters"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/api/v1/edges [get]
func (rh *RouterHandler) listEdges(c *gin.Context) {
	edges, err := rh.ucHandler.ListEdges(c)
	if err != nil {
		handleError(c, err)
		return
	}

	edgeApiModels := make([]*jsonformatter.EdgeResponse, 0)
	for _, edge := range edges {
		edgeApiModels = append(edgeApiModels, jsonformatter.NewEdgeFromDomain(edge))
	}

	c.JSON(http.StatusOK, edgeApiModels)
}
