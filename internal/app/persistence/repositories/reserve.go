package repositories

import "context"

type ReserveRepository interface {
	// Make reservation.
	MakeReservation(ctx context.Context, orders []*StoredItem) error
	// Free items from the latest reservation.
	FreeReservation(ctx context.Context, itemsCount map[int]int) error
}
