import { useState } from "react"

function App() {
  const [noOfStoriesPerSecond, setNoOfStoriesPerSecond] = useState(0);
  const [noOfStoriesDonePerSecond, setNoOfStoriesDonePerSecond] = useState(0);
  const [started, setStarted] = useState(false);

  const getRandomInt = (min, max) => {
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min) + min);
  }

  const generateStories = async () => {
    console.log('generateStories');
    let i = 0;
    let noOfStories = 0;

    setInterval(() => {
      setNoOfStoriesPerSecond(noOfStories);
      noOfStories = 0;
    }, 1000);

    while (true) {
      i++;
      noOfStories++;
      const noOfTasks = getRandomInt(20, 80);
      const tasks = [];
      for (let j = 0; j < noOfTasks; j++) {
        const estimate = getRandomInt(10, 20);
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
      await new Promise(r => setTimeout(r, getRandomInt(50, 100)));
    }
  }

  const startJob = () => {
    setStarted(true);
    generateStories();
    const ws = new WebSocket('ws://localhost:3000/ws/status');
    ws.onopen = () => { console.log('ws opened'); }
    ws.onmessage = async (evt) => { 
      const msg = await evt.data.text();
      setNoOfStoriesDonePerSecond(msg);
    }
  }
  
  return (
    <div>
      {!started && (<button onClick={startJob} style={{ margin: 16 }}>Start</button>)}
      {started && (<>
        <div style={{ padding: 16 }}>
          No of stories produced per second: {noOfStoriesPerSecond}
        </div>
        <div style={{ padding: 16 }}>
          No of stories done per second: {noOfStoriesDonePerSecond}
        </div>
      </>)}
    </div>
  )
}

export default App
