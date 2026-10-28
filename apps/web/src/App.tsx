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
  const [workoutName, setWorkoutName] = useState("");
  const [createStatus, setCreateStatus] = useState("");
  const [isCreating, setIsCreating] = useState(false);
  const [selectedWorkout, setSelectedWorkout] = useState<Workout | null>(null);
  const [detailStatus, setDetailStatus] = useState("");
  const [isLoadingDetail, setIsLoadingDetail] = useState(false);

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

  async function loadWorkout(id: string) {
    const url = `/api/workouts/${id}`;

    if (isLoadingDetail) {
      return;
    }

    setDetailStatus("Loading...");
    setIsLoadingDetail(true);
    setSelectedWorkout(null);

    try {
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error(`Response.status: ${response.status}`);
      }
      const result: Workout = await response.json();
      setSelectedWorkout(result);
      setDetailStatus("Loaded");
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      setDetailStatus(`Error: ${message}`);
    } finally {
      setIsLoadingDetail(false);
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

  async function createWorkout() {
    if (isCreating || workoutName.trim() === "") return;
    setIsCreating(true);
    setCreateStatus("Creating...");
    try {
      const response = await fetch("/api/workouts", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name: workoutName,
          started_at: new Date().toISOString(),
        }),
      });
      if (!response.ok) {
        throw new Error(`Response.status: ${response.status}`);
      }
      const workout: Workout = await response.json();
      setWorkoutName("");
      setCreateStatus(`Created ${workout.name}`);
      await loadWorkouts();
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      setCreateStatus(`Error: ${message}`);
    } finally {
      setIsCreating(false);
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
          <li key={workout.id}>
            <button
              type="button"
              disabled={isLoadingDetail}
              onClick={() => void loadWorkout(workout.id)}
            >
              {workout.name}
            </button>
          </li>
        ))}
      </ul>
      <p>{detailStatus}</p>
      <p>{selectedWorkout ? selectedWorkout.name : "Select a workout"}</p>
      <label htmlFor="workout-name">Workout name</label>
      <p>{createStatus}</p>
      <form
        onSubmit={(event) => {
          event.preventDefault();
          void createWorkout();
        }}
      >
        <input
          id="workout-name"
          disabled={isCreating}
          type="text"
          value={workoutName}
          onChange={(event) => setWorkoutName(event.target.value)}
        />
        <button type="submit" disabled={isCreating || workoutName.trim() === ""}>
          Create Workout
        </button>
      </form>
    </>
  );
}

export default App;
