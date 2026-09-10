package usecases_test

import (
	"context"
	"testing"

	"github.com/ficusinapot/ds/internal/domain/status/usecases"
	statuscontracts "github.com/ficusinapot/ds/internal/models/coreifc/status"
	"github.com/ficusinapot/ds/internal/models/dbifc"
	"github.com/ficusinapot/ds/internal/models/entities"

	"github.com/joomcode/errorx"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStatusUseCaseGetStatusWhenDatabaseIsAvailable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	provider := new(statusProviderMock)
	provider.On("IsConnectionAvailable", mock.Anything).Return(true).Once()

	status, err := usecases.NewStatusUseCase(provider).GetStatus(ctx)

	require.NoError(t, err)
	require.Equal(t, &entities.StatusData{
		DBAvailable:     true,
		OperatingStatus: entities.FullyOperational,
	}, status)
	provider.AssertExpectations(t)
}

func TestStatusUseCaseGetStatusWhenDatabaseIsUnavailable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	provider := new(statusProviderMock)
	provider.On("IsConnectionAvailable", mock.Anything).Return(false).Once()

	status, err := usecases.NewStatusUseCase(provider).GetStatus(ctx)

	require.NoError(t, err)
	require.Equal(t, &entities.StatusData{
		DBAvailable:     false,
		OperatingStatus: entities.NonOperational,
	}, status)
	provider.AssertExpectations(t)
}

func TestStatusUseCaseHealthz(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	provider := new(statusProviderMock)
	provider.On("IsConnectionAvailable", mock.Anything).Return(true).Once()

	err := usecases.NewStatusUseCase(provider).Healthz(ctx)

	require.NoError(t, err)
	provider.AssertExpectations(t)
}

func TestStatusUseCaseHealthzReturnsConnectionError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	provider := new(statusProviderMock)
	provider.On("IsConnectionAvailable", mock.Anything).Return(false).Once()

	err := usecases.NewStatusUseCase(provider).Healthz(ctx)

	require.True(t, errorx.IsOfType(err, statuscontracts.ErrConnectionNotAvailable), "unexpected error: %v", err)
	provider.AssertExpectations(t)
}

func TestStatusUseCaseReadyz(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	provider := new(statusProviderMock)
	provider.On("IsConnectionAvailable", mock.Anything).Return(true).Once()

	err := usecases.NewStatusUseCase(provider).Readyz(ctx)

	require.NoError(t, err)
	provider.AssertExpectations(t)
}

func TestStatusUseCaseReadyzReturnsConnectionError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	provider := new(statusProviderMock)
	provider.On("IsConnectionAvailable", mock.Anything).Return(false).Once()

	err := usecases.NewStatusUseCase(provider).Readyz(ctx)

	require.True(t, errorx.IsOfType(err, statuscontracts.ErrConnectionNotAvailable), "unexpected error: %v", err)
	provider.AssertExpectations(t)
}

type statusProviderMock struct {
	mock.Mock
}

func (m *statusProviderMock) GetDBStatus(ctx context.Context) dbifc.ConnectionStatus {
	args := m.Called(ctx)
	status, _ := args.Get(0).(dbifc.ConnectionStatus)

	return status
}

func (m *statusProviderMock) IsConnectionAvailable(ctx context.Context) bool {
	args := m.Called(ctx)

	return args.Bool(0)
}
