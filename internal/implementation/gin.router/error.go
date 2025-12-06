package ginrouter

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, map[string]string{"error": string(err.Error())})
}
