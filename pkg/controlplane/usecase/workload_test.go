package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/wernsiet/morchy/pkg/controlplane/domain"
	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
	workloadMocks "github.com/wernsiet/morchy/pkg/controlplane/domain/workload/mocks"
	dbMocks "github.com/wernsiet/morchy/pkg/db.utils/mocks"
	"github.com/wernsiet/morchy/pkg/runtime"
	"github.com/wernsiet/morchy/pkg/testutil"
)

func TestInteractor_CreateWorkload(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockRepo := workloadMocks.NewMockRepository(t)
		mockDB := dbMocks.NewMockDB(t)
		mockFactory := workloadMocks.NewMockRepositoryFactory(t)

		logger := testutil.TestLogger()
		interactor := NewHandler(logger, mockRepo, mockFactory, mockDB, 30, 10, 300)

		expectedWorkload := testutil.NewTestWorkload("test-id")
		mockRepo.EXPECT().CreateWorkload(mock.Anything, mock.MatchedBy(func(w workload.Workload) bool {
			return w.Spec.Name == "test-workload" && w.Status == workload.NewWorkloadStatus
		})).Return(expectedWorkload, nil)

		spec := testutil.NewTestWorkloadSpec()
		result, err := interactor.CreateWorkload(testutil.TestContext(), spec)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, "test-id", result.ID)
		require.Equal(t, workload.NewWorkloadStatus, result.Status)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := workloadMocks.NewMockRepository(t)
		mockDB := dbMocks.NewMockDB(t)
		mockFactory := workloadMocks.NewMockRepositoryFactory(t)

		logger := testutil.TestLogger()
		interactor := NewHandler(logger, mockRepo, mockFactory, mockDB, 30, 10, 300)

		mockRepo.EXPECT().CreateWorkload(mock.Anything, mock.AnythingOfType("workload.Workload")).
			Return(nil, errors.New("database error"))

		spec := testutil.NewTestWorkloadSpec()
		result, err := interactor.CreateWorkload(testutil.TestContext(), spec)

		require.Error(t, err)
		require.Nil(t, result)
	})
}

func TestInteractor_ListWorkloads(t *testing.T) {
	t.Run("successful listing", func(t *testing.T) {
		mockRepo := workloadMocks.NewMockRepository(t)
		mockDB := dbMocks.NewMockDB(t)
		mockFactory := workloadMocks.NewMockRepositoryFactory(t)

		logger := testutil.TestLogger()
		interactor := NewHandler(logger, mockRepo, mockFactory, mockDB, 30, 10, 300)

		expectedWorkloads := []*workload.Workload{
			testutil.NewTestWorkload("workload-1"),
			testutil.NewTestWorkload("workload-2"),
		}

		mockRepo.EXPECT().ListWorkloads(mock.Anything, (*string)(nil), (*runtime.ResourceLimits)(nil)).
			Return(expectedWorkloads, nil)

		result, err := interactor.ListWorkloads(context.Background(), nil, nil)

		require.NoError(t, err)
		require.Equal(t, 2, len(result))
		require.Equal(t, "workload-1", result[0].ID)
		require.Equal(t, "workload-2", result[1].ID)
	})

	t.Run("listing with status filter", func(t *testing.T) {
		mockRepo := workloadMocks.NewMockRepository(t)
		mockDB := dbMocks.NewMockDB(t)
		mockFactory := workloadMocks.NewMockRepositoryFactory(t)

		logger := testutil.TestLogger()
		interactor := NewHandler(logger, mockRepo, mockFactory, mockDB, 30, 10, 300)

		status := "active"
		expectedWorkloads := []*workload.Workload{
			testutil.NewTestWorkload("workload-1"),
		}

		mockRepo.EXPECT().ListWorkloads(mock.Anything, &status, (*runtime.ResourceLimits)(nil)).
			Return(expectedWorkloads, nil)

		result, err := interactor.ListWorkloads(context.Background(), &status, nil)

		require.NoError(t, err)
		require.Equal(t, 1, len(result))
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := workloadMocks.NewMockRepository(t)
		mockDB := dbMocks.NewMockDB(t)
		mockFactory := workloadMocks.NewMockRepositoryFactory(t)

		logger := testutil.TestLogger()
		interactor := NewHandler(logger, mockRepo, mockFactory, mockDB, 30, 10, 300)

		mockRepo.EXPECT().ListWorkloads(mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("database error"))

		result, err := interactor.ListWorkloads(context.Background(), nil, nil)

		require.Error(t, err)
		require.Nil(t, result)
	})
}

func TestInteractor_DeleteWorkload(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		mockRepo := workloadMocks.NewMockRepository(t)
		mockDB := dbMocks.NewMockDB(t)
		mockFactory := workloadMocks.NewMockRepositoryFactory(t)

		logger := testutil.TestLogger()
		interactor := NewHandler(logger, mockRepo, mockFactory, mockDB, 30, 10, 300)

		mockRepo.EXPECT().DeleteWorkload(mock.Anything, "workload-1").
			Return(nil)

		err := interactor.DeleteWorkload(context.Background(), "workload-1")

		require.NoError(t, err)
	})

	t.Run("workload not found", func(t *testing.T) {
		mockRepo := workloadMocks.NewMockRepository(t)
		mockDB := dbMocks.NewMockDB(t)
		mockFactory := workloadMocks.NewMockRepositoryFactory(t)

		logger := testutil.TestLogger()
		interactor := NewHandler(logger, mockRepo, mockFactory, mockDB, 30, 10, 300)

		notFoundErr := domain.ErrorWorkloadRepositoryNotFound.New("workload not found")
		mockRepo.EXPECT().DeleteWorkload(mock.Anything, "workload-1").
			Return(notFoundErr)

		err := interactor.DeleteWorkload(context.Background(), "workload-1")

		require.Error(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := workloadMocks.NewMockRepository(t)
		mockDB := dbMocks.NewMockDB(t)
		mockFactory := workloadMocks.NewMockRepositoryFactory(t)

		logger := testutil.TestLogger()
		interactor := NewHandler(logger, mockRepo, mockFactory, mockDB, 30, 10, 300)

		mockRepo.EXPECT().DeleteWorkload(mock.Anything, "workload-1").
			Return(errors.New("database error"))

		err := interactor.DeleteWorkload(context.Background(), "workload-1")

		require.Error(t, err)
	})
}
