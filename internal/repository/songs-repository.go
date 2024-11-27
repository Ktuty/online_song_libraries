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

// Структура SongsRepository, которая инкапсулирует подключение к базе данных
type SongsRepository struct {
	db *pgxpool.Pool
}

// Функция для создания нового экземпляра SongsRepository с подключением к базе данных
func NewSongsRepository(db *pgxpool.Pool) *SongsRepository {
	return &SongsRepository{db: db}
}

// Формат даты для парсинга и форматирования
var FORMAT_DATE = "02.01.2006"

// Метод для получения всех песен с фильтрацией, пагинацией и возвратом общего количества страниц
func (r *SongsRepository) GetAllSongs(filter models.Songs, page, pageSize int) ([]models.Songs, int, error) {
	// Получение ID групп по имени группы
	groupIDs, err := r.getGroupIDsByName(filter.Group)
	if err != nil {
		return nil, 0, err
	}

	// Построение SQL-запроса и аргументов
	query, args := buildQuery(filter, page, pageSize, groupIDs)

	// Выполнение запроса к базе данных
	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Сканирование строк и преобразование их в модели песен
	songs, err := scanRows(rows, r)
	if err != nil {
		return nil, 0, err
	}

	// Подсчет общего количества строк
	totalRows, err := countRows(r.db, filter)
	if err != nil {
		return nil, 0, err
	}

	// Вычисление общего количества страниц
	totalPages := (totalRows + pageSize - 1) / pageSize
	return songs, totalPages, nil
}

// Метод для получения песни по ID
func (r *SongsRepository) GetSongByID(id int) (models.Songs, error) {
	// Построение SQL-запроса и аргументов
	var query string
	var args []interface{}
	query = `SELECT id, group_id, song, text, release_date, link FROM songs WHERE id = $1`
	args = append(args, id)

	// Выполнение запроса к базе данных
	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return models.Songs{}, err
	}
	defer rows.Close()

	// Сканирование строки и преобразование её в модель песни
	var song models.Songs
	var releaseDate time.Time
	var groupID int

	if rows.Next() {
		err := rows.Scan(&song.ID, &groupID, &song.Song, &song.Text, &releaseDate, &song.Link)
		if err != nil {
			return models.Songs{}, err
		}
		song.Group = r.groupNameByID(groupID)
		song.ReleaseDate = releaseDate.Format(FORMAT_DATE)
	} else {
		return models.Songs{}, fmt.Errorf("song with id %d not found", id)
	}

	if err := rows.Err(); err != nil {
		return models.Songs{}, err
	}

	return song, nil
}

// Метод для создания новой песни
func (r *SongsRepository) PostSong(song models.Songs) error {
	// Убедиться, что группа существует или создать её
	groupID, err := r.ensureGroupExists(song.Group)
	if err != nil {
		return err
	}

	// Парсинг даты релиза
	releaseDate, err := parseStrToTime(song.ReleaseDate)
	if err != nil {
		return fmt.Errorf("invalid release date format: %v", err)
	}

	// Построение SQL-запроса для вставки новой песни
	query := `INSERT INTO songs (group_id, song, text, release_date, link) VALUES ($1, $2, $3, $4, $5)`
	_, err = r.db.Exec(context.Background(), query, groupID, song.Song, song.Text, releaseDate, song.Link)
	return err
}

