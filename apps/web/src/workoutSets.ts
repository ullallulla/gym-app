export type WorkoutSet = {
  id: string;
  user_id: string;
  workout_id: string;
  exercise_name: string;
  reps: number;
  weight_kg: number;
  created_at: string;
  revision: number;
};

export type CreateWorkoutSetRequest = {
  exercise_name: string;
  reps: number;
  weight_kg: number;
};

export async function fetchWorkoutSets(workoutId: string): Promise<WorkoutSet[]> {
  const url = `/api/workouts/${workoutId}/sets`;

  const response = await fetch(url);

  if (!response.ok) {
    throw new Error(`Response status: ${response.status}`);
  }
  const result = await response.json();
  return result;
}

export async function createWorkoutSet(
  workoutId: string,
  request: CreateWorkoutSetRequest,
): Promise<WorkoutSet> {
  const url = `/api/workouts/${workoutId}/sets`;
  const response = await fetch(url, {
    method: "POST",
    headers: { "Content-type": "application/json" },
    body: JSON.stringify(request),
  });
  if (!response.ok) {
    throw new Error(`Response status: ${response.status}`);
  }

  const workoutSet: WorkoutSet = await response.json();

  return workoutSet;
}
