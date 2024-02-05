package services_test

import (
	"context"
	"lamoda_task/internal/app/persistence/repositories"
	"lamoda_task/internal/app/services"
	"testing"

	appErrors "lamoda_task/internal/app/errors"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type test struct {
	name   string
	data   any
	mock   func()
	result interface{}
	err    error
}

func fixture(t *testing.T) (services.ItemService, *MockTransactional, *MockItemRepository, *MockWarehouseRepository, *MockReserveRepository) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	txManager := NewMockTransactional(ctrl)
	itemRepo := NewMockItemRepository(ctrl)
	whRepo := NewMockWarehouseRepository(ctrl)
	reserveRepo := NewMockReserveRepository(ctrl)
	service := services.NewItemService(
		txManager,
		itemRepo,
		whRepo,
		reserveRepo)

	return service, txManager, itemRepo, whRepo, reserveRepo
}

func TestItemService_MakeReservation(t *testing.T) {
	t.Parallel()

	svc, txManager, itemRepo, _, reserveRepo := fixture(t)

	tests := []test{
		{
			name: "empty stock",
			data: []int{1, 2, 3, 4},
			mock: func() {
				txManager.EXPECT().WithinTransaction(context.Background(), gomock.Any()).DoAndReturn(func(ctx context.Context, f func(ctx context.Context) error) error {
					return f(ctx)
				}).AnyTimes()
				itemRepo.EXPECT().GetStoredAt(context.Background(), []int{1, 2, 3, 4}).Return(nil, nil)
			},
			result: nil,
			err:    appErrors.ErrNotFound,
		},
		{
			name: "success",
			data: []int{1, 1, 1, 2, 2},
			mock: func() {
				orders := []*repositories.StoredItem{
					{
						ItemCode:    1,
						WarehouseId: 1,
						Quantity:    4,
					},
					{
						ItemCode:    2,
						WarehouseId: 1,
						Quantity:    5,
					},
				}
				txManager.EXPECT().WithinTransaction(context.Background(), gomock.Any()).DoAndReturn(func(ctx context.Context, f func(ctx context.Context) error) error {
					return f(ctx)
				}).AnyTimes()
				itemRepo.EXPECT().GetStoredAt(context.Background(), []int{1, 2}).Return(orders, nil).AnyTimes()
				reserveRepo.EXPECT().MakeReservation(context.Background(), orders).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.mock()

			data, ok := tt.data.([]int)
			require.True(t, ok)

			err := svc.MakeReservation(context.Background(), data)
			require.Equal(t, tt.err, err)
		})
	}

}
