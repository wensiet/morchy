package domain

import "github.com/samber/oops"

var (
	errorBaseNotFound            = oops.Code(SCodeNotFound)
	errorBaseInternalServerError = oops.Code(SInternalServerError)

	ErrorWorkloadNotFound     = errorBaseNotFound.With(SDomain, SWorkload)
	ErrorBaseWorkloadInternal = errorBaseInternalServerError.With(SDomain, SWorkload)

	ErrorWorkloadTerminatedOnControlPlane = oops.Code(STerminatedOnControlPlane)
	ErrorWorkloadOwnedByAnotherNode       = oops.Code(SOwnedByAnotherNode)

	ErrorInsufficientResources = oops.Code(SInsufficientResources)

	ErrorBaseWorkloadHealthcheckFailed = oops.Code(SHealthcheckFailed)
)
