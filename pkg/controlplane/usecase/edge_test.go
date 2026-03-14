package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
	workloadMocks "github.com/wernsiet/morchy/pkg/controlplane/domain/workload/mocks"
	dbMocks "github.com/wernsiet/morchy/pkg/db.utils/mocks"
	"github.com/wernsiet/morchy/pkg/testutil"
)

func TestInteractor_ListEdges(t *testing.T) {
	t.Run("successful listing", func(t *testing.T) {
		mockRepo := workloadMocks.NewMockRepository(t)
		mockDB := dbMocks.NewMockDB(t)
		mockFactory := workloadMocks.NewMockRepositoryFactory(t)

		logger := testutil.TestLogger()
		interactor := NewHandler(logger, mockRepo, mockFactory, mockDB, 30, 10, 300)

		expectedEdges := []*workload.Edge{
			testutil.NewTestEdge("http://localhost:8080", "/api/v1"),
			testutil.NewTestEdge("http://localhost:8081", "/api/v2"),
		}

		mockRepo.EXPECT().ListEdges(mock.Anything).
			Return(expectedEdges, nil)

		result, err := interactor.ListEdges(context.Background())

		require.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, "http://localhost:8080", result[0].UpstreamAddress)
		assert.Equal(t, "/api/v1", result[0].ProxyPath)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo := workloadMocks.NewMockRepository(t)
		mockDB := dbMocks.NewMockDB(t)
		mockFactory := workloadMocks.NewMockRepositoryFactory(t)

		logger := testutil.TestLogger()
		interactor := NewHandler(logger, mockRepo, mockFactory, mockDB, 30, 10, 300)

		mockRepo.EXPECT().ListEdges(mock.Anything).
			Return([]*workload.Edge{}, nil)

		result, err := interactor.ListEdges(context.Background())

		require.NoError(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := workloadMocks.NewMockRepository(t)
		mockDB := dbMocks.NewMockDB(t)
		mockFactory := workloadMocks.NewMockRepositoryFactory(t)

		logger := testutil.TestLogger()
		interactor := NewHandler(logger, mockRepo, mockFactory, mockDB, 30, 10, 300)

		mockRepo.EXPECT().ListEdges(mock.Anything).
			Return(nil, errors.New("database error"))

		result, err := interactor.ListEdges(context.Background())

		require.Error(t, err)
		assert.Nil(t, result)
	})
}
