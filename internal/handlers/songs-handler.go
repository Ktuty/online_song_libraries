package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Ktuty/internal/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func (h *Handler) Songs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page := getQueryParamAsInt(r, "page", 1)
	pageSize := getQueryParamAsInt(r, "pageSize", 10)

	filter := models.Songs{
		Song:        r.URL.Query().Get("song"),
		Group:       r.URL.Query().Get("group"),
		Text:        r.URL.Query().Get("text"),
		ReleaseDate: r.URL.Query().Get("releaseDate"),
		Link:        r.URL.Query().Get("link"),
	}

	songs, totalPages, err := h.services.GetAll(filter, page, pageSize)
	if err != nil {
		log.Printf("Error getting songs: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Songs       []models.Songs `json:"songs"`
		TotalPages  int            `json:"totalPages"`
		CurrentPage int            `json:"currentPage"`
		PageSize    int            `json:"pageSize"`
	}{
		Songs:       songs,
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) NewSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var song models.Songs
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		log.Printf("Error decoding json: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("New song: %v date: %#v", song, song.ReleaseDate)

	if err := h.services.Create(song); err != nil {
		log.Printf("Error creating song: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) SongByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vers := getQueryParamAsInt(r, "vers", 0)

	vars := mux.Vars(r)
	songID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Error converting songID to int: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	song, err := h.services.GetByID(songID)
	if err != nil {
		log.Printf("Error getting song: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if vers != 0 {
		verse, err := splitTextIntoVerses(song.Text, vers-1)
		if err != nil {
			log.Printf("Error splitting text into verses: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		song.Text = verse
	}

	if err := json.NewEncoder(w).Encode(song); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	songID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Error converting songID to int: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var song models.Songs
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		log.Printf("Error decoding json: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.services.Update(songID, song); err != nil {
		log.Printf("Error updating song: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteSongs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	songID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Error converting songID to int: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.services.Delete(songID); err != nil {
		log.Printf("Error deleting song: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Получение параметра из запроса
func getQueryParamAsInt(r *http.Request, param string, defaultValue int) int {
	valueStr := r.URL.Query().Get(param)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func splitTextIntoVerses(text string, vers int) (string, error) {
	// Регулярное выражение для нахождения куплетов (тегов)
	re := regexp.MustCompile(`(?m)(\[[^]]+\])`)
	// Разделяем текст на части по тегам
	parts := re.Split(text, -1)

	// Создание среза для куплетов
	verses := make([]string, 0)
	for i := 0; i < len(parts); i++ {
		// Добавляем текущее разделение в срез, если оно не пустое
		verse := strings.TrimSpace(parts[i])
		if verse != "" {
			verses = append(verses, verse)

			// Если следующий элемент - это тег, добавляем его к куплету
			if i < len(parts)-1 {
				tag := strings.TrimSpace(re.FindString(text))
				if tag != "" {
					verses[len(verses)-1] = tag + "\n" + verses[len(verses)-1]
					// Удаляем использованный тег для следующей итерации
					text = strings.Replace(text, tag, "", 1)
				}
			}
		}
	}

	// Проверяем индекс куплета
	if vers < 0 || vers >= len(verses) {
		return "", fmt.Errorf("invalid verse index: %d", vers)
	}
	return strings.TrimSpace(verses[vers]), nil
}
