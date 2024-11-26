package repository

import (
	"github.com/Ktuty/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Songs interface {
	GetAllSongs(filter models.Songs, page, pageSize int) ([]models.Songs, int, error)
	GetSongByID(id int) (models.Songs, error)
	PostSong(song models.Songs) error
	UpdateSong(songID int, song models.Songs) error
	DeleteSong(songID int) error
}

type Repository struct {
	Songs
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Songs: NewSongsRepository(db),
	}
}
