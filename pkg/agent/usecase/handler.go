package usecase

import (
	"context"
	"net/http"
)

type JoinLogic interface {
	ApplyWorkloadJoin(context.Context) error
	ReconcileWorkloads(context.Context) error
}

type Handler interface {
	JoinLogic
}

type interactor struct {
	httpClient      *http.Client
	controlPlaneURL string
}

func NewHandler(
	httpClient *http.Client, controlPlaneURL string,
) Handler {
	return &interactor{
		httpClient:      httpClient,
		controlPlaneURL: controlPlaneURL,
	}
}
