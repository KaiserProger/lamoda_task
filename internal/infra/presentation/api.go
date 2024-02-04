package presentation

import (
	"context"
	"encoding/json"
	"lamoda_task/internal/app/services"
	"net/http"

	"go.uber.org/zap"
)

type ReserveRequest struct {
	ItemCodes []int `json:"item_codes"`
}

type FreeReserveRequest = ReserveRequest

type GetWarehouseItems struct {
}

type _apiHandler struct {
	svc    services.ItemService
	logger *zap.Logger
}

func NewApiHandler(svc services.ItemService) *_apiHandler {
	return &_apiHandler{
		svc: svc,
	}
}

func (handler *_apiHandler) MakeReservation(w http.ResponseWriter, req *http.Request) {
	var request ReserveRequest

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		handler.logger.Error("decode request fail", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := handler.svc.MakeReservation(context.Background(), request.ItemCodes); err != nil {
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := handler.svc.FreeReservation(context.Background(), request.ItemCodes); err != nil {
		handler.logger.Error("service execution fail", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
