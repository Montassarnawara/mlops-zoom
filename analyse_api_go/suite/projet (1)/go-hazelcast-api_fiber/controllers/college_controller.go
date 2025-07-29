package controllers

import (
	"go-hazelcast-api/cache"

	"github.com/gofiber/fiber/v2"
)

func GetCollegeByID(c *fiber.Ctx) error {
	id := c.Params("id")

	college, err := cache.GetFromCache(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erreur lors de l'accès au cache",
		})
	}

	if college == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Collège non trouvé dans le cache",
		})
	}

	return c.Status(fiber.StatusOK).JSON(college)
}
