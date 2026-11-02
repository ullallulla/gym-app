# Learning checkpoint

Updated: 2026-11-02.

## Current progress

- Milestone 1 editing/conflict recovery is complete, with code review and user-reported browser verification.
- Milestone 2 starts with client-generated UUIDs and create retries, as specified in `PLAN.md`. Schema has migrations 1–3.

## Active exercise — Understand create retries

Next small step: trace `createWorkout` in `services/api/cmd/server/main.go`. Consider a successful transaction whose HTTP response never reaches the client, followed by the same POST again. Explain what the current INSERT and database UUID default do to workout rows and change-log events, and which identity must remain stable across attempts.

Then define same-ID retry/payload-mismatch behavior and ownership responses before changing SQL. Start implementation with workouts. Give one exercise at a time; full examples when asked. No milestone 2 code changes yet.

## Essential constraints

- Ownership is server-selected (`local-dev` until authentication).
- Allocate the per-user counter before entity lookup/writes; counter, entity, and event share one transaction. Pattern: `createWorkout`.
- Check expected revisions; increment stored revisions. Stale edits must not overwrite newer data.
- Keep the current weight representation until an explicit precision/representation decision.

## Latest verification

- Code review (2026-11-02): conflict recovery/guards complete. Actual checks: `vp install`, web `vp check`, TypeScript/build, and environment doctor passed. Root `vp check` reports formatting in unchanged `PLAN.md`/`package.json`; `vp test` found no tests (exit 1). Diff checks passed. Go suite last passed 2026-10-31.
- User-reported browser results (2026-11-02): persistence, failed saves/retries, two-tab conflicts, and failed reloads/retries all work as intended, including preserved drafts and manual saving. Earlier WSL results: migrations through version 3 and empty workout listing verified.

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
