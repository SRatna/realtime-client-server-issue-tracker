package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type TaskPayload struct {
	Id       int `json:"id"`
	Estimate int `json:"estimate"`
}

type StoryPayload struct {
	Id        int           `json:"id"`
	NoOfTasks int           `json:"noOfTasks"`
	Completed bool          `json:"completed"`
	Tasks     []TaskPayload `json:"tasks"`
}

func main() {
	app := fiber.New()

	app.Post("api/jobs", func(c *fiber.Ctx) error {
		fmt.Println("cool")
		story := new(StoryPayload)
		if err := c.BodyParser(&story); err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("%v", story)
		return nil
	})

	log.Fatal(app.Listen(":3000"))
}
