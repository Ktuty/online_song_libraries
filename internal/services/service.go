package services

import (
	"github.com/Ktuty/internal/models"
	"github.com/Ktuty/internal/repository"
)

// Интерфейс Songs, определяющий методы для работы с песнями
type Songs interface {
	// Метод для получения всех песен с фильтрацией, пагинацией и возвратом общего количества страниц
	GetAll(filter models.Songs, page, pageSize int) ([]models.Songs, int, error)
	// Метод для получения песни по ID
	GetByID(id int) (models.Songs, error)
	// Метод для создания новой песни
	Create(song models.Songs) error
	// Метод для обновления песни по ID
	Update(songID int, song models.Songs) error
	// Метод для удаления песни по ID
	Delete(songID int) error
}

// Структура Service, реализующая интерфейс Songs
type Service struct {
	Songs
}

// Функция для создания нового экземпляра Service с заданным репозиторием
func NewService(repo *repository.Repository) *Service {
	return &Service{
		Songs: NewSongsService(repo), // Инициализация сервиса песен с заданным репозиторием
	}
}
