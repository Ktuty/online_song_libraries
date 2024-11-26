package services

import (
	"github.com/Ktuty/internal/models"
	"github.com/Ktuty/internal/repository"
)

type Songs interface {
	GetAll(filter models.Songs, page, pageSize int) ([]models.Songs, int, error)
	GetByID(id int) (models.Songs, error)
	Create(song models.Songs) error
	Update(songID int, song models.Songs) error
	Delete(songID int) error
}

type Service struct {
	Songs
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Songs: NewSongsService(repo),
	}
}
