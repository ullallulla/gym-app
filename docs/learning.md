# Learning checkpoint

Updated: 2026-11-01.

## Current progress

- Go workout/set APIs, React creation/listing, the PUT helper, and the edit form exist. `App.tsx` does not import or call `updateWorkoutSet`; edit submission only prevents default.
- Exercise library/renaming, blocks, ordering, completion, prefilling, and deletion safeguards are planned in `PLAN.md`, not implemented. Custom renames will update current and past workout labels. Current schema has migrations 1–3. Counter/revision details remain under review in `docs/proposed-changes.md`.

## Active exercise — Save an edited set

In `apps/web/src/App.tsx`, import `updateWorkoutSet` and implement async `handleUpdateWorkoutSet`, following the create handler. Submit the editing set's workout/record IDs, trimmed name, numeric reps/weight, and `expected_revision`. Replace the saved row and close editing on success; preserve drafts and show errors on failure. Add saving/status state, `try/catch/finally`, duplicate-save/loading guards, and disable conflicting controls while saving. Wire edit-form submission to the handler. Give hints; full solutions only when asked.

Next: explicit revision-conflict recovery preserving drafts, then verify persistence, failed saves, and stale edits before milestone 2.

## Essential constraints

- Ownership is server-selected (`local-dev` until authentication).
- Allocate the per-user counter before entity lookup/writes; counter, entity, and event share one transaction. Pattern: `createWorkout`.
- Check expected revisions; increment stored revisions. Stale edits must not overwrite newer data.
- Keep the current weight representation during this exercise.

## Latest verification

- Code inspection (2026-11-01): confirmed missing edit-save wiring and current schema. This documentation update: diff/word count checked; no runtime checks rerun.
- Prior recorded checks: `vp check` and TypeScript/build passed; `vp test` found no tests. Go suite last passed 2026-10-31.
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
