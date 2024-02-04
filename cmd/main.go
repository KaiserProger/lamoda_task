package main

import (
	"database/sql"
	"errors"
	"fmt"
	"lamoda_task/internal/app/services"
	"lamoda_task/internal/infra/persistence/postgres"
	"lamoda_task/internal/infra/presentation"
	"net/http"
	"os"
	"strconv"

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

	httpHost := os.Getenv("HOST_ADDRESS")

	driverName := "pgx"
	dbUser, dbPassword := os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD")
	dbHost, _port, dbName := os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB")

	dbPort, err := strconv.Atoi(_port)
	if err != nil {
		logger.Fatal("cannot receive database port or is not a number", zap.Error(err))
	}

	logger.Info("configuration read")

	database_dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open(driverName, database_dsn)
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
