// Package config отвечает за чтение конфигурации приложения из переменных окружения и флагов.
package config

import (
	"flag"
	"os"
)

// Структура для хранения конфигурации из JSON-файла.
type Config struct {
	ServerAddress   string `json:"server_address"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
}

// Переменные для хранения значений env и флагов.
var (
	// FlagRunAddr содержит адрес и порт для запуска сервера.
	FlagRunAddr string
	// FlagLogLevel задает уровень логирования приложения.
	FlagLogLevel string
	// FileStoragePath определяет путь к файлу для хранения данных.
	FileStoragePath string
	// DatabaseDSN содержит строку подключения к базе данных.
	DatabaseDSN string
)

// ParseFlags читает флаги командной строки и переменные окружения.
// Если указаны как флаги, так и переменные окружения, приоритет имеют значения из переменных окружения.
func ParseFlags() {
	// Инициализация флагов командной строки.
	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&FlagLogLevel, "l", "info", "log level")
	flag.StringVar(&FileStoragePath, "f", "", "file storage path")
	flag.StringVar(&DatabaseDSN, "d", "", "database dsn")
	flag.Parse()

	// Переопределение значений флагов переменными окружения (если они заданы).
	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}

	if envStoragePath := os.Getenv("FILE_STORAGE_PATH"); envStoragePath != "" {
		FileStoragePath = envStoragePath
	}

	if databaseDsn := os.Getenv("DATABASE_DSN"); databaseDsn != "" {
		DatabaseDSN = databaseDsn
	}
}
