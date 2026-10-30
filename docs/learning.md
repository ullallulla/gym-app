# Learning checkpoint

Updated: 2026-10-30.

## Current progress

- Working: Go health and workout create/list/detail endpoints; React workout creation, list, and detail views.
- Workout-set logging is in progress. Migration 3, request/response types, validation, rejection tests, and owner-filtered parent lookup are implemented.
- `POST /api/workouts/{id}/sets` validates input, allocates the counter, checks the parent, inserts/scans a set and snapshot event, commits, then returns 201 JSON. Live persistence is user-confirmed.
- `GET /api/workouts/{id}/sets` is implemented and reviewed: UUID validation, owner-filtered parent lookup, query, scanning, iteration error handling, and JSON response. Populated, empty, and missing-workout cases are confirmed live. React integration is next.

## Active exercise

Create `apps/web/src/workoutSets.ts`. Export a `WorkoutSet` type matching the Go response (IDs/timestamps as strings; reps, weight, revision as numbers) and an async `fetchWorkoutSets(workoutId: string): Promise<WorkoutSet[]>` helper. Fetch `/api/workouts/${workoutId}/sets`, throw on non-OK responses, and return parsed JSON. Use the fetch/error pattern in `App.tsx` as a reference. Review the user's attempt next.

Then connect loading and displaying sets to workout selection in React. Add set creation afterward; roadmap: `PLAN.md`.

## Essential constraints

- Ownership is server-selected (`local-dev` until authentication).
- For writes, allocate the per-user counter first, before parent lookup/entity writes; keep counter, entity, and event in one transaction. Follow `createWorkout` as the existing pattern.

## Latest verification

- Agent review (2026-10-30): completed GET sets handler is correct. Checks run: `gofmt -d` and `git diff --check` clean; `go -C services/api test ./... -count=1` passed. Tests cover health/input rejection without database integration.
- User evidence: populated GET returns 200 JSON with sets in ascending creation order. User confirms empty-workout GET returns 200 with `[]` and missing-workout GET returns 404. POST persistence and rollback/recovery also passed (set/event/counter all +1 on success, unchanged on event failure). Migration Down remains untested.

## Startup

From the root, use separate terminals. User shell: Fish.

```fish
sudo docker compose up --wait db
set -gx DATABASE_URL 'postgres://gym_app:gym_app_dev@127.0.0.1:5432/gym_app?sslmode=disable'
go -C services/api run ./cmd/server
```

Web: from `apps/web`, run `vp run dev`.
