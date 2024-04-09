package main

import (
	"os"
	"time"

	"github.com/dkshi/bwgtest"
	_ "github.com/dkshi/bwgtest/docs"
	"github.com/dkshi/bwgtest/internal/handler"
	"github.com/dkshi/bwgtest/internal/repository"
	"github.com/dkshi/bwgtest/internal/service"
	"github.com/sirupsen/logrus"
)

// @title REST API test service for BWG
// @version 1.0
// @description API Сервер для управления котировками

// @host localhost:8080
// @BasePath /

func main() {
	// if err := godotenv.Load(); err != nil {
	// 	logrus.Fatalf("error loading .env file: %s", err.Error())
	// }

	if err := loadTimeLocation(os.Getenv("TIME_LOCATION")); err != nil {
		logrus.Fatalf("error loading time location: %s", err.Error())
	}

	connConfig := repository.Config{
		Host:     os.Getenv("DB_HOST"),
		DBName:   os.Getenv("DB_NAME"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		Port:     os.Getenv("DB_PORT"),
	}

	db, err := repository.NewPostgresDB(connConfig)
	if err != nil {
		logrus.Fatalf("error connecting to database: %s", err.Error())
	}

	lis, err := repository.NewListener(connConfig, "new_quotation_update")
	if err != nil {
		logrus.Fatalf("error creating listener: %s", err.Error())
	}

	repo := repository.NewRepository(db, lis)
	if err = repo.InitSchema("migrations/create_all.up.sql"); err != nil {
		logrus.Fatalf("error initializing db schema: %s", err.Error())
	}

	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	srv := new(bwgtest.Server)
	if err = srv.Run(os.Getenv("PORT"), handler.InitRoutes()); err != nil {
		logrus.Fatalf("error while running app: %s", err.Error())
	}

}

// Loads time zone, UTC prefered because it's postgres default timezone
func loadTimeLocation(location string) error {
	timeLocation, err := time.LoadLocation(location)
	if err != nil {
		return err
	}
	time.Local = timeLocation
	return nil
}
