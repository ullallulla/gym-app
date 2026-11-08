# Learning checkpoint

Updated: 2026-11-02.

## Current progress

- Milestone 1 editing/conflict recovery is complete.
- Milestone 2 workout-create retries are verified on the server; schema version 4 is applied. The web caller still needs stable IDs and payload retention. Architecture and retry contracts are in `PLAN.md`.

## Active exercise — Type the web create request

This exercise has not started. `createWorkout` in `apps/web/src/App.tsx` still sends no ID and generates a timestamp per attempt.

Next: define `CreateWorkoutRequest` with optional string `id`, string `name`, and string `started_at`. Inside the existing `try`, extract the POST body into a typed `request` object before `fetch`; send `JSON.stringify(request)`.

Then generate one UUID and retain the complete pending body across retries until a successful response. UUID generation must work over the plain HTTP WSL IP. Give one small exercise at a time.

## Essential constraints

- Ownership is server-selected (`local-dev` until authentication).
- Allocate the per-user counter before entity lookup/writes; counter, entity, and event share one transaction. Retries/conflicts roll back provisional allocation.
- Enforce expected revisions; preserve the current weight representation.
- Keep reviews brief with existing checks; broader database checks belong at completed feature boundaries.

## Latest verification

- Code review and actual checks (2026-11-02): server flow reviewed; 19 actual-handler/PostgreSQL requests passed, covering conflicts, ownership, normalization, later edits, concurrency, and counter/event preservation. Temporary schema removed. Go tests, formatting, diff checks, and isolated migration Up/Down verification passed.
- User results: Goose confirmed migration 4; milestone 1 browser persistence/error/conflict recovery and WSL listing verified.

## Startup

From root, separate terminals; user shell: Fish.

```fish
sudo docker compose up --wait db
set -gx DATABASE_URL 'postgres://gym_app:gym_app_dev@127.0.0.1:5432/gym_app?sslmode=disable'
goose -dir services/api/migrations postgres "$DATABASE_URL" up
go -C services/api run ./cmd/server
```

Web: `apps/web`, `vp run dev`; build: `vp run build`. Root has no build script.

WSL: `vp run dev --host <WSL-IP> --port 5173 --strictPort`; IP: `hostname -I`. Windows direct IP works; localhost forwarding fails. Stop existing server first.

Restart the Go API after Go source changes before browser verification.
