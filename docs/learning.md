# Learning checkpoint

Updated: 2026-11-02.

## Current progress

- Milestone 1 editing/conflict recovery is complete, with code review and user-reported browser verification.
- Milestone 2 starts with client-generated UUIDs and create retries, as specified in `PLAN.md`. User's Goose output confirms schema version 4.

## Active exercise — Store original values on creation

Optional ID validation in `services/api/cmd/server/main.go` and both rejection tests pass review. Empty/malformed IDs return HTTP 400 before database work; omitted/null IDs skip validation. SQL still generates fresh UUIDs, so lost-response retries can duplicate workouts/events.

Retry and storage rules are in `PLAN.md`. Migration 4 is applied. Next exercise: extend `createWorkout`'s INSERT with `original_name` and `original_started_at`, using the existing normalized name/time values. Reuse placeholders `$2` and `$3`; the existing three Go arguments and public RETURNING/Scan fields already suffice. Review this change before adding client-ID insertion and retry handling. Give one small exercise at a time; examples when asked.

## Essential constraints

- Ownership is server-selected (`local-dev` until authentication).
- Allocate the per-user counter before entity lookup/writes; counter, entity, and event share one transaction. Pattern: `createWorkout`.
- Check expected revisions; increment stored revisions. Stale edits must not overwrite newer data.
- Keep the current weight representation until an explicit precision/representation decision.

## Latest verification

- Actual checks (2026-11-02): PostgreSQL 18.6 passed 36 fixtures plus Up/Down, paired-null rejection, repeat-backfill, and record/log/counter preservation checks in a rolled-back schema. Diff checks passed. Last Go tests and formatting checks passed earlier.
- User-reported results: Goose screenshot confirms migration 4 applied. Milestone 1 browser persistence, failed-save/reload recovery, and two-tab conflicts passed with preserved drafts/manual saving; WSL listing verified.

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
