# Learning checkpoint

Updated: 2026-10-30.

## Current progress

- Working: Go health and workout create/list/detail endpoints; React workout creation, list, and detail views.
- Workout-set migration 3, types, validation, rejection tests, and Go create/list endpoints are complete. Owner-filtered reads and transactional change logging are implemented and reviewed.
- React set loading, count, list, and empty state work in the browser. TypeScript GET/POST helpers are reviewed; the set-entry UI is next.

## Active exercise

In `App.tsx`, add labeled controlled inputs inside the selected workout section for exercise name, reps, and weight in kilograms. Use string states `exerciseName`, `reps`, and `weightKg`, initially empty; input `value` and `onChange` should use those states. Reps: number, min 1, step 1. Weight: number, min 0, step any. Reset these fields when `loadWorkout` starts after its guard. Review the inputs before wiring submission.

Then add form submission, convert numeric strings, call the POST helper, and update displayed sets. Roadmap: `PLAN.md`.

## Essential constraints

- Ownership is server-selected (`local-dev` until authentication).
- For writes, allocate the per-user counter first, before parent lookup/entity writes; keep counter, entity, and event in one transaction. Follow `createWorkout` as the existing pattern.

## Latest verification

- Agent review (2026-10-30): POST helper and error propagation are correct. Checks run: `vp install` up to date; `vp check` and `vp run build` pass; `vp test --passWithNoTests` finds no frontend tests. Last Go tests passed without database integration.
- User evidence: browser list/empty state and switching work. GET returns ordered sets, `[]` for empty, and 404 for missing workouts. POST persistence and rollback/recovery passed (set/event/counter all +1 on success, unchanged on event failure). Migration Down remains untested.

## Startup

From the root, use separate terminals. User shell: Fish.

```fish
sudo docker compose up --wait db
set -gx DATABASE_URL 'postgres://gym_app:gym_app_dev@127.0.0.1:5432/gym_app?sslmode=disable'
go -C services/api run ./cmd/server
```

Web: from `apps/web`, run `vp run dev`.
