package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/wernsiet/morchy/pkg/testutil"
	workloadMocks "github.com/wernsiet/morchy/pkg/controlplane/domain/workload/mocks"
	dbMocks "github.com/wernsiet/morchy/pkg/db.utils/mocks"
)

func TestNewEvent(t *testing.T) {
	t.Run("event with valid ID", func(t *testing.T) {
		nodeID := "node-1"
		payload := []byte(`{"status":"success"}`)

		event := newEvent(nodeID, payload)

		require.NotEmpty(t, event.ID)
		require.Equal(t, event.ID, event.SourceID)
		require.Equal(t, nodeID, event.NodeID)
		require.Equal(t, json.RawMessage(payload), event.Payload)
	})

	t.Run("event IDs are unique", func(t *testing.T) {
		nodeID := "node-1"
		payload := []byte(`{}`)

		event1 := newEvent(nodeID, payload)
		event2 := newEvent(nodeID, payload)

		require.NotEqual(t, event1.ID, event2.ID)
	})
}

func TestInteractor_PushEvent(t *testing.T) {
	t.Run("successful event push", func(t *testing.T) {
		mockRepo := workloadMocks.NewMockRepository(t)
		mockDB := dbMocks.NewMockDB(t)
		mockFactory := workloadMocks.NewMockRepositoryFactory(t)

		logger := testutil.TestLogger()
		interactor := NewHandler(logger, mockRepo, mockFactory, mockDB, 30, 10, 300)

		testEvent := testutil.NewTestEvent("event-1", "source-1", "node-1", []byte(`{"status":"success"}`))

		mockRepo.EXPECT().SaveEvent(mock.Anything, *testEvent).
			Return(nil)

		err := interactor.PushEvent(context.Background(), *testEvent)

		require.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := workloadMocks.NewMockRepository(t)
		mockDB := dbMocks.NewMockDB(t)
		mockFactory := workloadMocks.NewMockRepositoryFactory(t)

		logger := testutil.TestLogger()
		interactor := NewHandler(logger, mockRepo, mockFactory, mockDB, 30, 10, 300)

		testEvent := testutil.NewTestEvent("event-1", "source-1", "node-1", []byte(`{}`))

		mockRepo.EXPECT().SaveEvent(mock.Anything, *testEvent).
			Return(errors.New("database error"))

		err := interactor.PushEvent(context.Background(), *testEvent)

		require.Error(t, err)
	})
}
