package main

import (
	"database/sql"
	"errors"
	"lamoda_task/internal/app/services"
	"lamoda_task/internal/infra/persistence/postgres"
	"lamoda_task/internal/infra/presentation"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(errors.Join(errors.New("logger did not initialize"), err))
	}

	logger.Info("logger initialized")

	driverName, dbName, httpHost := "pgx", os.Getenv("DATABASE_URL"), os.Getenv("HOST_ADDRESS")

	db, err := sql.Open(driverName, dbName)
	if err != nil {
		logger.Fatal("database did not initialize", zap.Error(err))
	}

	logger.Info("database initialized")

	router := mux.NewRouter()
	logger.Info("router initialized")

	txManager := postgres.NewTransactional(db)

	itemRepo := postgres.NewItemRepository()
	reserveRepo := postgres.NewReserveRepository()
	warehouseRepo := postgres.NewWarehouseRepository()

	itemService := services.NewItemService(txManager, itemRepo, warehouseRepo, reserveRepo)

	handler := presentation.NewApiHandler(itemService, logger)
	logger.Info("all dependencies initialized")

	handler.RegisterHandle(router)
	logger.Info("all routes registered")

	logger.Info("starting server...")
	if err := http.ListenAndServe(httpHost, router); err != nil {
		logger.Fatal("critical server error", zap.Error(err))
	}
}
