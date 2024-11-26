package main

import (
	"context"
	"github.com/Ktuty/internal/handlers"
	"github.com/Ktuty/internal/repository"
	"github.com/Ktuty/internal/services"
	"github.com/Ktuty/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgres(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	repo := repository.NewRepository(db)
	service := services.NewService(repo)
	handler := handlers.NewHandler(service)

	srv := new(server.Server)
	go func() {
		if err := srv.Run(os.Getenv("port"), handler.InitRouts()); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	log.Printf("Server Started on port %s", os.Getenv("port"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Printf("Server Shutdown")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("error server Shutdown Failed: %s", err.Error())
	}

	db.Close()
}
