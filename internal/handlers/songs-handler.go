package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Ktuty/internal/models"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

//	@Summary		Get all songs
//	@Description	Get a list of songs with optional filtering and pagination
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Page number"	default(1)
//	@Param			pageSize	query		int		false	"Page size"		default(10)
//	@Param			song		query		string	false	"Song name"
//	@Param			group		query		string	false	"Group name"
//	@Param			text		query		string	false	"Song text"
//	@Param			releaseDate	query		string	false	"Release date"
//	@Param			link		query		string	false	"Link"
//	@Success		200			{object}	models.Songs
//	@Router			/songs [get]
func (h *Handler) Songs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получение номера страницы и размера страницы из параметров запроса
	page := getQueryParamAsInt(r, "page", 1)
	pageSize := getQueryParamAsInt(r, "pageSize", 10)
	logrus.Debugf("Songs: page=%d, pageSize=%d", page, pageSize)

	// Создание фильтра для поиска песен
	filter := models.Songs{
		Song:        r.URL.Query().Get("song"),
		Group:       r.URL.Query().Get("group"),
		Text:        r.URL.Query().Get("text"),
		ReleaseDate: r.URL.Query().Get("releaseDate"),
		Link:        r.URL.Query().Get("link"),
	}
	logrus.Debugf("Songs: filter=%+v", filter)

	// Получение списка песен с использованием сервиса
	songs, totalPages, err := h.services.GetAll(filter, page, pageSize)
	if err != nil {
		logrus.Printf("Ошибка при получении песен: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Формирование ответа
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
	logrus.Debugf("Songs: response=%+v", response)

	// Кодирование ответа в JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.Printf("Ошибка при кодировании ответа: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//	@Summary		Create a new song
//	@Description	Create a new song
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			song	body	models.Songs	true	"Song details"
//	@Success		201
//	@Router			/songs [post]
func (h *Handler) NewSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var song models.Songs
	defer r.Body.Close()

	// Декодирование JSON из тела запроса
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		logrus.Printf("Ошибка при декодировании JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logrus.Printf("Новая песня: %v дата: %#v", song, song.ReleaseDate)

	// Создание новой песни с использованием сервиса
	if err := h.services.Create(song); err != nil {
		logrus.Printf("Ошибка при создании песни: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

//	@Summary		Get a song by ID
//	@Description	Get a song by its ID
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	true	"Song ID"
//	@Param			vers	query		int	false	"Verse number"	default(0)
//	@Success		200		{object}	models.Songs
//	@Router			/songs/{id} [get]
func (h *Handler) SongByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получение номера куплета из параметров запроса
	vers := getQueryParamAsInt(r, "vers", 0)
	logrus.Debugf("SongByID: vers=%d", vers)

	// Получение ID песни из переменных маршрута
	vars := mux.Vars(r)
	songID, err := strconv.Atoi(vars["id"])
	if err != nil {
		logrus.Printf("Ошибка при преобразовании songID в int: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Debugf("SongByID: songID=%d", songID)

	// Получение песни по ID с использованием сервиса
	song, err := h.services.GetByID(songID)
	if err != nil {
		logrus.Printf("Ошибка при получении песни: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Разделение текста песни на куплеты, если указан номер куплета
	if vers != 0 {
		verse, err := splitTextDoubleMargins(song.Text, vers-1)
		if err != nil {
			logrus.Printf("Ошибка при разделении текста на куплеты: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		song.Text = verse
	}

	// Кодирование ответа в JSON
	if err := json.NewEncoder(w).Encode(song); err != nil {
		logrus.Printf("Ошибка при кодировании ответа: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//	@Summary		Update a song
//	@Description	Update a song by its ID
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int				true	"Song ID"
//	@Param			song	body	models.Songs	true	"Song details"
//	@Success		200
//	@Router			/songs/{id} [patch]
func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получение ID песни из переменных маршрута
	vars := mux.Vars(r)
	songID, err := strconv.Atoi(vars["id"])
	if err != nil {
		logrus.Printf("Ошибка при преобразовании songID в int: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Debugf("UpdateSong: songID=%d", songID)

	var song models.Songs
	defer r.Body.Close()

	// Декодирование JSON из тела запроса
	err = json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		logrus.Printf("Ошибка при декодировании JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновление песни с использованием сервиса
	if err := h.services.Update(songID, song); err != nil {
		logrus.Printf("Ошибка при обновлении песни: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

//	@Summary		Delete a song
//	@Description	Delete a song by its ID
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Song ID"
//	@Success		200
//	@Router			/songs/{id} [delete]
func (h *Handler) DeleteSongs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получение ID песни из переменных маршрута
	vars := mux.Vars(r)
	songID, err := strconv.Atoi(vars["id"])
	if err != nil {
		logrus.Printf("Ошибка при преобразовании songID в int: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Debugf("DeleteSongs: songID=%d", songID)

	// Удаление песни с использованием сервиса
	if err := h.services.Delete(songID); err != nil {
		logrus.Printf("Ошибка при удалении песни: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Функция для получения параметра из запроса в виде целого числа
func getQueryParamAsInt(r *http.Request, param string, defaultValue int) int {
	valueStr := r.URL.Query().Get(param)
	logrus.Debugf("getQueryParamAsInt: param=%s, valueStr=%s", param, valueStr)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		logrus.Debugf("getQueryParamAsInt: error converting valueStr to int: %v", err)
		return defaultValue
	}
	return value
}

// Функция для разделения текста на куплеты
func splitTextIntoVerses(text string, vers int) (string, error) {
	logrus.Debugf("splitTextIntoVerses: text=%s, vers=%d", text, vers)
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
		return "", fmt.Errorf("неверный индекс куплета: %d", vers+1)
	}
	return strings.TrimSpace(verses[vers]), nil
}

// Функция для разделения текста на куплеты
func splitTextDoubleMargins(text string, vers int) (string, error) {
	logrus.Debugf("splitTextDoubleMargins: text=%s, vers=%d", text, vers)
	parts := strings.Split(text, "\n\n")

	// Проверка границ массива
	if vers < 0 || vers >= len(parts) {
		return "", fmt.Errorf("invalid verse index: %d", vers+1)
	}

	return parts[vers], nil
}
