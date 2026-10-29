# Learning checkpoint

Updated: 2026-10-29.

## Current progress

- Working: Go health and workout create/list/detail endpoints; React workout creation, list, and detail views.
- Workout-set logging is in progress. Migration 3, request/response types, validation, rejection tests, and owner-filtered parent lookup are implemented.
- `POST /api/workouts/{id}/sets` validates input, allocates the counter, checks the parent, inserts/scans a set and snapshot event, commits, then returns 201 JSON. Live persistence is user-confirmed.
- Handler review and successful-creation verification are complete. No active exercise.

## Next step

On “continue,” guide a deliberate change-event insert failure: verify HTTP 500 and unchanged set/event counts and counter. Restore operation `create` afterward and verify successful creation again.

Set listing/UI follow later; roadmap: `PLAN.md`.

## Essential constraints

- Ownership is server-selected (`local-dev` until authentication).
- Allocate the per-user counter first, before parent lookup/entity writes; keep counter, entity, and event in one transaction. Follow `createWorkout` as the existing pattern.

## Latest verification

- Agent review (2026-10-29): duplicate rollback removed; the early defer remains. Formatting is clean. `go -C services/api test ./... -count=1` passed before this one-line cleanup. Tests cover health, input validation, and request rejection, without database integration.
- User reports: live POST returns 201 with the stored set, matching snapshot event, and exactly +1 counter increment. Prior migration/constraint and lookup/outage/recovery exercises passed. Migration Down remains untested.

## Startup

From the root, use separate terminals. User shell: Fish.

```fish
sudo docker compose up --wait db
set -gx DATABASE_URL 'postgres://gym_app:gym_app_dev@127.0.0.1:5432/gym_app?sslmode=disable'
go -C services/api run ./cmd/server
```

Web: from `apps/web`, run `vp run dev`.
