package ginrouter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wernsiet/morchy/pkg/controlplane/domain"
	"github.com/wernsiet/morchy/pkg/controlplane/implementation/jsonformatter"
	"github.com/wernsiet/morchy/pkg/runtime"
)

func validateWorkloadID(workloadID string) error {
	if workloadID == "" {
		return domain.ErrorRequestParamsValidation.Wrapf(nil, "workload_id cannot be empty")
	}
	if _, err := uuid.Parse(workloadID); err != nil {
		return domain.ErrorRequestParamsValidation.Wrapf(err, "invalid workload_id format")
	}
	return nil
}

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
//	@Success		200		{array}		jsonformatter.WorkloadResponse
//	@Failure		400		{object}	map[string]string	"Invalid request parameters"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/api/v1/workloads [get]
func (rh *RouterHandler) listWorkloads(c *gin.Context) {
	var queryParams struct {
		StatusEq        *string `form:"status" example:"new"`
		CPU             *uint   `form:"cpu" example:"100"`
		RAM             *uint   `form:"ram" example:"256"`
		SchedulableOnly *bool   `form:"schedulable_only" example:"true"`
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

	schedulableOnly := false
	if queryParams.SchedulableOnly != nil {
		schedulableOnly = *queryParams.SchedulableOnly
	}

	workloads, err := rh.ucHandler.ListWorkloads(c, queryParams.StatusEq, resorceFilter, schedulableOnly)
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

// getWorkload godoc
//
//	@Summary		Get workload
//	@Description	Retrieve a single workload by ID
//	@Tags			workloads
//	@Accept			json
//	@Produce		json
//	@Param			workload_id	path		string	true	"Workload ID"	minlength(1)
//	@Success		200			{object}	jsonformatter.WorkloadResponse
//	@Failure		400			{object}	map[string]string	"Invalid request parameters"
//	@Failure		500			{object}	map[string]string	"Internal server error"
//	@Router			/api/v1/workloads/{workload_id} [get]
func (rh *RouterHandler) getWorkload(c *gin.Context) {
	workloadID := c.Param("workload_id")
	if err := validateWorkloadID(workloadID); err != nil {
		handleError(c, err)
		return
	}

	workload, err := rh.ucHandler.GetWorkload(c, workloadID)
	if err != nil {
		handleError(c, err)
		return
	}
	workloadApiModel := jsonformatter.NewWorkloadResponseFromDomain(workload)
	c.JSON(http.StatusOK, workloadApiModel)
}

// createWorkload godoc
//
//	@Summary		Create workload
//	@Description	Create a new workload from provided spec
//	@Tags			workloads
//	@Accept			json
//	@Produce		json
//	@Param			workloadSpec	body		jsonformatter.WorkloadSpecRequest	true	"Workload specification"
//	@Success		201				{object}	jsonformatter.WorkloadResponse
//	@Failure		400				{object}	map[string]string	"Invalid request parameters"
//	@Failure		500				{object}	map[string]string	"Internal server error"
//	@Router			/api/v1/workloads [post]
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

// deleteWorkload godoc
//
//	@Summary		Delete workload
//	@Description	Delete a workload and all related models (lease, spec)
//	@Tags			workloads
//	@Accept			json
//	@Produce		json
//	@Param			workload_id	path	string	true	"Workload ID"	minlength(1)
//	@Success		204
//	@Failure		400	{object}	map[string]string	"Invalid request parameters"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/api/v1/workloads/{workload_id} [delete]
func (rh *RouterHandler) deleteWorkload(c *gin.Context) {
	workloadID := c.Param("workload_id")
	if err := validateWorkloadID(workloadID); err != nil {
		handleError(c, err)
		return
	}

	if err := rh.ucHandler.DeleteWorkload(c, workloadID); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
