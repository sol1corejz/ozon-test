package storage

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sol1corejz/ozon-test/cmd/config"
	"github.com/sol1corejz/ozon-test/internal/models"
	"log"
)

// Storage - интерфейс хранилища данных, предоставляющий методы для работы с пользователями и учетными данными.
// Используется для расширяемости и удобства тестирования.
type Storage interface {
	GetURL(shortUrl string) (string, error)
	PostURL(url models.URL) error
}

// DBStorageImpl - реализация интерфейса Storage, использующая базу данных PostgresSQL.
type DBStorageImpl struct {
	DB *sql.DB
}

// MemoryStorageImpl - реализация интерфейса Storage, использующая память компьютера для хранения.
type MemoryStorageImpl struct {
	Memory map[string]*models.URL
}

// Глобальные переменные хранилищ
var (
	DBStorage     DBStorageImpl
	MemoryStorage MemoryStorageImpl
	ActiveStorage Storage // Выбранное хранилище (БД или память)
)

// ErrNotFound - ошибка, возвращаемая при отсутствии данных.
var ErrNotFound = errors.New("not found")

// ErrAlreadyExists — ошибка, которая возвращается, если сокращённый URL уже существует.
var ErrAlreadyExists = errors.New("link already exists")

// InitializeStorage устанавливает соединение с базой данных и создает таблицу.
func InitializeStorage() error {
	if config.DatabaseDSN == "" {
		log.Print("Initializing memory storage")
		MemoryStorage.Memory = make(map[string]*models.URL)
		ActiveStorage = &MemoryStorage
		return nil
	}

	// Подключение к базе данных.
	db, err := sql.Open("pgx", config.DatabaseDSN)
	if err != nil {
		log.Fatal("Error opening database connection: ", err)
		return err
	}

	DBStorage.DB = db
	ActiveStorage = &DBStorage

	// Создание таблицы для хранения сокращённых URL, если она не существует.
	_, err = DBStorage.DB.Exec(`
			CREATE TABLE IF NOT EXISTS urls (
				id SERIAL PRIMARY KEY,
				url TEXT NOT NULL UNIQUE,
				short_url TEXT NOT NULL UNIQUE,
			)
		`)
	if err != nil {
		log.Fatal("Error creating table: ", err)
		return err
	}

	return nil

}

// GetURL получение ссылки
func (storage DBStorageImpl) GetURL(shortUrl string) (string, error) {
	var URLData models.URL
	row := storage.DB.QueryRow(`
		SELECT * FROM urls WHERE short_url = $1;
	`, shortUrl)

	err := row.Scan(&URLData.ID, &URLData.URL, &URLData.ShortUrl)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Print(ErrNotFound)
			return "", ErrNotFound
		}
		log.Print("failed to get url: ", err)
		return "", err
	}

	return URLData.URL, nil
}

// PostURL добавляет новую ссылку
func (storage DBStorageImpl) PostURL(u models.URL) error {
	_, err := storage.DB.Exec(`
		INSERT INTO urls (id, url, short_url) VALUES ($1, $2, $3)
	`, u.ID, u.URL, u.ShortUrl)

	if err != nil {
		// Проверяем, является ли ошибка нарушением ограничения уникальности
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Print(ErrAlreadyExists, u.URL)
			return ErrAlreadyExists
		}

		log.Print("failed to add URL: ", err)
		return err
	}

	return nil
}

// GetURL получение ссылки
func (storage MemoryStorageImpl) GetURL(shortUrl string) (string, error) {
	urlData, ok := storage.Memory[shortUrl]

	if !ok {
		log.Print(ErrNotFound)
		return "", ErrNotFound
	}

	return urlData.URL, nil
}

// PostURL добавляет новую ссылку
func (storage MemoryStorageImpl) PostURL(u models.URL) error {
	if _, ok := storage.Memory[u.ShortUrl]; ok {
		return ErrAlreadyExists
	}
	storage.Memory[u.ShortUrl] = &u
	return nil
}
