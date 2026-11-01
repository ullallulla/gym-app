import { useState } from "react";
import {
  fetchWorkoutSets,
  createWorkoutSet,
  updateWorkoutSet,
  type WorkoutSet,
  ApiError,
} from "./workoutSets";

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
  const [workoutSets, setWorkoutSets] = useState<WorkoutSet[]>([]);
  const [exerciseName, setExerciseName] = useState("");
  const [reps, setReps] = useState("");
  const [weightKg, setWeightKg] = useState("");
  const [isCreatingSet, setIsCreatingSet] = useState(false);
  const [createSetStatus, setCreateSetStatus] = useState("");
  const [editingSet, setEditingSet] = useState<WorkoutSet | null>(null);
  const [editExerciseName, setEditExerciseName] = useState("");
  const [editReps, setEditReps] = useState("");
  const [editWeightKg, setEditWeightKg] = useState("");
  const [isUpdatingSet, setIsUpdatingSet] = useState(false);
  const [updateSetStatus, setUpdateSetStatus] = useState("");

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

    if (isLoadingDetail || isCreatingSet || isUpdatingSet) {
      return;
    }

    setDetailStatus("Loading...");
    setIsLoadingDetail(true);
    setSelectedWorkout(null);
    setWorkoutSets([]);
    setExerciseName("");
    setReps("");
    setWeightKg("");
    setCreateSetStatus("");
    setEditingSet(null);
    setUpdateSetStatus("");

    try {
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error(`Response.status: ${response.status}`);
      }
      const result: Workout = await response.json();
      const sets = await fetchWorkoutSets(id);
      setSelectedWorkout(result);
      setDetailStatus("Loaded");
      setWorkoutSets(sets);
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

  async function handleCreateWorkoutSet() {
    if (!selectedWorkout || isCreatingSet) {
      return;
    }
    setIsCreatingSet(true);
    setCreateSetStatus("Saving...");
    try {
      const createdSet = await createWorkoutSet(selectedWorkout.id, {
        exercise_name: exerciseName.trim(),
        reps: Number(reps),
        weight_kg: Number(weightKg),
      });
      setWorkoutSets((currentSets) => [...currentSets, createdSet]);
      setExerciseName("");
      setReps("");
      setWeightKg("");
      setCreateSetStatus("Set added");
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      setCreateSetStatus(`Error: ${message}`);
    } finally {
      setIsCreatingSet(false);
    }
  }

  async function handleUpdateWorkoutSet() {
    if (!editingSet || isUpdatingSet || isLoadingDetail) {
      return;
    }
    setIsUpdatingSet(true);
    setUpdateSetStatus("Editing...");

    try {
      const updatedSet = await updateWorkoutSet(editingSet.workout_id, editingSet.id, {
        exercise_name: editExerciseName.trim(),
        reps: Number(editReps),
        weight_kg: Number(editWeightKg),
        expected_revision: editingSet.revision,
      });
      setWorkoutSets((currentSets) =>
        currentSets.map((set) => (set.id === updatedSet.id ? updatedSet : set)),
      );
      setUpdateSetStatus("Set updated");
      setEditingSet(null);
    } catch (error) {
      if (error instanceof ApiError && error.status === 409) {
        setUpdateSetStatus("This set has changed. Your edits have been kept.");
      } else {
        const message = error instanceof Error ? error.message : "Unknown error";
        setUpdateSetStatus(`Error: ${message}`);
      }
    } finally {
      setIsUpdatingSet(false);
    }
  }

  function handleEditButton(workoutSet: WorkoutSet) {
    setEditingSet(workoutSet);
    setEditExerciseName(workoutSet.exercise_name);
    setEditReps(String(workoutSet.reps));
    setEditWeightKg(String(workoutSet.weight_kg));
    setUpdateSetStatus("");
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
              disabled={isLoadingDetail || isCreatingSet || isUpdatingSet}
              onClick={() => void loadWorkout(workout.id)}
            >
              {workout.name}
            </button>
          </li>
        ))}
      </ul>
      <p>{detailStatus}</p>
      {selectedWorkout ? (
        <section>
          <h2>{selectedWorkout.name}</h2>
          <p>Started: {new Date(selectedWorkout.started_at).toLocaleString()}</p>
          <p>Sets: {workoutSets.length}</p>
          {workoutSets.length > 0 ? (
            <ul>
              {workoutSets.map((workoutSet) => (
                <li key={workoutSet.id}>
                  <button
                    type="button"
                    disabled={isUpdatingSet}
                    onClick={() => handleEditButton(workoutSet)}
                  >
                    Edit
                  </button>
                  {workoutSet.exercise_name}: {workoutSet.reps} reps * {workoutSet.weight_kg} kg.
                </li>
              ))}
            </ul>
          ) : (
            <p>No sets logged yet.</p>
          )}
          <p>{updateSetStatus}</p>
          {editingSet && (
            <form
              onSubmit={(event) => {
                event.preventDefault();
                handleUpdateWorkoutSet();
              }}
            >
              <p>
                Editing: <label htmlFor="edit-exercise">Exercise</label>
                <input
                  id="edit-exercise"
                  type="text"
                  required
                  disabled={isUpdatingSet}
                  value={editExerciseName}
                  onChange={(event) => setEditExerciseName(event.target.value)}
                />{" "}
                — <label htmlFor="edit-reps">Reps</label>
                <input
                  id="edit-reps"
                  type="number"
                  required
                  disabled={isUpdatingSet}
                  min="1"
                  step="1"
                  value={editReps}
                  onChange={(event) => setEditReps(event.target.value)}
                />{" "}
                reps ×<label htmlFor="edit-weight">Weight</label>
                <input
                  id="edit-weight"
                  type="number"
                  required
                  disabled={isUpdatingSet}
                  min="0"
                  step="any"
                  value={editWeightKg}
                  onChange={(event) => setEditWeightKg(event.target.value)}
                />{" "}
                kg
              </p>
              <button type="button" disabled={isUpdatingSet} onClick={() => setEditingSet(null)}>
                Cancel
              </button>
              <button type="submit" disabled={isUpdatingSet}>
                Save
              </button>
            </form>
          )}
          <p>{createSetStatus}</p>
          <form
            onSubmit={(event) => {
              event.preventDefault();
              void handleCreateWorkoutSet();
            }}
          >
            <label htmlFor="set-exercise">Exercise</label>
            <input
              id="set-exercise"
              required
              disabled={isCreatingSet}
              type="text"
              value={exerciseName}
              onChange={(event) => setExerciseName(event.target.value)}
            />
            <label htmlFor="set-reps">Reps</label>
            <input
              id="set-reps"
              required
              disabled={isCreatingSet}
              type="number"
              value={reps}
              min="1"
              step="1"
              onChange={(event) => setReps(event.target.value)}
            />
            <label htmlFor="set-weight">Weight</label>
            <input
              id="set-weight"
              required
              disabled={isCreatingSet}
              type="number"
              min="0"
              step="any"
              value={weightKg}
              onChange={(event) => setWeightKg(event.target.value)}
            />
            <button type="submit" disabled={isCreatingSet}>
              Add set
            </button>
          </form>
        </section>
      ) : (
        <p>Select a workout</p>
      )}
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
