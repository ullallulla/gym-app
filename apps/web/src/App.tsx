import { useState } from "react";

type Workout = {
  id: string;
  user_id: string;
  name: string;
  started_at: string;
  created_at: string;
  revision: number;
};

function App() {
  const [apiStatus, setApiStatus] = useState("Not checked");
  const [workouts, setWorkouts] = useState<Workout[]>([]);
  const [workoutStatus, setWorkoutStatus] = useState("Not loaded");

  async function checkApi() {
    const url = "/api/health";
    setApiStatus("Checking...");
    try {
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
      }
      const result = await response.json();
      setApiStatus(result.status);
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      setApiStatus(`Error: ${message}`);
    }
  }

  async function loadWorkouts() {
    const url = "/api/workouts";
    setWorkoutStatus("Loading...");
    try {
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error(`Response.status: ${response.status}`);
      }
      const result: Workout[] = await response.json();
      setWorkouts(result);
      setWorkoutStatus("Loaded");
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      setWorkoutStatus(`Error: ${message}`);
    }
  }

  return (
    <>
      <h1>Gym app</h1>
      <button onClick={checkApi}>Check API</button>
      <p>{apiStatus}</p>

      <button onClick={loadWorkouts}>Load workouts</button>
      <p>{workoutStatus}</p>

      <ul>
        {workouts.map((workout) => (
          <li key={workout.id}>{workout.name}</li>
        ))}
      </ul>
    </>
  );
}

export default App;
