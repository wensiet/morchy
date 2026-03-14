package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/wernsiet/morchy/pkg/controlplane/domain"
	workloadMocks "github.com/wernsiet/morchy/pkg/controlplane/domain/workload/mocks"
	dbMocks "github.com/wernsiet/morchy/pkg/db.utils/mocks"
	"github.com/wernsiet/morchy/pkg/testutil"
)

func TestInteractor_GetLeaseByNodeAndWorkloadID(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		mockRepo := workloadMocks.NewMockRepository(t)
		mockDB := dbMocks.NewMockDB(t)
		mockFactory := workloadMocks.NewMockRepositoryFactory(t)

		logger := testutil.TestLogger()
		interactor := NewHandler(logger, mockRepo, mockFactory, mockDB, 30, 10, 300)

		expectedLease := testutil.NewTestLease("node-1", "workload-1")

		mockRepo.EXPECT().GetLease(mock.Anything, "node-1", "workload-1").
			Return(expectedLease, nil)

		result, err := interactor.GetLeaseByNodeAndWorkloadID(context.Background(), "node-1", "workload-1")

		require.NoError(t, err)
		require.Equal(t, "node-1", result.NodeID)
		require.Equal(t, "workload-1", result.WorkloadID)
	})

	t.Run("lease not found", func(t *testing.T) {
		mockRepo := workloadMocks.NewMockRepository(t)
		mockDB := dbMocks.NewMockDB(t)
		mockFactory := workloadMocks.NewMockRepositoryFactory(t)

		logger := testutil.TestLogger()
		interactor := NewHandler(logger, mockRepo, mockFactory, mockDB, 30, 10, 300)

		notFoundErr := domain.ErrorWorkloadRepositoryNotFound.New("lease not found")
		mockRepo.EXPECT().GetLease(mock.Anything, "node-1", "workload-1").
			Return(nil, notFoundErr)

		result, err := interactor.GetLeaseByNodeAndWorkloadID(context.Background(), "node-1", "workload-1")

		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := workloadMocks.NewMockRepository(t)
		mockDB := dbMocks.NewMockDB(t)
		mockFactory := workloadMocks.NewMockRepositoryFactory(t)

		logger := testutil.TestLogger()
		interactor := NewHandler(logger, mockRepo, mockFactory, mockDB, 30, 10, 300)

		mockRepo.EXPECT().GetLease(mock.Anything, "node-1", "workload-1").
			Return(nil, errors.New("database error"))

		result, err := interactor.GetLeaseByNodeAndWorkloadID(context.Background(), "node-1", "workload-1")

		require.Error(t, err)
		require.Nil(t, result)
	})
}

func TestInteractor_ExpireLeases(t *testing.T) {
	t.Run("successful expiration", func(t *testing.T) {
		mockRepo := workloadMocks.NewMockRepository(t)
		mockDB := dbMocks.NewMockDB(t)
		mockFactory := workloadMocks.NewMockRepositoryFactory(t)

		logger := testutil.TestLogger()
		interactor := NewHandler(logger, mockRepo, mockFactory, mockDB, 30, 10, 300)

		mockRepo.EXPECT().DeleteExpiredLeases(mock.Anything, 30).
			Return(nil)

		err := interactor.ExpireLeases(context.Background())

		require.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := workloadMocks.NewMockRepository(t)
		mockDB := dbMocks.NewMockDB(t)
		mockFactory := workloadMocks.NewMockRepositoryFactory(t)

		logger := testutil.TestLogger()
		interactor := NewHandler(logger, mockRepo, mockFactory, mockDB, 30, 10, 300)

		mockRepo.EXPECT().DeleteExpiredLeases(mock.Anything, 30).
			Return(errors.New("database error"))

		err := interactor.ExpireLeases(context.Background())

		require.Error(t, err)
	})
}

func TestInteractor_DeleteLease(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		mockRepo := workloadMocks.NewMockRepository(t)
		mockDB := dbMocks.NewMockDB(t)
		mockFactory := workloadMocks.NewMockRepositoryFactory(t)

		logger := testutil.TestLogger()
		interactor := NewHandler(logger, mockRepo, mockFactory, mockDB, 30, 10, 300)

		mockRepo.EXPECT().DeleteLease(mock.Anything, "node-1", "workload-1").
			Return(nil)

		err := interactor.DeleteLease(context.Background(), "node-1", "workload-1")

		require.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := workloadMocks.NewMockRepository(t)
		mockDB := dbMocks.NewMockDB(t)
		mockFactory := workloadMocks.NewMockRepositoryFactory(t)

		logger := testutil.TestLogger()
		interactor := NewHandler(logger, mockRepo, mockFactory, mockDB, 30, 10, 300)

		mockRepo.EXPECT().DeleteLease(mock.Anything, "node-1", "workload-1").
			Return(errors.New("database error"))

		err := interactor.DeleteLease(context.Background(), "node-1", "workload-1")

		require.Error(t, err)
	})
}
