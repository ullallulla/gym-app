# Learning checkpoint

Updated: 2026-11-01.

## Current progress

- Go workout/set APIs and React creation/listing/edit saving passed code review, including guards, disabled controls, status resets, `ApiError`, HTTP 409 messaging, and conflict state. Drafts survive failed saves. Browser verification remains.
- Milestone 2 remains planned in `PLAN.md`. Schema has migrations 1–3. Counter/revision details remain under review in `docs/proposed-changes.md`.

## Active exercise — Recover from revision conflicts

Next small step: teach a conflict-reload handler in `App.tsx`, using `fetchWorkoutSets` and the editing set's ID to find the latest record. On an explicit reload, update that row and the editing snapshot/revision while preserving all draft inputs; clear the conflict and prompt review before manual saving. On fetch failure or a missing record, preserve drafts/conflict and show an error. Use loading state, guards, and `try/catch/finally`. Give hints; full solutions only when asked.

Then connect a Reload latest button shown during conflicts and disable conflicting controls during recovery. Avoid `loadWorkout`: it clears editing/drafts. No automatic retry or silent revision replacement. Verify persistence, failed saves, and stale edits before milestone 2.

## Essential constraints

- Ownership is server-selected (`local-dev` until authentication).
- Allocate the per-user counter before entity lookup/writes; counter, entity, and event share one transaction. Pattern: `createWorkout`.
- Check expected revisions; increment stored revisions. Stale edits must not overwrite newer data.
- Keep the current weight representation during this exercise.

## Latest verification

- Code review (2026-11-01): conflict guard, successful-save reset, and Cancel fixes approved. Actual checks: `vp check`, TypeScript/build, and diff checks passed; `vp test` found no tests. Earlier `vp install` passed. No browser verification. Go suite last passed 2026-10-31.
- User-reported WSL results: migrations through version 3 and empty workout listing verified.

## Startup

From root, separate terminals; user shell: Fish.

```fish
sudo docker compose up --wait db
set -gx DATABASE_URL 'postgres://gym_app:gym_app_dev@127.0.0.1:5432/gym_app?sslmode=disable'
go run -tags='no_clickhouse,no_libsql,no_mssql,no_mysql,no_sqlite3,no_vertica,no_ydb' github.com/pressly/goose/v3/cmd/goose@v3.29.0 -dir services/api/migrations postgres "$DATABASE_URL" up
go -C services/api run ./cmd/server
```

Web: `apps/web`, `vp run dev`; build: `vp run build`. Root has no build script.

WSL: `vp run dev --host <WSL-IP> --port 5173 --strictPort`; IP: `hostname -I`. Windows direct IP works; localhost forwarding fails. Stop existing server first.
