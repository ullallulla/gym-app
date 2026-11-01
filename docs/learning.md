# Learning checkpoint

Updated: 2026-11-01.

## Current progress

- Go workout/set APIs and React creation/listing/edit saving passed code review, including guards, disabled controls, status resets, `ApiError`, and HTTP 409 messaging. Drafts survive failed saves. Browser verification remains.
- Milestone 2 remains planned in `PLAN.md`. Schema has migrations 1–3. Counter/revision details remain under review in `docs/proposed-changes.md`.

## Active exercise — Recover from revision conflicts

Next small step: track conflicts explicitly in `App.tsx` to control recovery UI. Set conflict state on HTTP 409; clear on new edit, workout load, successful save, or cancel. Keep drafts and the editing revision unchanged. Give hints; full solutions only when asked.

Then add a user-triggered reload of the latest record, preserving drafts for review before a manual resave. Avoid `loadWorkout` for this reload: it clears editing/drafts. No automatic retry or silent revision replacement. Guard actions during recovery. Verify persistence, failed saves, and stale edits before milestone 2.

## Essential constraints

- Ownership is server-selected (`local-dev` until authentication).
- Allocate the per-user counter before entity lookup/writes; counter, entity, and event share one transaction. Pattern: `createWorkout`.
- Check expected revisions; increment stored revisions. Stale edits must not overwrite newer data.
- Keep the current weight representation during this exercise.

## Latest verification

- Code review (2026-11-01): conflict message/fallback approved; drafts and revision preserved. Actual checks: `vp check`, TypeScript/build, and diff checks passed; `vp test` found no tests. Earlier `vp install` passed. No browser verification. Go suite last passed 2026-10-31.
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
