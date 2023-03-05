package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
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
type Story struct {
	id        int
	noOfTasks int
	completed bool
	tasks     []Task
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

func worker(id int, tasks <-chan *Task) {
	task := <-tasks
	time.Sleep(time.Duration(task.estimate) * time.Millisecond)
}

func main() {
	app := fiber.New()

	noOfStoriesCompleted := 0
	noOfStoriesCompletedPerSecond := 0
	ticker := time.NewTicker(time.Second)
	stories := make(chan *Story, 1000)
	defer close(stories)

	go func() {
		for {
			select {
			case story := <-stories:
				go func() {
					noOfWorkers := story.noOfTasks
					var wg sync.WaitGroup
					for i := 0; i < noOfWorkers; i++ {
						wg.Add(1)
						i := i
						go func() {
							defer wg.Done()
							worker(i, story.tasks)
						}()
					}
					wg.Wait()
					noOfStoriesCompleted++
					noOfStoriesCompletedPerSecond++
					defer close(story.tasks)
				}()
			case <-ticker.C:
				fmt.Printf("no of stories completed %d\n", noOfStoriesCompleted)
			}
		}
	}()

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Static("/", "./dist")

	app.Post("api/jobs", func(c *fiber.Ctx) error {
		story := new(StoryPayload)
		if err := c.BodyParser(&story); err != nil {
			return err
		}
		go pushStory(stories, story)
		return nil
	})

	app.Get("/ws/status", websocket.New(func(c *websocket.Conn) {
		for {
			status := noOfStoriesCompletedPerSecond
			noOfStoriesCompletedPerSecond = 0
			if err := c.WriteMessage(websocket.BinaryMessage, []byte(strconv.Itoa(status))); err != nil {
				log.Println("write:", err)
				break
			}
			time.Sleep(time.Second)
		}

	}))

	log.Fatal(app.Listen(":3000"))
}
