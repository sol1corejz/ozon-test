package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sol1corejz/ozon-test/internal/storage"
)

// GetUrl принимает JSON {"short_url": "exmpl"} и возвращает {"original_url": "https://example.com"}
func GetUrl(c *fiber.Ctx) error {
	var request struct {
		ShortUrl string `json:"short_url"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	originalUrl, err := storage.ActiveStorage.GetURL(request.ShortUrl)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "URL not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.JSON(fiber.Map{
		"original_url": originalUrl,
	})
}
