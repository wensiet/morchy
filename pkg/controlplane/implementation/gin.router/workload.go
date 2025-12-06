package ginrouter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wernsiet/morchy/pkg/controlplane/implementation/jsonformatter"
	"github.com/wernsiet/morchy/pkg/runtime"
)

// listWorkloads godoc
//
//	@Summary		List workloads
//	@Description	Retrieve a list of workloads filtered by status, CPU, and RAM
//	@Tags			workloads
//	@Accept			json
//	@Produce		json
//	@Param			status	query		string	false	"Filter by workload status"		example(new)
//	@Param			cpu		query		int		false	"Filter by CPU (millicores)"	example(100)
//	@Param			ram		query		int		false	"Filter by RAM (MB)"			example(256)
//	@Success		200		{array}		[]jsonformatter.WorkloadResponse
//	@Failure		400		{object}	map[string]string	"Invalid request parameters"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/api/v1/workloads [get]
func (rh *RouterHandler) listWorkloads(c *gin.Context) {
	var queryParams struct {
		StatusEq *string `form:"status" example:"new"`
		CPU      *uint   `form:"cpu" example:"100"`
		RAM      *uint   `form:"ram" example:"256"`
	}
	err := c.ShouldBindQuery(&queryParams)
	if err != nil {
		handleError(c, err)
		return
	}
	var resorceFilter *runtime.ResourceLimits = nil
	if queryParams.CPU != nil && queryParams.RAM != nil {
		resorceFilter = &runtime.ResourceLimits{
			CPU: *queryParams.CPU,
			RAM: *queryParams.RAM,
		}
	}

	workloads, err := rh.ucHandler.ListWorkloads(c, queryParams.StatusEq, resorceFilter)
	if err != nil {
		handleError(c, err)
		return
	}
	workloadApiModels := make([]*jsonformatter.WorkloadResponse, 0)
	for _, workload := range workloads {
		workloadApiModels = append(workloadApiModels, jsonformatter.NewWorkloadResponseFromDomain(workload))
	}
	c.JSON(http.StatusOK, workloadApiModels)
}

func (rh *RouterHandler) getWorkload(c *gin.Context) {
	workloadID := c.Param("workload_id")

	workload, err := rh.ucHandler.GetWorkload(c, workloadID)
	if err != nil {
		handleError(c, err)
		return
	}
	workloadApiModel := jsonformatter.NewWorkloadResponseFromDomain(workload)
	c.JSON(http.StatusOK, workloadApiModel)
}

func (rh *RouterHandler) createWorkload(c *gin.Context) {
	var workloadSpec jsonformatter.WorkloadSpecRequest
	err := c.ShouldBindJSON(&workloadSpec)
	if err != nil {
		handleError(c, err)
		return
	}

	domainWorkloadSpec := workloadSpec.ToDomain()
	createdWorkload, err := rh.ucHandler.CreateWorkload(c, domainWorkloadSpec)

	if err != nil {
		handleError(c, err)
		return
	}

	workloadApiModel := jsonformatter.NewWorkloadResponseFromDomain(createdWorkload)
	c.JSON(http.StatusCreated, workloadApiModel)
}
