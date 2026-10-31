# Learning checkpoint

Updated: 2026-10-31.

## Current progress

- Working: Go health, workout create/list/detail, and workout-set create/list APIs with transactional change logging. React workout views and set creation/listing/recovery are complete.
- Go PUT set editing is implemented, reviewed, and manually verified against PostgreSQL. React editing is next.

## Active exercise — TypeScript update API client

In `apps/web/src/workoutSets.ts`, export `UpdateWorkoutSetRequest` combining `CreateWorkoutSetRequest` with numeric `expected_revision`. Add async `updateWorkoutSet(workoutId, setId, request): Promise<WorkoutSet>` following the create helper: PUT to the set URL, JSON Content-Type/body, reject non-OK responses including 409, and return the parsed set. Review the user's implementation before UI wiring.

Then add edit state/form, saving/error handling, conflict recovery, and replace the saved row with the API result. Roadmap: `PLAN.md`.

## Essential constraints

- Ownership is server-selected (`local-dev` until authentication).
- For writes, allocate the per-user counter first, before parent lookup/entity writes; keep counter, entity, and event in one transaction. Follow `createWorkout` as the existing pattern.
- Edits must check the client's expected revision and increment the stored revision; stale edits must not overwrite newer data.

## Latest verification

- Agent review (2026-10-31): completed PUT handler and fixture SQL reviewed. Last code checks: `gofmt -d` and `git diff --check` clean; Go suite passed (`go -C services/api test ./... -count=1`). Tests cover input rejection. Last frontend checks/build passed; no frontend tests exist.
- User evidence: update persistence, conflict rejection, workout/owner isolation, event consistency, forced event-insert rollback, and recovery pass. Test constraint was removed. Recovery leaves the local set at 1 rep/5 kg/revision 3, counter 20, two update events. Ownership fixture remains at 8 reps/20 kg/revision 1. Earlier POST rollback/recovery passed.

## Startup

From the root, use separate terminals. User shell: Fish.

```fish
sudo docker compose up --wait db
set -gx DATABASE_URL 'postgres://gym_app:gym_app_dev@127.0.0.1:5432/gym_app?sslmode=disable'
go -C services/api run ./cmd/server
```

Web: from `apps/web`, run `vp run dev`.
