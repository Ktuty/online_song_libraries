package services

import (
	"github.com/Ktuty/internal/models"
	"github.com/Ktuty/internal/repository"
)

type SongsService struct {
	rep *repository.Repository
}

func NewSongsService(rep *repository.Repository) *SongsService {
	return &SongsService{rep}
}

func (s *SongsService) GetAll(filter models.Songs, page, pageSize int) ([]models.Songs, int, error) {
	return s.rep.GetAllSongs(filter, page, pageSize)
}
func (s *SongsService) GetByID(id int) (models.Songs, error) {
	return s.rep.GetSongByID(id)
}
func (s *SongsService) Create(song models.Songs) error {
	return s.rep.PostSong(song)
}
func (s *SongsService) Update(songID int, song models.Songs) error {
	return s.rep.UpdateSong(songID, song)
}
func (s *SongsService) Delete(songID int) error {
	return s.rep.DeleteSong(songID)
}
