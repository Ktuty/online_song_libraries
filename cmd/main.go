package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_ "github.com/Ktuty/docs"
	"github.com/Ktuty/internal/handlers"
	"github.com/Ktuty/internal/repository"
	"github.com/Ktuty/internal/services"
	"github.com/Ktuty/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	_ "github.com/swaggo/http-swagger"
)

//	@title			Online Songs-lib
//	@version		1.0
//	@description	This is a sample server Online Songs-lib server.
//	@host			localhost:8080

func main() {
	// Загрузка переменных окружения из файла .env
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	// Установка уровня логирования из переменной окружения
	logLevel := os.Getenv("LOG_LEVEL")
	level, err := logrus.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		logrus.Fatalf("error parsing log level: %s", err.Error())
	}
	logrus.SetLevel(level)

	// Установка формата логирования
	logrus.SetFormatter(&logrus.JSONFormatter{})

	db, err := repository.NewPostgres(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repo := repository.NewRepository(db)
	service := services.NewService(repo)
	handler := handlers.NewHandler(service, os.Getenv("URL"))

	srv := new(server.Server)
	go func() {
		if err := srv.Run(os.Getenv("port"), handler.InitRouts()); err != nil {
			logrus.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()

	logrus.Printf("Server Started on port %s", os.Getenv("port"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Printf("Server Shutdown")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("error server Shutdown Failed: %s", err.Error())
	}

	db.Close()
}
