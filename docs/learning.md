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
- Connected React to Go through the development proxy. The user confirmed successful requests, visible errors when Go is stopped, and recovery after restarting Go. The production build passed during review.
- Added a PostgreSQL 18 Compose definition. The user installed Docker and Compose and reported that running `SELECT 1;` through the container's `psql` returned `1`.

Runtime success above was reported by the user. At the initial checkpoint, code was inspected without builds or automated tests. Checks from the resumed session are recorded below.

## Completed exercise: connect React to Go

Status: complete. The user confirmed successful requests, visible failures while Go is stopped, and recovery after restarting it. The production build passed during review, and the component now passes formatting and lint checks.

Current code state:

- `App.tsx` owns `apiStatus` state, initially `"Not checked"`, and renders the button and status paragraph. The unused import was removed from `main.tsx`.
- `checkApi` is inside `App`, fetches `/api/health`, checks `response.ok`, parses the JSON, and stores `result.status`.
- The relative URL uses the configured development proxy. The earlier absolute URL bypassed the proxy and triggered a browser CORS error.
- The handler sets `"Checking..."` before fetching. Its catch block narrows `error` with `instanceof Error`, supplies an unknown-error fallback, and displays the error through state.

Initial review and verification (2026-10-26):

- Inspected the current component, entry point, proxy, API handler, scripts, and TypeScript configuration.
- User-reported runtime result: the request now works. Failure and recovery behavior have not been confirmed.
- Ran `pnpm --filter web exec vp run build`: failed with `TS18046` at `App.tsx:17` because `error` is `unknown`.
- Ran `pnpm --filter web exec vp check`: failed on formatting in `package.json`, `src/App.tsx`, `src/main.tsx`, and `vite.config.ts`. This did not establish a passing lint result.
- Ran `pnpm --filter web exec vp test`: exited with code 1 because no test files exist.
- Initially only the local Vite+ CLI was accessible through `pnpm --filter web exec vp`; its `vp env doctor` command required the global CLI.
- Tooling follow-up: the user installed global Vite+. Sourcing `~/.config/vite-plus/env` makes it available in the agent's existing Bash environment. Verified `vp v1.4.0` and ran global `vp env doctor`: all checks passed, using managed Node.js v24.22.1 and pnpm v12.3.4.
- The required Vite+ install hit restricted registry access and left the local executable unavailable. Restored dependencies successfully with approved `pnpm install --frozen-lockfile`; no manifest or lockfile changes were needed.
- Application code was reviewed without edits.

Follow-up review (2026-10-26):

- Inspected the user's updated `App.tsx`: loading state, error narrowing, and visible errors are now implemented correctly for this exercise. Explained ternary expressions, `instanceof Error`, and the `unknown` catch variable.
- Ran `vp run build` from `apps/web`: TypeScript and the production build both passed.
- Ran `vp check src/App.tsx`: formatting still needs correction.
- Ran `vp check --no-fmt src/App.tsx`: passed with no warnings or lint errors.
- Did not repeat the test runner: no test files were added since its previous no-tests result. Browser failure and recovery behavior remain unverified.
- Application code was reviewed without edits.

Completion confirmation (2026-10-26):

- User-reported runtime results: stopping Go produces a visible error; restarting Go and retrying returns the status to `ok`. These browser checks were performed by the user.
- Inspected the latest `App.tsx`: the user has formatted it, and the reviewed request/loading/error behavior remains present.
- Ran `vp check src/App.tsx`: formatting and lint both passed.
- Formatting findings in other web files from the initial review have not been rechecked. No new builds or tests were needed for this confirmation.

Completion checks:

- Clicking the button displays `ok` after the request completes (success reported by the user).
- Stopping the Go server and clicking again displays an error (reported by the user).
- Restarting Go and retrying recovers successfully (reported by the user).
- The web production build succeeds (verified in the follow-up review).

## Completed exercise: run PostgreSQL locally

