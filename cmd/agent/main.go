package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/wernsiet/morchy/pkg/agent/usecase"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	httpClient := &http.Client{}
	handler := usecase.NewHandler(httpClient, "http://localhost:8080")
	err := handler.ApplyWorkloadJoin(ctx)
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		err := handler.ReconcileWorkloads(ctx)
		if err != nil {
			fmt.Printf("Got error: %v\n", err)
		}
	}
}
