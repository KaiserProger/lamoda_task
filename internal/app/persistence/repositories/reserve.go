package repositories

import (
	"context"
	"lamoda_task/internal/app/models"
)

type ReserveRepository interface {
	// Make reservation.
	MakeReservation(ctx context.Context, orders []*models.ReservationItem) error
	// Free items from the latest reservation.
	FreeReservation(ctx context.Context, reservation []*models.ReservationItem) error
	GetReservation(ctx context.Context, itemsCount map[int]int) ([]*models.ReservationItem, error)
}
