package domain

import "github.com/samber/oops"

var (
	errorBaseInternalServerError = oops.Code(string(InternalServerError))
	errorBaseNotFound            = oops.Code(string(NotFound))
	errorBaseBadRequest          = oops.Code(string(BadRequest))
	errorBaseConflict            = oops.Code(string(Conflict))

	ErrorUnknownServerError = errorBaseInternalServerError.With(SDomain, SUnknown)

	ErrorRequestParamsValidation = errorBaseBadRequest.With(SReason, SValidation)

	ErrorWorkloadRepositoryInternalError = errorBaseInternalServerError.With(SDomain, SWorkload)
	ErrorWorkloadRepositoryNotFound      = errorBaseNotFound.With(SDomain, SWorkload)

	ErrorWorkloadLeaseOwnedByAnotherNode = errorBaseConflict.With(SDomain, SWorkload).With(SReason, SOwnedByAnotherNode)
)
