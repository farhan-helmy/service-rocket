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
	app.Post("/api/v1/image/multiple-upload", controller.UploadMultipleImage)
}
func main() {
	app := fiber.New()
	//width, height := utils.GetImageDimension("jpegtest.jpeg")
	//utils.DecodeImage("jpegtest.jpeg")
	//fmt.Println("Width:", width, "Height:", height)
	setupRoutes(app)

	log.Fatal(app.Listen(":3001"))
}
