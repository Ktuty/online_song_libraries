package services

import (
	"github.com/Ktuty/internal/models"
	"github.com/Ktuty/internal/repository"
)

// Структура SongsService, которая инкапсулирует репозиторий для работы с песнями
type SongsService struct {
	rep *repository.Repository
}

// Функция для создания нового экземпляра SongsService с заданным репозиторием
func NewSongsService(rep *repository.Repository) *SongsService {
	return &SongsService{rep}
}

// Метод для получения всех песен с фильтрацией, пагинацией и возвратом общего количества страниц
func (s *SongsService) GetAll(filter models.Songs, page, pageSize int) ([]models.Songs, int, error) {
	return s.rep.GetAllSongs(filter, page, pageSize)
}

// Метод для получения песни по ID
func (s *SongsService) GetByID(id int) (models.Songs, error) {
	return s.rep.GetSongByID(id)
}

// Метод для создания новой песни
func (s *SongsService) Create(song models.Songs) error {
	return s.rep.PostSong(song)
}

// Метод для обновления песни по ID
func (s *SongsService) Update(songID int, song models.Songs) error {
	return s.rep.UpdateSong(songID, song)
}

// Метод для удаления песни по ID
func (s *SongsService) Delete(songID int) error {
	return s.rep.DeleteSong(songID)
}
