# Gym app stack with Vite+

## Summary

Build a phone-friendly website first, with permanent records stored in PostgreSQL on your computer or home server. Later, add an offline-capable Android app using SQLite and custom synchronization. AWS provides optional encrypted backups.

## Stack and development tooling

| Area | Choice |
|---|---|
| Website | React + TypeScript, CSS Modules |
| Web toolchain | [Vite+](https://viteplus.dev/guide/): Vite, Rolldown, Vitest, Oxlint, and Oxfmt |
| Web routing | TanStack Router |
| Web data fetching | TanStack Query with an in-memory cache |
| Backend | Go, standard-library `net/http`, and pgx |
| Server database | PostgreSQL on your own hardware |
| Authentication | Clerk with Google sign-in |
| Mobile app | React Native + Expo, Expo Router, and Metro |
| Mobile database | `expo-sqlite` |
| Charts | ECharts on web; Victory Native on mobile |
| Cloud backups | Private AWS S3 storage |

Use one repository with a pnpm workspace for web, mobile, and shared TypeScript packages. Keep the backend in its own Go module.

Use `vp dev`, `vp check`, `vp test`, and `vp build` for the website and shared TypeScript code. Enable type-aware linting and TypeScript checks. Use `vp run` to orchestrate Go, Expo, and other project tasks.

Expo retains its [Metro build pipeline](https://docs.expo.dev/guides/customizing-metro/). Mobile component tests use Jest, `jest-expo`, and React Native Testing Library; Go retains its native test tools.

## Application behavior and data flow

- **Initial features:** Workout logging, manual calories/macros, reusable meals, food search, barcode scanning, and progress dashboards.
- **Workout sets (first version):** Record an exercise name, repetitions, and weight in kilograms for each set belonging to a workout.
- **Web storage:** React calls the Go JSON API. PostgreSQL stores permanent fitness records; the browser holds temporary interface state.
- **Mobile storage:** Save immediately into SQLite. Existing local records remain usable when the home server or authentication services are unreachable.
- **Analytics:** Go queries PostgreSQL for web dashboards. Mobile calculates equivalent summaries from local SQLite. Share TypeScript models, validation, and pure calculations where applicable.
- **Custom synchronization:** Persist pending mobile operations, assign unique operation IDs, and acknowledge uploads so retries cannot duplicate records.
- **Change tracking:** All server writes feed a transactional, per-user change log. Download cursors and deletion records let mobile devices catch up after being offline.
- **Conflicts:** Check record revisions and preserve conflicting edits for explicit resolution.
- **Access control:** Go verifies Clerk tokens and enforces ownership for records, synchronization, and backups.
- **Import/export:** Provide versioned, per-user JSON archives.
- **Cloud backups:** Encrypt archives on the client using AES-GCM and a separate recovery key before uploading through temporary S3 URLs. Keep recovery keys off the server.

## Validation, delivery, and defaults

- Test offline edits, interrupted synchronization, duplicate retries, conflicting edits, deletions, migrations, and account isolation.
- Test import/export round trips and backup restoration, including corrupted archives and incorrect recovery keys.
- Run `vp check`, `vp test`, `vp build`, Playwright browser tests, mobile component tests, and Go tests as appropriate.
- Benchmark with 50,000 workout sets and 50,000 food entries. Initial targets: mobile saves under 100 ms and standard charts under one second after data retrieval. Measure server queries and network latency separately.
- Deliver the website first, establishing IDs and change tracking from the beginning. Add Android and synchronization afterward; the Expo architecture supports iOS later.
- Keep AWS backups optional. The web app requires access to your home server.
- Treat Samsung health integration as a later feasibility project, retaining your requirement to avoid Samsung or Google health-data uploads.
