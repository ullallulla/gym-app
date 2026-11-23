# Learning checkpoint

Updated: 2026-11-23.

## Current progress

- Milestone 1 editing/conflict recovery is complete.
- Milestone 2 workout-create retries are verified on the server; schema version 4 is applied. The web caller still needs stable IDs and payload retention. Contracts and review decisions are in `PLAN.md`.

## Active exercise — Type the web create request

Not started. Define `CreateWorkoutRequest` (optional `id`, `name`, `started_at` strings), build a typed `request` inside the existing `try` in `createWorkout`, and send `JSON.stringify(request)`.

Then generate one UUID and retain the complete pending body until a successful response. `crypto.randomUUID()` needs a secure context, so it fails over the plain HTTP WSL IP.

## Queued fix exercises

One at a time. Do 3 before UUID retention: retries must tell 409 from network/5xx failures.

1. Set `"strict": true` in `apps/web/tsconfig.json` (`tsc --strict` already passes).
2. Import `defineConfig` from `vite-plus`; enable `lint.options.typeAware` and `typeCheck`; fix the floating `handleUpdateWorkoutSet()` promise.
3. Throw `ApiError` with status and server message from every fetch helper.
4. Conflict reload replaces the whole set list; Cancel clears the update status.
5. Go: log `lookupErr`; use one "Could not create workout" message; `go mod tidy`; return scanned timestamps in UTC.
6. Migration 5: rename constraint `original_name`; require nonblank `workouts.name`; check `change_log.entity_type`. Never edit applied migrations.
7. Fixture: build change-log payloads explicitly; comment that null `original_*` intentionally tests the unverifiable-retry 409.
8. Automate PostgreSQL handler checks, gated by `TEST_DATABASE_URL`.

## Essential constraints

- Ownership is server-selected (`local-dev` until authentication).
- Allocate the per-user counter before entity lookup/writes; counter, entity, and event share one transaction. Retries/conflicts roll back provisional allocation.
- Enforce expected revisions; preserve the current weight representation.
- Keep reviews brief; broader database checks belong at completed feature boundaries.

## Latest verification

- Actual checks: 2026-11-02, 19 handler/PostgreSQL requests passed (conflicts, ownership, normalization, edits, concurrency, counter/event preservation); migration 4 Up/Down verified. 2026-11-23, `go vet`, `go test`, `tsc`, and `vp check` pass.
- User results: Goose applied migration 4; milestone 1 browser flows and WSL listing verified.

## Startup

From root, separate terminals; user shell: Fish.

```fish
sudo docker compose up --wait db
set -gx DATABASE_URL 'postgres://gym_app:gym_app_dev@127.0.0.1:5432/gym_app?sslmode=disable'
goose -dir services/api/migrations postgres "$DATABASE_URL" up
go -C services/api run ./cmd/server
```

Web: `apps/web`, `vp run dev`. WSL: `vp run dev --host <WSL-IP> --port 5173 --strictPort` (IP: `hostname -I`); Windows direct IP works, localhost forwarding fails. Restart the Go API after Go changes.
