package domain

import "github.com/samber/oops"

var (
	errorBaseInternalServerError = oops.Code(string(InternalServerError))
	errorBaseNotFound            = oops.Code(string(NotFound))
	errorBaseBadRequest          = oops.Code(string(BadRequest))

	ErrorUnknownServerError = errorBaseInternalServerError.With(SDomain, SUnknown)

	ErrorRequestParamsValidation = errorBaseBadRequest.With(SReason, SValidation)

	ErrorWorkloadRepositoryInternalError = errorBaseInternalServerError.With(SDomain, SWorkload)
	ErrorWorkloadRepositoryNotFound      = errorBaseNotFound.With(SDomain, SWorkload)
)
