package status_test

import (
	"context"
	"testing"

	"github.com/ficusinapot/ds/internal/domain/status"
	statuscontracts "github.com/ficusinapot/ds/internal/models/coreifc/status"
	dbifcmocks "github.com/ficusinapot/ds/internal/models/dbifc/mocks"
	"github.com/ficusinapot/ds/internal/models/entities"

	"github.com/joomcode/errorx"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestStatusUseCaseGetStatusWhenDatabaseIsAvailable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	provider := dbifcmocks.NewMockStatusProvider(gomock.NewController(t))
	provider.EXPECT().IsConnectionAvailable(gomock.Any()).Return(true)

	status, err := status.NewStatusUseCase(provider).GetStatus(ctx)

	require.NoError(t, err)
	require.Equal(t, &entities.StatusData{
		DBAvailable:     true,
		OperatingStatus: entities.FullyOperational,
	}, status)
}

func TestStatusUseCaseGetStatusWhenDatabaseIsUnavailable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	provider := dbifcmocks.NewMockStatusProvider(gomock.NewController(t))
	provider.EXPECT().IsConnectionAvailable(gomock.Any()).Return(false)

	status, err := status.NewStatusUseCase(provider).GetStatus(ctx)

	require.NoError(t, err)
	require.Equal(t, &entities.StatusData{
		DBAvailable:     false,
		OperatingStatus: entities.NonOperational,
	}, status)
}

func TestStatusUseCaseHealthz(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	provider := dbifcmocks.NewMockStatusProvider(gomock.NewController(t))
	provider.EXPECT().IsConnectionAvailable(gomock.Any()).Return(true)

	err := status.NewStatusUseCase(provider).Healthz(ctx)

	require.NoError(t, err)
}

func TestStatusUseCaseHealthzReturnsConnectionError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	provider := dbifcmocks.NewMockStatusProvider(gomock.NewController(t))
	provider.EXPECT().IsConnectionAvailable(gomock.Any()).Return(false)

	err := status.NewStatusUseCase(provider).Healthz(ctx)

	require.True(t, errorx.IsOfType(err, statuscontracts.ErrConnectionNotAvailable), "unexpected error: %v", err)
}

func TestStatusUseCaseReadyz(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	provider := dbifcmocks.NewMockStatusProvider(gomock.NewController(t))
	provider.EXPECT().IsConnectionAvailable(gomock.Any()).Return(true)

	err := status.NewStatusUseCase(provider).Readyz(ctx)

	require.NoError(t, err)
}

func TestStatusUseCaseReadyzReturnsConnectionError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	provider := dbifcmocks.NewMockStatusProvider(gomock.NewController(t))
	provider.EXPECT().IsConnectionAvailable(gomock.Any()).Return(false)

	err := status.NewStatusUseCase(provider).Readyz(ctx)

	require.True(t, errorx.IsOfType(err, statuscontracts.ErrConnectionNotAvailable), "unexpected error: %v", err)
}
