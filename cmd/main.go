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
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logger, err := config.Build()
	if err != nil {
		panic(errors.Join(errors.New("logger did not initialize"), err))
	}
	defer logger.Sync()

	logger.Info("logger initialized")

	httpHost := os.Getenv("HOST_ADDRESS")

	driverName := "pgx"
	dbUser, dbPassword := os.Getenv("DBUSER"), os.Getenv("DBPASSWORD")
	dbHost, _port, dbName := os.Getenv("DBHOST"), os.Getenv("DBPORT"), os.Getenv("DBNAME")

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
