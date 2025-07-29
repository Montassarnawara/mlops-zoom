package main

import (
	"fmt"
	"go-hazelcast-api/cache"
	"go-hazelcast-api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cache.InitHazelcast()
	defer cache.Client.Shutdown(cache.Ctx)

	fmt.Println("  API d  marr  e avec Hazelcast + Fiber")

	app := fiber.New()
	routes.SetupRoutes(app)

	app.Listen(":8080")
}
