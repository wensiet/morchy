package ginrouter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wernsiet/morchy/pkg/controlplane/domain"
	"github.com/wernsiet/morchy/pkg/controlplane/implementation/jsonformatter"
)

// pushEvent godoc
//
//	@Summary		Push event
//	@Description	Create a new event for a node
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			node_id	query		string								true	"Node ID"
//	@Param			body	body		jsonformatter.EventCreateRequest	true	"Event payload"
//	@Success		201		{object}	nil
//	@Router			/api/v1/events [post]
func (rh *RouterHandler) pushEvent(c *gin.Context) {
	node_id := c.Query("node_id")
	if node_id == "" {
		handleError(c, domain.ErrorRequestParamsValidation.New("node_id field is required"))
		return
	}

	var eventCreate jsonformatter.EventCreateRequest
	err := c.ShouldBindBodyWithJSON(&eventCreate)
	if err != nil {
		handleError(c, domain.ErrorRequestParamsValidation.Wrapf(err, "request body is invalid"))
		return
	}

	err = rh.ucHandler.PushEvent(c, eventCreate.ToDomain(node_id))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, nil)
}
