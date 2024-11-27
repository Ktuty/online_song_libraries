package handlers

import (
	"github.com/Ktuty/internal/services"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
)

// Структура обработчика
type Handler struct {
	services *services.Service
	router   *mux.Router
}

// Функция для создания нового обработчика с заданными сервисами
func NewHandler(services *services.Service) *Handler {
	return &Handler{services: services}
}

// Функция для инициализации маршрутов
func (h *Handler) InitRouts() *mux.Router {
	h.router = mux.NewRouter()
	h.endpoints()

	return h.router
}

// Функция для настройки конечных точек маршрутизатора
func (h *Handler) endpoints() {
	h.router.HandleFunc("/songs", h.Songs).Methods(http.MethodGet)
	h.router.HandleFunc("/songs", h.NewSong).Methods(http.MethodPost)
	h.router.HandleFunc("/songs/{id}", h.SongByID).Methods(http.MethodGet)
	h.router.HandleFunc("/songs/{id}", h.UpdateSong).Methods(http.MethodPatch)
	h.router.HandleFunc("/songs/{id}", h.DeleteSongs).Methods(http.MethodDelete)

	// Swagger маршрут
	h.router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}
