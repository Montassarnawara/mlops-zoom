package routes

import (
	"go-hazelcast-api/controllers"
	"go-hazelcast-api/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Use(middlewares.RequestLogger()) // middleware logging

	// Route principale
	app.Get("/college/:id", controllers.GetCollegeByID)
}
