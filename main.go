package main

import (
	"log"

	"github.com/farhan-helmy/service-rocket/controller"
	"github.com/gofiber/fiber/v2"
)

func welcome(c *fiber.Ctx) error {
	return c.SendString("welcome to my AWSOME API")
}

func setupRoutes(app *fiber.App) {
	app.Get("/api/v1", welcome)
	app.Post("/api/v1/image/upload", controller.UploadImage)
}
func main() {
	app := fiber.New()

	setupRoutes(app)

	log.Fatal(app.Listen(":3001"))
}
