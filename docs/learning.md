# Learning checkpoint

Updated: 2026-10-30.

## Current progress

- Working: Go health and workout create/list/detail endpoints; React workout creation, list, and detail views.
- Workout-set migration 3, types, validation, rejection tests, and Go create/list endpoints are complete. Owner-filtered reads and transactional change logging are implemented and reviewed.
- React set listing, creation, persistence, validation-error recovery, slow-save controls, and server-outage recovery work in the browser. Set creation is complete. Go editing request validation and its focused tests are implemented and reviewed; the HTTP handler is next.

## Active exercise

In `main.go`, add `app.updateWorkoutSet` and register `PUT /api/workouts/{id}/sets/{setID}`. Parse both path values as `pgtype.UUID`. Return 400/return for an invalid workout ID or workout-set ID, with distinct messages. After both IDs validate, return 501 while building the endpoint. Follow `createWorkoutSet`'s UUID validation pattern. Review the user's attempt next.

Then decode/validate the update JSON and implement ownership/revision checks plus transactional change logging. Connect editing to React afterward. Roadmap: `PLAN.md`.

## Essential constraints

- Ownership is server-selected (`local-dev` until authentication).
- For writes, allocate the per-user counter first, before parent lookup/entity writes; keep counter, entity, and event in one transaction. Follow `createWorkout` as the existing pattern.
- Edits must check the client's expected revision and increment the stored revision; stale edits must not overwrite newer data.

## Latest verification

- Agent review (2026-10-30): update validation tests correctly assert success/name trimming and each expected error; case names are clear. Checks run: `gofmt -d` and `git diff --check` clean; `go -C services/api test ./... -count=1` passes. Last frontend checks/build passed; no frontend tests or database integration tests exist.
- User evidence: browser creation/persistence, listing, validation recovery, and slow-save controls pass. User reports outage error handling and retry also work. GET API cases and POST transactional rollback/recovery passed. Migration Down remains untested.

## Startup

From the root, use separate terminals. User shell: Fish.

```fish
sudo docker compose up --wait db
set -gx DATABASE_URL 'postgres://gym_app:gym_app_dev@127.0.0.1:5432/gym_app?sslmode=disable'
go -C services/api run ./cmd/server
```

Web: from `apps/web`, run `vp run dev`.
