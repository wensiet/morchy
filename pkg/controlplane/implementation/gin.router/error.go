package ginrouter

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/oops"
	"github.com/wernsiet/morchy/pkg/controlplane/domain"
)

func mapStatusCode(status string) int {
	switch status {
	case string(domain.NotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func handleError(c *gin.Context, err error) {
	if _, ok := oops.AsOops(err); !ok {
		err = domain.ErrorUnknownServerError.Wrap(err)
	}

	oopsErr, _ := oops.AsOops(err)

	raw, _ := oopsErr.MarshalJSON()
	var safeErr map[string]any
	json.Unmarshal(raw, &safeErr)
	delete(safeErr, "stacktrace")

	c.JSON(
		mapStatusCode(oopsErr.Code()),
		safeErr,
	)
}
