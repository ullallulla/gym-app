# Learning log

Last checkpoint: 2026-10-26.

## How we are working

- Learning through building is an explicit project goal. The user knows Python and TypeScript, wants to learn Go, and wants no Python in the application stack.
- Use small exercises with a concrete result. Explain the relevant concept, let the user attempt the code, review it, and verify the behavior together.
- Give hints and explanations, and provide complete examples when requested. Honor explicit requests to implement changes.
- `PLAN.md` contains the architecture and feature roadmap. This file records exercises, progress, and the next step.

## Completed so far

- Initialized a pnpm workspace with `apps/*` and `packages/*`.
- Added the Vite+ web starter, installed React and its TypeScript types, configured the React plugin, and enabled JSX.
- The user confirmed that the React frontend runs.
- Initialized the Go module `gym-app/api` under `services/api`.
- Implemented a Go HTTP server in `services/api/cmd/server/main.go`, listening on `127.0.0.1:8080`.
- Registered `GET /api/health`, returning HTTP 200 and JSON `{"status":"ok"}` with the JSON content type.
- The user confirmed that the backend runs.
- Configured the frontend development proxy to forward `/api` to `http://127.0.0.1:8080`.
- Separated `App` into `apps/web/src/App.tsx`; `main.tsx` mounts it into the HTML `#app` element.
- Removed unused starter assets from `src`. `public/favicon.svg` remains referenced by `index.html`.

Runtime success above was reported by the user. The code was inspected at this checkpoint; no builds or automated tests were run during the checkpoint.

## Current exercise: connect React to Go

Status: unfinished. The user stopped for the night before implementing the API button.

Current code state:

- `App.tsx` renders only the `Gym app` heading.
- `main.tsx` still has an unused `useState` import. Move that import into `App.tsx` when adding state.
- The development proxy is configured, but an end-to-end frontend request has not yet been confirmed.

Resume with these steps:

1. Import `useState` in `App.tsx` and remove the unused import from `main.tsx`.
2. Add `apiStatus` state, initially `"Not checked"`.
3. Write an asynchronous `checkApi` handler that sets `"Checking..."`, calls `fetch("/api/health")`, checks `response.ok`, reads the JSON, and displays its `status` value.
4. Handle failures with `try/catch` and display a useful error message. Remember that `fetch` does not reject automatically for HTTP error statuses.
5. Render a `Check API` button with `onClick={checkApi}` and a paragraph showing the current status.
6. Keep both servers running and restart the frontend if needed to apply its proxy configuration.

Completion checks:

- Clicking the button displays `ok` after the request completes.
- Stopping the Go server and clicking again displays an error.
- Restarting Go and retrying recovers successfully.
- The web production build succeeds after completing the exercise.

## Commands for resuming

Run these from the project root in separate terminals.

Go API:

```bash
go -C services/api run ./cmd/server
```

React website:

```bash
pnpm --filter web run dev
```

Check the API directly:

```bash
curl -i http://localhost:8080/api/health
```

Check the web production build after finishing the exercise:

```bash
pnpm --filter web run build
```

## Concepts discussed

- pnpm workspaces, root `node_modules/.pnpm`, dependency links, and keeping dependencies out of Git.
- Go module names and where `go mod init` creates `go.mod`.
- React components, JSX/TSX, and the roles of `main.tsx` and `App.tsx`.
- Go's `main` package, entry point, HTTP router, response writer, and response headers.
- `json.Encoder` writes JSON to a writer; `json.Marshal` returns JSON bytes. The health handler uses `Encoder` for convenience.
- Development proxying and the planned use of React state, events, and `fetch` in the current exercise.

## Later milestones

After completing the frontend-to-backend connection, introduce PostgreSQL and migrations, then build a small persistent workout flow. Clerk, TanStack Router/Query, dashboards, food logging, Expo/SQLite synchronization, and encrypted S3 backups are still future work described in `PLAN.md`.
