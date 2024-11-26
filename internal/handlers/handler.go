package handlers

import (
	"github.com/Ktuty/internal/services"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	services *services.Service
	router   *mux.Router
}

func NewHandler(services *services.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRouts() *mux.Router {
	h.router = mux.NewRouter()
	h.endpoints()

	return h.router
}

func (h *Handler) endpoints() {
	h.router.HandleFunc("/songs", h.Songs).Methods(http.MethodGet)
	h.router.HandleFunc("/songs", h.NewSong).Methods(http.MethodPost)
	h.router.HandleFunc("/songs/{id}", h.SongByID).Methods(http.MethodGet)
	h.router.HandleFunc("/songs/{id}", h.UpdateSong).Methods(http.MethodPatch)
	h.router.HandleFunc("/songs/{id}", h.DeleteSongs).Methods(http.MethodDelete)

}
