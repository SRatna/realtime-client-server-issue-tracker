# Realtime client/server issue tracking statistics

This is a simple implementation of a realtime issue tracking web application powered by React in the frontend and Go in the backend.

## Tech Stack

This app uses a number of third party open-source tools:

### Frontend
- [Vite](https://vitejs.dev/) for building the [React](https://reactjs.org/) frontend.

### Backend
- [Air](https://github.com/cosmtrek/air) for live reloading the [Go](https://go.dev/) backend.
- [Fiber](https://docs.gofiber.io/) as web framework.
- [WebSocket middleware for Fiber](https://github.com/gofiber/websocket) for websocket connection.

## Getting started

### Requirements
You must install [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) to run the application.

### Up and Running
Run following command from root directory of the project to run the overall application.
```shell
docker-compose build
docker-compose up -d
```

It will build and start the docker image written for the app (which can be found in Dockerfile). The web app can be loaded by visiting [http://localhost:3000/](http://localhost:3000/).

## Data model
We have `Task` model with id and estimate as attributes as shown below:
```go
type Task struct {
	id       int
	estimate int // in ms
}
```
Then we have `Story` model with id, noOfTasks, completed and array of tasks as attributes as shown below:
```go
type Story struct {
	id        int
	noOfTasks int
	completed bool
	tasks     []Task
}
```
Each story can have multiple tasks and we process these tasks concurrently in the server.

## API endpoints
We have two endpoints:
- `/api/jobs` is a RESTful post endpoint which expects json data in the body of the request
which consists of the `Story` data model, for e.g.:
```json
{
  "id":1,
  "noOfTasks":3,
  "completed":false,
  "tasks":[
    {"id":0,"estimate":10},
    {"id":1,"estimate":14},
    {"id":2,"estimate":14}
  ]
}
```
- `/ws/status` is a websocket endpoint which return number of stories completed by the server every second. The client opens connection and listens for the message from the server. Once it receives the message, it shows it to the user.

## Implementation details
### Client
The client create story data and sends it to the server continuously with a random delay of 50 to 100 ms. Each story contains of 20 to 80 tasks, each task with estimated completion time of 10 to 20 ms.

### Server
As soon as server receive a story, it pushes it in a queue and as soon as the queue has a story, a thread with start working on it. The thread also spawns other threads based on the number of tasks and waits for all the tasks to finish.

## Further improvements
### Data persistance
Use [MongoDB](https://www.mongodb.com/) to save incoming stories. Run a thread that removes the completed story from DB. Also auto populate the story queue with the incomplete stories if exists in the DB during application startup.

### Presentation
Use [ChartJS](https://www.chartjs.org/) to plot bar charts of generated stories vs completed stories per second. Keep latest 20 of these stats.