package server

import (
	"context"
	"net/http"
	"time"
)

// Структура Server, которая инкапсулирует HTTP-сервер
type Server struct {
	HttpServer *http.Server
}

// Метод для запуска HTTP-сервера на указанном порту с заданным обработчиком
func (server *Server) Run(port string, handler http.Handler) error {
	// Настройка HTTP-сервера
	server.HttpServer = &http.Server{
		Addr:           ":" + port,       // Адрес и порт для прослушивания
		Handler:        handler,          // Обработчик запросов
		MaxHeaderBytes: 1 << 20,          // Максимальный размер заголовков (1MB)
		ReadTimeout:    10 * time.Second, // Таймаут на чтение запроса
		WriteTimeout:   10 * time.Second, // Таймаут на запись ответа
	}

	// Запуск HTTP-сервера и прослушивание входящих запросов
	return server.HttpServer.ListenAndServe()
}

// Метод для корректного завершения работы HTTP-сервера
func (server *Server) Shutdown(ctx context.Context) error {
	// Завершение работы HTTP-сервера с использованием контекста
	return server.HttpServer.Shutdown(ctx)
}
