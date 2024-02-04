package services

import (
	"context"
	"errors"
	"lamoda_task/internal/app/models"
	"lamoda_task/internal/app/persistence"
	"lamoda_task/internal/app/persistence/repositories"
)

type Item struct {
	Code     int
	Name     string
	Size     int
	Quantity int
}

type Warehouse struct {
	Id         int
	Name       string
	Accessible bool
	Items      []*Item
}

type ItemService interface {
	MakeReservation(ctx context.Context, itemCodes []int) error
	FreeReservation(ctx context.Context, itemCodes []int) error
	Warehouse(ctx context.Context, warehouseId int) (*Warehouse, error)
}

type _itemServiceImpl struct {
	txManager     persistence.Transactional
	itemRepo      repositories.ItemRepository
	warehouseRepo repositories.WarehouseRepository
	reserveRepo   repositories.ReserveRepository
}

func NewItemService(txManager persistence.Transactional,
	itemRepo repositories.ItemRepository,
	warehouseRepo repositories.WarehouseRepository,
	reserveRepo repositories.ReserveRepository) ItemService {
	return &_itemServiceImpl{
		txManager:     txManager,
		itemRepo:      itemRepo,
		warehouseRepo: warehouseRepo,
		reserveRepo:   reserveRepo,
	}
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
	reserveOrders := make([]*models.ReservationItem, 0)
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

			reserveOrders = append(reserveOrders, &models.ReservationItem{
				ItemCode:    storedItem.ItemCode,
				WarehouseId: storedItem.WarehouseId,
				Quantity:    finalQty,
			})
		}

		if err := svc.reserveRepo.MakeReservation(txCtx, reserveOrders); err != nil {
			return errors.Join(errors.New("make reservation fail"), err)
		}
		return nil
	})
}

func (svc *_itemServiceImpl) FreeReservation(ctx context.Context, itemCodes []int) error {
	countMap := svc.itemCodesAsUniqueMap(itemCodes)
	return svc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		reserveItems, err := svc.reserveRepo.GetReservation(ctx, countMap)
		if err != nil {
			return errors.Join(errors.New("get reservation items fail"), err)
		}

		err = svc.reserveRepo.FreeReservation(txCtx, reserveItems)
		if err != nil {
			return errors.Join(errors.New("free reservation fail"), err)
		}

		err = svc.warehouseRepo.UpdateStock(ctx, reserveItems)
		if err != nil {
			return errors.Join(errors.New("update stock fail"), err)
		}

		return nil
	})
}

func (svc *_itemServiceImpl) Warehouse(ctx context.Context, warehouseId int) (*Warehouse, error) {
	var warehouse *models.Warehouse
	err := svc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		_warehouse, err := svc.warehouseRepo.Get(txCtx, warehouseId)
		if err != nil {
			return errors.Join(errors.New("get warehouse fail"), err)
		}
		warehouse = _warehouse
		return nil
	})
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, nil
	}

	response := &Warehouse{
		Id:         warehouseId,
		Name:       warehouse.Name,
		Accessible: warehouse.Accessible,
		Items:      []*Item{},
	}

	for _, item := range warehouse.Items {
		response.Items = append(response.Items, &Item{
			Code:     item.Code,
			Name:     item.Name,
			Size:     item.Size,
			Quantity: item.Quantity,
		})
	}

	return response, err
}
