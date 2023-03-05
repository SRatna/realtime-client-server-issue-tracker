import { useEffect } from "react"

function App() {
  const payload = {
    id: 1,
    noOfTasks: 2,
    completed: false,
    tasks: [
      {
        id: 1,
        estimate: 10
      },
      {
        id: 2,
        estimate: 10
      }
    ]
  };

  const getRandomInt = (min, max) => {
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min) + min);
  }

  const generateStories = async () => {
    let i = 0
    while (true) {
      i++
      const noOfTasks = getRandomInt(50, 100);
      const tasks = [];
      for (let j = 0; j < noOfTasks; j++) {
        const estimate = getRandomInt(1, 50);
        tasks.push({ id: j, estimate });
      }
      const story = { id: i, noOfTasks, tasks, completed: false };
      fetch('/api/jobs', {
        method: 'post',
        body: JSON.stringify(story),
        headers: {
          'Content-Type': 'application/json'
        },
      })
      await new Promise(r => setTimeout(r, getRandomInt(1, 10)));
    }
  }

  useEffect(() => {
    // generateStories();
  }, [])
  return (
    <div>
      <button>Start</button>
    </div>
  )
}

export default App
