package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Ktuty/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"strconv"
)

// Структура SongsRepository, которая инкапсулирует подключение к базе данных
type SongsRepository struct {
	db *pgxpool.Pool
}

// Функция для создания нового экземпляра SongsRepository с подключением к базе данных
func NewSongsRepository(db *pgxpool.Pool) *SongsRepository {
	return &SongsRepository{db: db}
}

// Метод для получения всех песен с фильтрацией, пагинацией и возвратом общего количества страниц
func (r *SongsRepository) GetAllSongs(filter models.Songs, page, pageSize int) ([]models.Songs, int, error) {
	offset := (page - 1) * pageSize

	query := `
	SELECT s.id, s.song, g."group", s.text, s.release_date, s.link
	FROM songs s
	INNER JOIN groups g ON s.group_id = g.id
	WHERE s.song ILIKE $1 AND g."group" ILIKE $2 AND s.text ILIKE $3 AND s.release_date ILIKE $4 AND s.link ILIKE $5
	LIMIT $6 OFFSET $7`

	logrus.WithFields(logrus.Fields{
		"query":  query,
		"params": []interface{}{"%" + filter.Song + "%", "%" + filter.Group + "%", "%" + filter.Text + "%", "%" + filter.ReleaseDate + "%", "%" + filter.Link + "%", pageSize, offset},
	}).Debug("Executing query")

	rows, err := r.db.Query(context.Background(), query, "%"+filter.Song+"%", "%"+filter.Group+"%", "%"+filter.Text+"%", "%"+filter.ReleaseDate+"%", "%"+filter.Link+"%", pageSize, offset)
	if err != nil {
		logrus.WithError(err).Error("SongsRepository.GetAllSongs query error")
		return nil, 0, fmt.Errorf("SongsRepository.GetAllSongs query error: %w", err)
	}
	defer rows.Close()

	var songs []models.Songs

	for rows.Next() {
		var song models.Songs
		if err := rows.Scan(&song.ID, &song.Song, &song.Group, &song.Text, &song.ReleaseDate, &song.Link); err != nil {
			logrus.WithError(err).Error("SongsRepository.GetAllSongs scan error")
			return nil, 0, fmt.Errorf("SongsRepository.GetAllSongs scan error: %w", err)
		}
		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		logrus.WithError(err).Error("SongsRepository.GetAllSongs rows error")
		return nil, 0, fmt.Errorf("SongsRepository.GetAllSongs rows error: %w", err)
	}

	// Запрос для получения общего количества записей
	countQuery := `
	SELECT COUNT(*)
	FROM songs s
	INNER JOIN groups g ON s.group_id = g.id
	WHERE s.song ILIKE $1 AND g."group" ILIKE $2 AND s.text ILIKE $3 AND s.release_date ILIKE $4 AND s.link ILIKE $5`

	logrus.WithFields(logrus.Fields{
		"query":  countQuery,
		"params": []interface{}{"%" + filter.Song + "%", "%" + filter.Group + "%", "%" + filter.Text + "%", "%" + filter.ReleaseDate + "%", "%" + filter.Link + "%"},
	}).Debug("Executing count query")

	var totalRecords int
	err = r.db.QueryRow(context.Background(), countQuery, "%"+filter.Song+"%", "%"+filter.Group+"%", "%"+filter.Text+"%", "%"+filter.ReleaseDate+"%", "%"+filter.Link+"%").Scan(&totalRecords)
	if err != nil {
		logrus.WithError(err).Error("SongsRepository.GetAllSongs count query error")
		return nil, 0, fmt.Errorf("SongsRepository.GetAllSongs count query error: %w", err)
	}

	// Вычисление общего количества страниц
	totalPages := (totalRecords + pageSize - 1) / pageSize
	return songs, totalPages, nil
}

// Метод для получения песни по ID
func (r *SongsRepository) GetSongByID(id int) (models.Songs, error) {
	// Построение SQL-запроса для получения песни
	query := `SELECT s.id, g."group", s.song, s.text, s.release_date, s.link
	          FROM songs s
	          INNER JOIN groups g ON s.group_id = g.id
	          WHERE s.id = $1`

	logrus.WithFields(logrus.Fields{
		"query":  query,
		"params": id,
	}).Debug("Executing query")

	// Выполнение запроса к базе данных
	var song models.Songs
	err := r.db.QueryRow(context.Background(), query, id).Scan(&song.ID, &song.Group, &song.Song, &song.Text, &song.ReleaseDate, &song.Link)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.WithField("id", id).Info("Song not found")
			return models.Songs{}, fmt.Errorf("song with id %d not found", id)
		}
		logrus.WithError(err).Error("SongsRepository.GetSongByID query error")
		return models.Songs{}, fmt.Errorf("SongsRepository.GetSongByID query error: %w", err)
	}

	return song, nil
}

// Метод для создания новой песни
func (r *SongsRepository) PostSong(song models.Songs) error {

	// Убедиться, что группа существует или создать её
	groupID, err := r.ensureGroupExists(song.Group)
	if err != nil {
		logrus.WithError(err).Error("Error ensuring group exists")
		return err
	}

	// Построение SQL-запроса для вставки новой песни
	query := `INSERT INTO songs (group_id, song, text, release_date, link) VALUES ($1, $2, $3, $4, $5)`
	logrus.WithFields(logrus.Fields{
		"query":  query,
		"params": []interface{}{groupID, song.Song, song.Text, song.ReleaseDate, song.Link},
	}).Debug("Executing query")

	_, err = r.db.Exec(context.Background(), query, groupID, song.Song, song.Text, song.ReleaseDate, song.Link)
	if err != nil {
		logrus.WithError(err).Error("Error inserting song")
	}
	return err
}

