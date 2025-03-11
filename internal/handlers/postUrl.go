package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sol1corejz/ozon-test/internal/models"
	"github.com/sol1corejz/ozon-test/internal/storage"
)

// PostUrl принимает JSON {"url": "https://example.com"} и возвращает короткий URL
func PostUrl(c *fiber.Ctx) error {
	var request struct {
		URL      string `json:"url"`
		ShortUrl string `json:"short_url"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	url := models.URL{
		URL:      request.URL,
		ShortUrl: request.ShortUrl,
	}

	err := storage.ActiveStorage.PostURL(url)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Short URL already exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save URL"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"short_url": request.ShortUrl,
	})
}
