package presentation

import (
	"context"
	"encoding/json"
	"lamoda_task/internal/app/services"
	"net/http"
	"strconv"

	appErrors "lamoda_task/internal/app/errors"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ReserveRequest struct {
	ItemCodes []int `json:"item_codes"`
}

type FreeReserveRequest = ReserveRequest

type ItemResponse struct {
	Code     int    `json:"code"`
	Name     string `json:"name"`
	Size     int    `json:"size"`
	Quantity int    `json:"quantity"`
}

type WarehouseRequest struct {
	Id int `json:"warehouse_id"`
}

type WarehouseResponse struct {
	Name       string          `json:"name"`
	Accessible bool            `json:"accessible"`
	Items      []*ItemResponse `json:"items"`
}

type _apiHandler struct {
	svc    services.ItemService
	logger *zap.Logger
}

func NewApiHandler(svc services.ItemService, logger *zap.Logger) *_apiHandler {
	return &_apiHandler{
		svc:    svc,
		logger: logger,
	}
}

func (handler *_apiHandler) MakeReservation(w http.ResponseWriter, req *http.Request) {
	var request ReserveRequest

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		handler.logger.Error("decode request fail", zap.Error(err))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := handler.svc.MakeReservation(context.Background(), request.ItemCodes); err != nil {
		if err == appErrors.ErrNotFound {
			handler.logger.Error("not found", zap.Error(err))
			w.WriteHeader(http.StatusNotFound)
			return
		}
		handler.logger.Error("service execution fail", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *_apiHandler) FreeReservation(w http.ResponseWriter, req *http.Request) {
	var request FreeReserveRequest

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		handler.logger.Error("decode request fail", zap.Error(err))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := handler.svc.FreeReservation(context.Background(), request.ItemCodes); err != nil {
		if err == appErrors.ErrNotFound {
			handler.logger.Error("not found", zap.Error(err))
			w.WriteHeader(http.StatusNotFound)
			return
		}
		handler.logger.Error("service execution fail", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *_apiHandler) GetWarehouse(w http.ResponseWriter, req *http.Request) {
	warehouseParamId := req.URL.Query().Get("id")

	warehouseId, err := strconv.Atoi(warehouseParamId)
	if err != nil {
		handler.logger.Error("decode warehouse id param fail", zap.Error(err))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	request := &WarehouseRequest{
		Id: warehouseId,
	}

	warehouse, err := handler.svc.Warehouse(context.Background(), request.Id)
	if err != nil {
		handler.logger.Error("service execution fail", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if warehouse == nil {
		handler.logger.Error("warehouse not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := WarehouseResponse{
		Name:       warehouse.Name,
		Accessible: warehouse.Accessible,
		Items:      []*ItemResponse{},
	}

	for _, item := range warehouse.Items {
		_item := ItemResponse(*item)
		response.Items = append(response.Items, &_item)
	}

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		handler.logger.Error("encode response fail", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (handler *_apiHandler) RegisterHandle(_mux *mux.Router) {
	_mux.HandleFunc("/reserve", handler.MakeReservation).Methods("POST")
	_mux.HandleFunc("/reserve", handler.FreeReservation).Methods("PATCH")
	_mux.HandleFunc("/warehouses", handler.GetWarehouse).Methods("GET")
}
