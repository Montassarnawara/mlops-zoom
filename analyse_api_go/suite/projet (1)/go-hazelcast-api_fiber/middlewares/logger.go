package middlewares

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RequestLogger() fiber.Handler {
	// Création du dossier /logs/ si absent
	os.MkdirAll("logs", os.ModePerm)

	// Ouverture/append du fichier de log
	logFile, err := os.OpenFile("logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Erreur d'ouverture du fichier de log : %v", err)
	}

	logger := log.New(logFile, "", log.LstdFlags)

	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Traiter la requête
		err := c.Next()

		duration := time.Since(start)

		logger.Printf("%s %s - %d - %s",
			c.Method(),
			c.OriginalURL(),
			c.Response().StatusCode(),
			duration)

		return err
	}
}
