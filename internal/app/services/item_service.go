package services

import (
	"context"
	"errors"
	"lamoda_task/internal/app/models"
	"lamoda_task/internal/app/persistence"
	"lamoda_task/internal/app/persistence/repositories"
)

type ItemService interface {
	MakeReservation(ctx context.Context, itemCodes []int) error
	FreeReservation(ctx context.Context, itemCodes []int) error
	GetItems(ctx context.Context, warehouseId int) ([]*models.Warehouse, error)
}

type _itemServiceImpl struct {
	txManager     persistence.Transactional
	itemRepo      repositories.ItemRepository
	warehouseRepo repositories.WarehouseRepository
	reverseRepo   repositories.ReserveRepository
}

func NewItemService() ItemService {
	return &_itemServiceImpl{}
}

func (svc *_itemServiceImpl) uniqueItemCodes(itemCodes []int) []int {
	uniqueCodesMap := map[int]bool{}
	uniqueCodes := make([]int, 0)

	for _, itemCode := range itemCodes {
		uniqueCodesMap[itemCode] = true
	}

	for key := range uniqueCodesMap {
		uniqueCodes = append(uniqueCodes, key)
	}

	return uniqueCodes
}

func (svc *_itemServiceImpl) itemCodesAsUniqueMap(itemCodes []int) map[int]int {
	countMap := map[int]int{}

	for _, itemCode := range itemCodes {
		_, exists := countMap[itemCode]
		if !exists {
			countMap[itemCode] = 0
		}
		countMap[itemCode] += 1
	}

	return countMap
}

func (svc *_itemServiceImpl) MakeReservation(ctx context.Context, itemCodes []int) error {
	uniqueCodes := svc.uniqueItemCodes(itemCodes)
	countMap := svc.itemCodesAsUniqueMap(itemCodes)
	reserveOrders := make([]*repositories.StoredItem, 0)
	return svc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		storedItems, err := svc.itemRepo.GetStoredAt(txCtx, uniqueCodes)
		if err != nil {
			return errors.Join(errors.New("get warehouses fail"), err)
		}

		for _, storedItem := range storedItems {
			qty := countMap[storedItem.ItemCode]
			if qty == 0 {
				continue
			}

			finalQty := min(qty, storedItem.Quantity)
			countMap[storedItem.ItemCode] -= finalQty

			reserveOrders = append(reserveOrders, &repositories.StoredItem{
				ItemCode:    storedItem.ItemCode,
				WarehouseId: storedItem.WarehouseId,
				Quantity:    finalQty,
			})
		}

		if err := svc.reverseRepo.MakeReservation(txCtx, reserveOrders); err != nil {
			return errors.Join(errors.New("make reservation fail"), err)
		}
		return nil
	})
}

func (svc *_itemServiceImpl) FreeReservation(ctx context.Context, itemCodes []int) error {
	return svc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		return svc.reverseRepo.FreeReservation(txCtx, itemCodes)
	})
}

func (svc *_itemServiceImpl) GetItems(ctx context.Context, warehouseId int) ([]*models.Warehouse, error) {
	var warehouses []*models.Warehouse
	err := svc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		_warehouses, err := svc.warehouseRepo.Get(warehouseId)
		if err != nil {
			return errors.Join(errors.New("get warehouse fail"), err)
		}
		warehouses = _warehouses
		return nil
	})
	return warehouses, err
}