Status: complete. The user confirmed installing Docker and Compose, creating `compose.yaml`, and successfully running `SELECT 1;`.

- Inspected the local environment: CachyOS (Arch-based). Neither Docker, Podman, nor PostgreSQL executables are currently on this terminal's PATH; no existing Compose or database setup files were found in the repository.
- Recommended Docker Compose and compared it with native installation. The user then requested Docker/PostgreSQL installation instructions for Arch Linux.
- Provided installation steps using the Arch `docker` and `docker-compose` packages, the Docker systemd service, and a `hello-world` check. Docker commands use `sudo`; no Docker group membership change is required for these steps.
- Proposed a root `compose.yaml` using `postgres:18`, a `db` service, database/user `gym_app`, local development password `gym_app_dev`, a port binding at `127.0.0.1:5432`, a named volume mounted at `/var/lib/postgresql`, and a `pg_isready` health check.
- The user should create the Compose file, run `sudo docker compose up --wait`, then verify `SELECT 1;` through `sudo docker compose exec db psql -U gym_app -d gym_app -c "SELECT 1;"`.
- At the instruction checkpoint, no database software had been installed, no Compose file had been created, and no container or database query had been run by the agent.

Completion review (2026-10-26):

- User-reported runtime result: Docker and Compose were installed, the Compose file was created, and the container `psql` command returned `1` for `SELECT 1;`.
- Inspected the actual root `compose.yaml`: PostgreSQL 18, the expected database/user, loopback port binding, a named volume at the correct mount location, and the readiness check are present.
- Ran `docker compose version`: v5.6.1. Ran `docker compose config --quiet`: passed.
- Did not run a container query or independently verify database connectivity. No application code or Compose configuration was edited during this review.
- Discussed `docker compose stop` versus `down` and how named volumes preserve database data by default.

## Current exercise: connect Go to PostgreSQL

Status: not yet implemented. Start with a small configuration exercise before adding the database driver.

1. In the terminal used to run Go, export `DATABASE_URL='postgres://gym_app:gym_app_dev@127.0.0.1:5432/gym_app?sslmode=disable'`.
2. Add the standard-library `os` import in `services/api/cmd/server/main.go`.
3. At the start of `main`, read the connection string with `os.Getenv("DATABASE_URL")`.
4. If it is empty, stop startup with `log.Fatal` and a clear message. If it is set, continue starting the existing API normally.
5. Review the user's attempt and verify startup with and without the environment variable.

The next part will introduce the pgxpool dependency, a context with a timeout, pool creation, and an explicit `Ping` to verify connectivity. Explain multiple return values and resource cleanup as they arise. Migrations and the first workout table follow successful Go-to-database connectivity.

## Commands for resuming

Run these from the project root in separate terminals.

Go API:

```bash
go -C services/api run ./cmd/server
```

React website:

```bash
cd apps/web
vp run dev
```

Check the API directly:

```bash
curl -i http://localhost:8080/api/health
```

Check the web production build after finishing the exercise:

```bash
cd apps/web
vp run build
```

## Concepts discussed

- pnpm workspaces, root `node_modules/.pnpm`, dependency links, and keeping dependencies out of Git.
- Go module names and where `go mod init` creates `go.mod`.
- React components, JSX/TSX, and the roles of `main.tsx` and `App.tsx`.
- Go's `main` package, entry point, HTTP router, response writer, and response headers.
- `json.Encoder` writes JSON to a writer; `json.Marshal` returns JSON bytes. The health handler uses `Encoder` for convenience.
- Development proxying, React state, events, and `fetch` in the API button exercise.
- Ternary expressions, `instanceof Error`, and narrowing a caught value from `unknown` before reading `.message`.

## Later milestones

After completing the frontend-to-backend connection, introduce PostgreSQL and migrations, then build a small persistent workout flow. Clerk, TanStack Router/Query, dashboards, food logging, Expo/SQLite synchronization, and encrypted S3 backups are still future work described in `PLAN.md`.
