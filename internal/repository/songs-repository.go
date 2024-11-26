package repository

import (
	"context"
	"fmt"
	"github.com/Ktuty/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"strconv"
	"time"
)

type SongsRepository struct {
	db *pgxpool.Pool
}

func NewSongsRepository(db *pgxpool.Pool) *SongsRepository {
	return &SongsRepository{db: db}
}

var FORMAT_DATE = "02.01.2006"

func (r *SongsRepository) GetAllSongs(filter models.Songs, page, pageSize int) ([]models.Songs, int, error) {
	query, args := buildQuery(filter, page, pageSize)

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	songs, err := scanRows(rows)
	if err != nil {
		return nil, 0, err
	}

	totalRows, err := countRows(r.db, filter)
	if err != nil {
		return nil, 0, err
	}

	totalPages := (totalRows + pageSize - 1) / pageSize
	return songs, totalPages, nil
}

func (r *SongsRepository) GetSongByID(id int) (models.Songs, error) {
	var query string
	var args []interface{}
	query = `SELECT id, song, "group", text, release_date, link FROM songs WHERE id = $1`
	args = append(args, id)

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return models.Songs{}, err
	}
	defer rows.Close()

	var song models.Songs
	var releaseDate time.Time

	if rows.Next() {
		err := rows.Scan(&song.ID, &song.Song, &song.Group, &song.Text, &releaseDate, &song.Link)
		if err != nil {
			return models.Songs{}, err
		}
		song.ReleaseDate = releaseDate.Format(FORMAT_DATE)
	} else {
		return models.Songs{}, fmt.Errorf("song with id %d not found", id)
	}

	if err := rows.Err(); err != nil {
		return models.Songs{}, err
	}

	return song, nil
}

func (r *SongsRepository) PostSong(song models.Songs) error {

	releaseDate, err := parseStrToTime(song.ReleaseDate)
	if err != nil {
		return fmt.Errorf("invalid release date format: %v", err)
	}

	query := `INSERT INTO songs (song, "group", text, release_date, link) VALUES ($1, $2, $3, $4, $5)`
	_, err = r.db.Exec(context.Background(), query, song.Song, song.Group, song.Text, releaseDate, song.Link)
	return err
}

func (r *SongsRepository) UpdateSong(songID int, song models.Songs) error {
	query := `UPDATE songs SET `
	var args []interface{}
	var argIndex = 2

	if song.Song != "" {
		query += `song = $` + strconv.Itoa(argIndex)
		args = append(args, song.Song)
		argIndex++
	}
	if song.Group != "" {
		if len(args) > 0 {
			query += `, `
		}
		query += `"group" = $` + strconv.Itoa(argIndex)
		args = append(args, song.Group)
		argIndex++
	}
	if song.Text != "" {
		if len(args) > 0 {
			query += `, `
		}
		query += `text = $` + strconv.Itoa(argIndex)
		args = append(args, song.Text)
		argIndex++
	}
	if song.ReleaseDate != "" {
		releaseDate, err := parseStrToTime(song.ReleaseDate)
		if err != nil {
			return fmt.Errorf("invalid release date format: %v", err)
		}
		if len(args) > 0 {
			query += `, `
		}
		query += `release_date = $` + strconv.Itoa(argIndex)
		args = append(args, releaseDate)
		argIndex++
	}
	if song.Link != "" {
		if len(args) > 0 {
			query += `, `
		}
		query += `link = $` + strconv.Itoa(argIndex)
		args = append(args, song.Link)
		argIndex++
	}

	query += ` WHERE id = $1`
	args = append([]interface{}{songID}, args...)

	_, err := r.db.Exec(context.Background(), query, args...)
	if err != nil {
		return fmt.Errorf("invalid update id: %v data: %v", songID, err)
	}

	return nil
}

func (r *SongsRepository) DeleteSong(songID int) error {
	query := `DELETE FROM songs WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, songID)
	if err != nil {
		return fmt.Errorf("invalid delete id: %v data: %v", songID, err)
	}
	return nil
}

func parseStrToTime(strTime string) (releaseDate time.Time, err error) {
	if strTime == "" {
		return time.Now(), nil
	}

	releaseDate, err = time.Parse(FORMAT_DATE, strTime)
	if err != nil {
		return releaseDate, err
	}

	return releaseDate, nil

}

func buildQuery(filter models.Songs, page, pageSize int) (string, []interface{}) {
	query := `SELECT id, song, "group", text, release_date, link FROM songs WHERE 1=1`
	args := []interface{}{}

	query, args = appendFilter(query, args, filter)

	query += ` LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	offset := (page - 1) * pageSize

	args = append(args, pageSize, offset)
	return query, args
}

func appendFilter(query string, args []interface{}, filter models.Songs) (string, []interface{}) {
	argIndex := 1
	if filter.Song != "" {
		query += ` AND song ILIKE $` + strconv.Itoa(argIndex)
		args = append(args, "%"+filter.Song+"%")
		argIndex++
	}
	if filter.Group != "" {
		query += ` AND "group" ILIKE $` + strconv.Itoa(argIndex)
		args = append(args, "%"+filter.Group+"%")
		argIndex++
	}
	if filter.Text != "" {
		query += ` AND text ILIKE $` + strconv.Itoa(argIndex)
		args = append(args, "%"+filter.Text+"%")
		argIndex++
	}
	if filter.ReleaseDate != "" {
		releaseDate, err := parseStrToTime(filter.ReleaseDate)
		if err != nil {
			return "", nil
		}
		query += ` AND release_date = $` + strconv.Itoa(argIndex)
		args = append(args, releaseDate)
		argIndex++
	}
	if filter.Link != "" {
		query += ` AND link ILIKE $` + strconv.Itoa(argIndex)
		args = append(args, "%"+filter.Link+"%")
		argIndex++
	}
	return query, args
}

func scanRows(rows pgx.Rows) ([]models.Songs, error) {
	var songs []models.Songs
	for rows.Next() {
		var song models.Songs
		var releaseDate time.Time
		err := rows.Scan(&song.ID, &song.Song, &song.Group, &song.Text, &releaseDate, &song.Link)
		if err != nil {
			return nil, err
		}
		song.ReleaseDate = releaseDate.Format(FORMAT_DATE)
		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}

func countRows(db *pgxpool.Pool, filter models.Songs) (int, error) {
	query, args := buildCountQuery(filter)
	var totalRows int
	err := db.QueryRow(context.Background(), query, args...).Scan(&totalRows)
	if err != nil {
		return 0, err
	}
	return totalRows, nil
}

func buildCountQuery(filter models.Songs) (string, []interface{}) {
	query := `SELECT COUNT(*) FROM songs WHERE 1=1`
	args := []interface{}{}

	query, args = appendFilter(query, args, filter)

	return query, args
}
