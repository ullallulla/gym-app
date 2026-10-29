# Learning checkpoint

Updated: 2026-10-29.

## Current progress

- Working: Go health and workout create/list/detail endpoints; React workout creation, list, and detail views.
- Workout-set logging is in progress. Migration 3, request/response types, validation, rejection tests, and owner-filtered parent lookup are implemented.
- `POST /api/workouts/{id}/sets` returns 400 for invalid input, 404 for a missing parent, 500 on lookup failure, and temporary 501 for a valid request. It saves no set yet.
- All assigned exercises are complete. No exercise is currently assigned.

## Next step

On “continue,” guide transaction setup in `createWorkoutSet`: begin with the existing context, handle errors, defer rollback, and move the parent lookup to `tx.QueryRow`. Reuse `err` with `=` after `tx, err :=`. Keep 501 until persistence is ready.

Then add counter allocation, set insertion, a full-snapshot change event, and commit before returning 201. Set UI follows later; roadmap: `PLAN.md`.

## Essential constraints

- Ownership is server-selected (`local-dev` until authentication).
- Allocate the per-user counter first, before parent lookup/entity writes; keep counter, entity, and event in one transaction. Follow `createWorkout` as the existing pattern.

## Latest verification

- Agent review: parent lookup is correct; Go formatting and `go -C services/api test ./... -v -count=1` passed (2026-10-29). Tests cover health, input validation, and request rejection, without database integration.
- User reports: migration/constraint exercises and lookup/outage/recovery checks passed. Latest screenshot shows the expected 501 placeholder. Migration Down remains untested.

## Startup

From the root, use separate terminals. User shell: Fish.

```fish
sudo docker compose up --wait db
set -gx DATABASE_URL 'postgres://gym_app:gym_app_dev@127.0.0.1:5432/gym_app?sslmode=disable'
go -C services/api run ./cmd/server
```

Web: from `apps/web`, run `vp run dev`.
