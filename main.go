package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("api/jobs", func(c *fiber.Ctx) error {

		return nil
	})

	log.Fatal(app.Listen(":3000"))
}