// Метод для обновления песни по ID
func (r *SongsRepository) UpdateSong(songID int, song models.Songs) error {

	currentGroupID, err := r.getGroupIDByName(songID)
	if err != nil {
		logrus.WithError(err).Error("Error getting group ID")
		return err
	}

	// Убедиться, что группа существует или создать её
	groupID, err := r.ensureGroupExists(song.Group)
	if err != nil {
		logrus.WithError(err).Error("Error ensuring group exists")
		return err
	}

	// Построение SQL-запроса для обновления песни
	query := `UPDATE songs SET `
	var args []interface{}
	var argIndex = 2

	if song.Song != "" {
		query += `song = $` + strconv.Itoa(argIndex)
		args = append(args, song.Song)
		argIndex++
	}
	if groupID != 0 {
		if len(args) > 0 {
			query += `, `
		}
		query += `group_id = $` + strconv.Itoa(argIndex)
		args = append(args, groupID)
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
		if len(args) > 0 {
			query += `, `
		}
		query += `release_date = $` + strconv.Itoa(argIndex)
		args = append(args, song.ReleaseDate)
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

	logrus.WithFields(logrus.Fields{
		"query":  query,
		"params": args,
	}).Debug("Executing query")

	_, err = r.db.Exec(context.Background(), query, args...)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"songID": songID,
			"error":  err,
		}).Error("Error updating song")
		return fmt.Errorf("invalid update id: %v data: %v", songID, err)
	}

	if currentGroupID != groupID {
		if err = r.ensureGroupUsed(currentGroupID, songID); err != nil {
			logrus.WithFields(logrus.Fields{
				"songID": songID,
				"error":  err,
			}).Error("Error checking group_id for song")
		}
	}

	return nil
}

// Метод для удаления песни по ID
func (r *SongsRepository) DeleteSong(songID int) error {
	// Получение group_id для песни
	groupID, err := r.getGroupIDByName(songID)
	if err != nil {
		logrus.WithError(err).Error("Error getting group ID")
		return err
	}

	if err = r.ensureGroupUsed(groupID, songID); err != nil {
		logrus.WithFields(logrus.Fields{
			"songID": songID,
			"error":  err,
		}).Error("Error checking group_id for song")
	}

	// Удаление песни
	query := `DELETE FROM songs WHERE id = $1`
	logrus.WithFields(logrus.Fields{
		"query":  query,
		"params": songID,
	}).Debug("Executing query")

	_, err = r.db.Exec(context.Background(), query, songID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"songID": songID,
			"error":  err,
		}).Error("Error deleting song")
		return fmt.Errorf("error deleting song with id %d: %v", songID, err)
	}

	return nil
}

// Метод для обеспечения существования группы
func (r *SongsRepository) ensureGroupExists(groupName string) (int, error) {
	if groupName == "" {
		return 0, nil
	}

	// Построение SQL-запроса для поиска группы
	query := `SELECT id FROM groups WHERE "group" ILIKE '%' || $1 || '%' LIMIT 1`
	logrus.WithFields(logrus.Fields{
		"query":  query,
		"params": groupName,
	}).Debug("Executing query")

	var groupID int
	err := r.db.QueryRow(context.Background(), query, groupName).Scan(&groupID)
	if err == pgx.ErrNoRows {
		// Группа не существует, создать её
		insertQuery := `INSERT INTO groups ("group") VALUES ($1) RETURNING id`
		logrus.WithFields(logrus.Fields{
			"query":  insertQuery,
			"params": groupName,
		}).Debug("Executing query")

		err := r.db.QueryRow(context.Background(), insertQuery, groupName).Scan(&groupID)
		if err != nil {
			logrus.WithError(err).Error("Error inserting group")
			return 0, fmt.Errorf("error inserting group: %v", err)
		}
		return groupID, nil
	} else if err != nil {
		logrus.WithError(err).Error("Error querying group")
		return 0, fmt.Errorf("error querying group: %v", err)
	}
	return groupID, nil
}

func (r *SongsRepository) ensureGroupUsed(groupID, songID int) error {
	// Подсчет количества песен в группе
	var count int
	query := `SELECT COUNT(*) FROM songs WHERE group_id = $1 AND id <> $2`
	logrus.WithFields(logrus.Fields{
		"query":  query,
		"params": []interface{}{groupID, songID},
	}).Debug("Executing query")

	err := r.db.QueryRow(context.Background(), query, groupID, songID).Scan(&count)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"groupID": groupID,
			"error":   err,
		}).Error("Error counting songs for group")
		return fmt.Errorf("error counting songs for group_id %d: %v", groupID, err)
	}

	// Удаление группы, если она больше не используется
	if count == 0 {
		query = `DELETE FROM groups WHERE id = $1`
		logrus.WithFields(logrus.Fields{
			"query":  query,
			"params": groupID,
		}).Debug("Executing query")

		_, err = r.db.Exec(context.Background(), query, groupID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"groupID": groupID,
				"error":   err,
			}).Error("Error deleting group")
			return fmt.Errorf("error deleting group with id %d: %v", groupID, err)
		}
	}

	return nil
}

func (r *SongsRepository) getGroupIDByName(songID int) (int, error) {
	var groupID int

	query := `SELECT group_id FROM songs WHERE id = $1`
	logrus.WithFields(logrus.Fields{
		"query":  query,
		"params": songID,
	}).Debug("Executing query")

	err := r.db.QueryRow(context.Background(), query, songID).Scan(&groupID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"songID": songID,
			"error":  err,
		}).Error("Error getting group_id for song")
		return 0, fmt.Errorf("error getting group_id for song id %d: %v", songID, err)
	}

	return groupID, nil
}
