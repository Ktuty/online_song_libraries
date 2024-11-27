package main

import (
	"context"
	"os"
	"os/signal"
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

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
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
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repo := repository.NewRepository(db)
	service := services.NewService(repo)
	handler := handlers.NewHandler(service)

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
