# Learning checkpoint

Updated: 2026-10-30.

## Current progress

- Working: Go health and workout create/list/detail endpoints; React workout creation, list, and detail views.
- Workout-set migration 3, types, validation, rejection tests, and Go create/list endpoints are complete. Owner-filtered reads and transactional change logging are implemented and reviewed.
- React set listing, creation, persistence, validation-error recovery, slow-save controls, and server-outage recovery work in the browser. TypeScript helpers and async submit flow are reviewed. The set-creation form is complete.

## Next step

The set-creation form is ready to commit: `Add workout set creation form`. No active coding exercise. Choose the next small feature using `PLAN.md` after this milestone.

## Essential constraints

- Ownership is server-selected (`local-dev` until authentication).
- For writes, allocate the per-user counter first, before parent lookup/entity writes; keep counter, entity, and event in one transaction. Follow `createWorkout` as the existing pattern.

## Latest verification

- Agent review (2026-10-30): submit flow and guards are correct. Saving is cleared in `finally`; an immediate outage error can make the disabled transition too brief to notice. Last checks: `vp install` up to date; `vp check` and `vp run build` pass; frontend test runner found no tests. Last Go tests passed without database integration.
- User evidence: browser creation/persistence, listing, validation recovery, and slow-save controls pass. User reports outage error handling and retry also work. GET API cases and POST transactional rollback/recovery passed. Migration Down remains untested.

## Startup

From the root, use separate terminals. User shell: Fish.

```fish
sudo docker compose up --wait db
set -gx DATABASE_URL 'postgres://gym_app:gym_app_dev@127.0.0.1:5432/gym_app?sslmode=disable'
go -C services/api run ./cmd/server
```

Web: from `apps/web`, run `vp run dev`.