// Метод для обновления песни по ID
func (r *SongsRepository) UpdateSong(songID int, song models.Songs) error {
	// Убедиться, что группа существует или создать её
	groupID, err := r.ensureGroupExists(song.Group)
	if err != nil {
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

	_, err = r.db.Exec(context.Background(), query, args...)
	if err != nil {
		return fmt.Errorf("invalid update id: %v data: %v", songID, err)
	}

	return nil
}

// Метод для удаления песни по ID
func (r *SongsRepository) DeleteSong(songID int) error {
	// Получение group_id для песни
	var groupID int
	query := `SELECT group_id FROM songs WHERE id = $1`
	err := r.db.QueryRow(context.Background(), query, songID).Scan(&groupID)
	if err != nil {
		return fmt.Errorf("error getting group_id for song id %d: %v", songID, err)
	}

	// Подсчет количества песен в группе
	var count int
	query = `SELECT COUNT(*) FROM songs WHERE group_id = $1 AND id <> $2`
	err = r.db.QueryRow(context.Background(), query, groupID, songID).Scan(&count)
	if err != nil {
		return fmt.Errorf("error counting songs for group_id %d: %v", groupID, err)
	}

	// Удаление группы, если она больше не используется
	if count == 0 {
		query = `DELETE FROM groups WHERE id = $1`
		_, err = r.db.Exec(context.Background(), query, groupID)
		if err != nil {
			return fmt.Errorf("error deleting group with id %d: %v", groupID, err)
		}
	}

	// Удаление песни
	query = `DELETE FROM songs WHERE id = $1`
	_, err = r.db.Exec(context.Background(), query, songID)
	if err != nil {
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
	var groupID int
	err := r.db.QueryRow(context.Background(), query, groupName).Scan(&groupID)
	if err == pgx.ErrNoRows {
		// Группа не существует, создать её
		insertQuery := `INSERT INTO groups ("group") VALUES ($1) RETURNING id`
		err := r.db.QueryRow(context.Background(), insertQuery, groupName).Scan(&groupID)
		if err != nil {
			return 0, fmt.Errorf("error inserting group: %v", err)
		}
		return groupID, nil
	} else if err != nil {
		return 0, fmt.Errorf("error querying group: %v", err)
	}
	return groupID, nil
}

// Метод для получения ID групп по имени группы
func (r *SongsRepository) getGroupIDsByName(groupName string) ([]int, error) {
	if groupName == "" {
		return nil, nil
	}

	// Построение SQL-запроса для поиска групп
	query := `SELECT id FROM groups WHERE "group" ILIKE '%' || $1 || '%'`
	rows, err := r.db.Query(context.Background(), query, groupName)
	if err != nil {
		return nil, fmt.Errorf("error querying groups: %v", err)
	}
	defer rows.Close()

	var groupIDs []int
	for rows.Next() {
		var groupID int
		err := rows.Scan(&groupID)
		if err != nil {
			return nil, fmt.Errorf("error scanning group ID: %v", err)
		}
		groupIDs = append(groupIDs, groupID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error scanning group IDs: %v", err)
	}

	return groupIDs, nil
}

// Метод для получения имени группы по ID
func (r *SongsRepository) groupNameByID(ID int) string {
	var groupName string
	query := `SELECT "group" FROM groups WHERE id = $1`
	err := r.db.QueryRow(context.Background(), query, ID).Scan(&groupName)
	if err != nil {
		return ""
	}

	return groupName
}

// Функция для парсинга строки в время
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

// Функция для построения SQL-запроса с фильтрацией, пагинацией и аргументами
func buildQuery(filter models.Songs, page, pageSize int, groupIDs []int) (string, []interface{}) {
	query := `SELECT id, group_id, song, text, release_date, link FROM songs WHERE 1=1`
	args := []interface{}{}

	query, args = appendFilter(query, args, filter, groupIDs)

	query += ` LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	offset := (page - 1) * pageSize

	args = append(args, pageSize, offset)
	return query, args
}

// Функция для добавления фильтров к SQL-запросу
func appendFilter(query string, args []interface{}, filter models.Songs, groupIDs []int) (string, []interface{}) {
	argIndex := 1
	if filter.Song != "" {
		query += ` AND song ILIKE $` + strconv.Itoa(argIndex)
		args = append(args, "%"+filter.Song+"%")
		argIndex++
	}
	if len(groupIDs) > 0 {
		query += ` AND group_id = ANY($` + strconv.Itoa(argIndex) + `)`
		args = append(args, groupIDs)
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

// Функция для сканирования строк и преобразования их в модели песен
func scanRows(rows pgx.Rows, r *SongsRepository) ([]models.Songs, error) {
	var songs []models.Songs

	for rows.Next() {
		var song models.Songs
		var groupID int
		var releaseDate time.Time
		err := rows.Scan(&song.ID, &groupID, &song.Song, &song.Text, &releaseDate, &song.Link)
		if err != nil {
			return nil, err
		}
		song.Group = r.groupNameByID(groupID)
		song.ReleaseDate = releaseDate.Format(FORMAT_DATE)
		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}

// Функция для подсчета общего количества строк
func countRows(db *pgxpool.Pool, filter models.Songs) (int, error) {
	query, args := buildCountQuery(filter)
	var totalRows int
	err := db.QueryRow(context.Background(), query, args...).Scan(&totalRows)
	if err != nil {
		return 0, err
	}
	return totalRows, nil
}

// Функция для построения SQL-запроса для подсчета строк
func buildCountQuery(filter models.Songs) (string, []interface{}) {
	query := `SELECT COUNT(*) FROM songs WHERE 1=1`
	args := []interface{}{}

	query, args = appendFilter(query, args, filter, nil)

	return query, args
}
