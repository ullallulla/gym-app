# Learning checkpoint

Updated: 2026-11-02.

## Current progress

- Milestone 1 editing/conflict recovery is complete, with code review and user-reported browser verification.
- Milestone 2 workout-create retries are verified on the server. Schema version 4 is applied; the web caller still needs stable IDs and payload retention.

## Active exercise — Type the web create request

The actual server handler passed isolated API/database checks, including JSON responses, ownership, unknown originals, normalization, later edits, concurrency, and counter/event preservation. The web `createWorkout` in `apps/web/src/App.tsx` still sends no ID and generates a timestamp on each attempt.

Retry rules are in `PLAN.md`. Next exercise: define a TypeScript create-request type with optional `id`, plus `name` and `started_at`; extract the inline POST body into a typed `request` object before `fetch` and send `JSON.stringify(request)`. Then add ID generation and retain the full pending body across retries. Account for development over a plain HTTP WSL IP when choosing UUID generation. Give one small exercise at a time; examples when asked.

## Essential constraints

- Ownership is server-selected (`local-dev` until authentication).
- Allocate the per-user counter before entity lookup/writes; counter, entity, and event share one transaction. Pattern: `createWorkout`.
- Check expected revisions; increment stored revisions. Stale edits must not overwrite newer data.
- Keep the current weight representation until an explicit precision/representation decision.
- Keep incremental reviews brief with existing Go checks; reserve broader database checks for the completed retry flow.

## Latest verification

- Actual checks (2026-11-02): 19 requests through the actual handler and PostgreSQL passed, including simultaneous identical creates and complete data/counter/event comparisons for retries/conflicts. The temporary schema was removed. Existing Go tests, formatting, and diff checks passed; migration 4 previously passed its isolated Up/Down verification.
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

Restart the Go API after Go source changes before browser verification.
