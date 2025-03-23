package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"user-actions-api/handlers"
	"user-actions-api/storage"
)

func main() {
	if err := storage.LoadData(); err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName:      "User Actions API",
		ErrorHandler: customErrorHandler,
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	app.Use(cache.New(cache.Config{
		Expiration:   5 * time.Minute,
		CacheControl: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API is running")
	})

	app.Get("/users/:id", handlers.GetUserByID)

	app.Get("/users/:id/actions/count", handlers.GetUserActionCount)

	app.Get("/actions/:type/next", handlers.GetNextActionBreakdown)

	app.Get("/referral-indices", handlers.GetReferralIndices)


	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}