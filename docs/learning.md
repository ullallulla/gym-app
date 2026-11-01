# Learning checkpoint

Updated: 2026-11-01.

## Current progress

- Go health, workout create/list/detail, and set create/list/update APIs work. React creation/listing/recovery, PUT client, and edit form are approved. Save currently only prevents default submission.

## Active exercise — Save an edited set

Resume in `apps/web/src/App.tsx`: implement async `handleUpdateWorkoutSet` following the create handler. Await the approved `updateWorkoutSet` helper with `editingSet.workout_id`, `editingSet.id`, and a request containing trimmed draft name, numeric reps/weight, and `expected_revision` from `editingSet.revision`. On success, replace the matching row with the returned set and close editing. On failure, preserve drafts and show an error. Use saving/status state, `try/catch/finally`, duplicate-save and workout-loading guards, and disable edit/navigation controls while saving. Call the handler from the edit form's submission. Give clear hints; full solutions only when asked.

Then add explicit conflict recovery and manual verification. Roadmap: `PLAN.md`.

## Essential constraints

- Ownership is server-selected (`local-dev` until authentication).
- Writes allocate the per-user counter before entity lookup/writes; counter, entity, and event share one transaction. Pattern: `createWorkout`.
- Edits must check the client's expected revision and increment the stored revision; stale edits must not overwrite newer data.

## Latest verification

- Agent review (2026-11-01): form submission approved. Actual checks: `vp check`, TypeScript/build, and diff checks passed; `vp test` found no test files. Last Go suite passed on 2026-10-31.
- WSL verification (2026-11-01): Goose v3.29.0 applied all three migrations; version 3 and tables confirmed. `GET /api/workouts` returned HTTP 200 with `[]` after setup. Previous Linux records were not transferred.
- User evidence from the previous Linux database: PUT persistence, conflict rejection, ownership isolation, transactional logging, rollback, and recovery passed.

## Startup

From the root, use separate terminals. User shell: Fish.

```fish
sudo docker compose up --wait db
set -gx DATABASE_URL 'postgres://gym_app:gym_app_dev@127.0.0.1:5432/gym_app?sslmode=disable'
go run -tags='no_clickhouse,no_libsql,no_mssql,no_mysql,no_sqlite3,no_vertica,no_ydb' github.com/pressly/goose/v3/cmd/goose@v3.29.0 -dir services/api/migrations postgres "$DATABASE_URL" up
go -C services/api run ./cmd/server
```

Web: from `apps/web`, `vp run dev` starts the server; `vp run build` only builds files. The root has no build script.

WSL workaround: `vp run dev --host <WSL-IP> --port 5173 --strictPort` (IP: `hostname -I`). Direct IP verified from Windows; localhost forwarding fails. Stop any existing server first.
