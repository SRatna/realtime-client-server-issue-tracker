package main

import (
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

type Task struct {
	id       int
	estimate int // in ms
}

type Story struct {
	id        int
	noOfTasks int
	tasks     chan *Task
}

func pushStory(stories chan<- *Story, storyPayload *StoryPayload) {
	tasks := make(chan *Task, storyPayload.NoOfTasks)
	for j := 0; j < storyPayload.NoOfTasks; j++ {
		taskPayload := storyPayload.Tasks[j]
		task := &Task{id: taskPayload.Id, estimate: taskPayload.Estimate}
		tasks <- task
	}
	story := &Story{id: storyPayload.Id, noOfTasks: storyPayload.NoOfTasks, tasks: tasks}
	stories <- story
}

func main() {
	app := fiber.New()

	stories := make(chan *Story)

	app.Post("api/jobs", func(c *fiber.Ctx) error {
		story := new(StoryPayload)
		if err := c.BodyParser(&story); err != nil {
			return err
		}
		pushStory(stories, story)
		return nil
	})

	log.Fatal(app.Listen(":3000"))
}
