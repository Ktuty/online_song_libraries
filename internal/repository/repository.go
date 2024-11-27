package repository

import (
	"github.com/Ktuty/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Интерфейс Songs, определяющий методы для работы с песнями
type Songs interface {
	// Метод для получения всех песен с фильтрацией, пагинацией и возвратом общего количества страниц
	GetAllSongs(filter models.Songs, page, pageSize int) ([]models.Songs, int, error)
	// Метод для получения песни по ID
	GetSongByID(id int) (models.Songs, error)
	// Метод для создания новой песни
	PostSong(song models.Songs) error
	// Метод для обновления песни по ID
	UpdateSong(songID int, song models.Songs) error
	// Метод для удаления песни по ID
	DeleteSong(songID int) error
}

// Структура Repository, реализующая интерфейс Songs
type Repository struct {
	Songs
}

// Функция для создания нового экземпляра Repository с подключением к базе данных
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Songs: NewSongsRepository(db), // Инициализация репозитория песен с подключением к базе данных
	}
}
