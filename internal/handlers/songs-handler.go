package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Ktuty/internal/models"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
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
//	@Failure		400			{object}	models.ErrorResponse
//	@Failure		500			{object}	models.ErrorResponse
//	@Router			/songs [get]
func (h *Handler) Songs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получение номера страницы и размера страницы из параметров запроса
	page := getQueryParamAsInt(r, "page", 1)
	pageSize := getQueryParamAsInt(r, "pageSize", 10)
	logrus.WithFields(logrus.Fields{
		"page":     page,
		"pageSize": pageSize,
	}).Info("Songs: page and pageSize parameters")

	// Создание фильтра для поиска песен
	filter := models.Songs{
		Song:        r.URL.Query().Get("song"),
		Group:       r.URL.Query().Get("group"),
		Text:        r.URL.Query().Get("text"),
		ReleaseDate: r.URL.Query().Get("releaseDate"),
		Link:        r.URL.Query().Get("link"),
	}
	logrus.WithField("filter", filter).Info("Songs: filter parameters")

	// Получение списка песен с использованием сервиса
	songs, totalPages, err := h.services.GetAll(filter, page, pageSize)
	if err != nil {
		logrus.WithError(err).Error("Ошибка при получении песен")
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
	logrus.WithField("response", response).Info("Songs: response")

	// Кодирование ответа в JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.WithError(err).Error("Ошибка при кодировании ответа")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//	@Summary		Create a new song
//	@Description	Create a new song by sending a request to an external API and then saving the song details.
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			song	body	models.Songs	true	"Song details"
//	@Success		201
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/songs [post]
func (h *Handler) NewSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var song models.Songs
	defer r.Body.Close()

	// Декодирование JSON из тела запроса
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		logrus.WithError(err).Error("Ошибка при декодировании JSON")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	songStr := url.QueryEscape(song.Song)
	groupStr := url.QueryEscape(song.Group)

	apiURL := fmt.Sprintf("%s?group=%s&song=%s", h.url, groupStr, songStr)
	logrus.WithField("apiURL", apiURL).Info("Requesting data from external API")

	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	logrus.WithField("status", resp.StatusCode).Info("Received response from external API")

	if resp.StatusCode != http.StatusOK {
		logrus.WithField("status", resp.StatusCode).Error("External API returned non-OK status")
		http.Error(w, "External API returned non-OK status", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error("Failed to read response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.WithField("responseBody", string(body)).Info("Response body from external API")

	if err := json.Unmarshal(body, &song); err != nil {
		logrus.WithError(err).Error("Failed to unmarshal song data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.WithFields(logrus.Fields{
		"song":        song,
		"releaseDate": song.ReleaseDate,
	}).Info("Новая песня")

	// Создание новой песни с использованием сервиса
	if err := h.services.Create(song); err != nil {
		logrus.WithError(err).Error("Ошибка при создании песни")
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
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/songs/{id} [get]
func (h *Handler) SongByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получение номера куплета из параметров запроса
	vers := getQueryParamAsInt(r, "vers", 0)
	logrus.WithField("vers", vers).Info("SongByID: verse number")

	// Получение ID песни из переменных маршрута
	vars := mux.Vars(r)
	songID, err := strconv.Atoi(vars["id"])
	if err != nil {
		logrus.WithError(err).Error("Ошибка при преобразовании songID в int")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.WithField("songID", songID).Info("SongByID: songID")

	// Получение песни по ID с использованием сервиса
	song, err := h.services.GetByID(songID)
	if err != nil {
		logrus.WithError(err).Error("Ошибка при получении песни")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Разделение текста песни на куплеты, если указан номер куплета
	if vers != 0 {
		verse, err := splitTextDoubleMargins(song.Text, vers-1)
		if err != nil {
			logrus.WithError(err).Error("Ошибка при разделении текста на куплеты")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		song.Text = verse
	}

	// Кодирование ответа в JSON
	if err := json.NewEncoder(w).Encode(song); err != nil {
		logrus.WithError(err).Error("Ошибка при кодировании ответа")
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
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/songs/{id} [patch]
func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получение ID песни из переменных маршрута
	vars := mux.Vars(r)
	songID, err := strconv.Atoi(vars["id"])
	if err != nil {
		logrus.WithError(err).Error("Ошибка при преобразовании songID в int")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.WithField("songID", songID).Info("UpdateSong: songID")

	var song models.Songs
	defer r.Body.Close()

	// Декодирование JSON из тела запроса
	err = json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		logrus.WithError(err).Error("Ошибка при декодировании JSON")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновление песни с использованием сервиса
	if err := h.services.Update(songID, song); err != nil {
		logrus.WithError(err).Error("Ошибка при обновлении песни")
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
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/songs/{id} [delete]
func (h *Handler) DeleteSongs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получение ID песни из переменных маршрута
	vars := mux.Vars(r)
	songID, err := strconv.Atoi(vars["id"])
	if err != nil {
		logrus.WithError(err).Error("Ошибка при преобразовании songID в int")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.WithField("songID", songID).Info("DeleteSongs: songID")

	// Удаление песни с использованием сервиса
	if err := h.services.Delete(songID); err != nil {
		logrus.WithError(err).Error("Ошибка при удалении песни")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Функция для получения параметра из запроса в виде целого числа
func getQueryParamAsInt(r *http.Request, param string, defaultValue int) int {
	valueStr := r.URL.Query().Get(param)
	logrus.WithFields(logrus.Fields{
		"param":    param,
		"valueStr": valueStr,
	}).Debug("getQueryParamAsInt: param and valueStr")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"param": param,
			"error": err,
			"value": valueStr,
		}).Debug("getQueryParamAsInt: error converting valueStr to int")
		return defaultValue
	}
	return value
}

// Функция для разделения текста на куплеты
func splitTextDoubleMargins(text string, vers int) (string, error) {
	logrus.WithFields(logrus.Fields{
		"text": text,
		"vers": vers,
	}).Debug("splitTextDoubleMargins: text and vers")
	parts := strings.Split(text, "\n\n")

	// Проверка границ массива
	if vers < 0 || vers >= len(parts) {
		return "", fmt.Errorf("invalid verse index: %d", vers+1)
	}

	return parts[vers], nil
}
