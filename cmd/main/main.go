package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sol1corejz/ozon-test/cmd/config"
	handlers2 "github.com/sol1corejz/ozon-test/internal/handlers"
	"github.com/sol1corejz/ozon-test/internal/storage"
	"log"
)

func main() {
	// Считывает флаги конфигурации и обновляет параметры запуска.
	config.ParseFlags()

	// Инициализация хранилища
	if err := storage.InitializeStorage(); err != nil {
		log.Fatal("Failed to initialize storage: ", err)
	}

	// Создаем приложение
	app := fiber.New()

	// Регестрируем маршруты
	app.Get("/", handlers2.GetUrl)
	app.Post("/", handlers2.PostUrl)
}
