package domain

import "github.com/samber/oops"

var (
	errorBaseInternalServerError = oops.Code(string(InternalServerError))
	errorBaseNotFound            = oops.Code(string(NotFound))

	ErrorUnknownServerError = errorBaseInternalServerError.With(SDomain, SUnknown)

	ErrorWorkloadRepositoryInternalError = errorBaseInternalServerError.With(SDomain, SWorkload)
	ErrorWorkloadRepositoryNotFound      = errorBaseNotFound.With(SDomain, SWorkload)
)
