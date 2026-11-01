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

- **Planned features:** Workout logging, manual calories/macros, reusable meals, food search, barcode scanning, and progress dashboards. Deliver workout logging first, following the milestones below.
- **Workout logging:** Each workout contains exercise blocks with ordered sets. Record repetitions, weight in kilograms, and completion for each set, following the model below.
- **Web storage:** React calls the Go JSON API. PostgreSQL stores permanent fitness records; the browser holds temporary interface state.
- **Mobile storage:** Save immediately into SQLite. Existing local records remain usable when the home server or authentication services are unreachable.
- **Analytics:** Go queries PostgreSQL for web dashboards. Mobile calculates equivalent summaries from local SQLite. Share TypeScript models, validation, and pure calculations where applicable.
- **Custom synchronization:** Persist pending mobile operations, assign unique operation IDs, and acknowledge uploads so retries cannot duplicate records.
- **Change tracking:** All server writes feed a transactional, per-user change log. Download cursors and deletion records let mobile devices catch up after being offline.
- **Conflicts:** Check record revisions and preserve conflicting edits for explicit resolution.
- **Access control:** Go verifies Clerk tokens and enforces ownership for records, synchronization, and backups.
- **Import/export:** Provide versioned, per-user JSON archives.
- **Server recovery:** Schedule PostgreSQL backups, keep a copy away from the database host, and verify restoration into a separate database before relying on the app for daily records.
- **Cloud backups:** Encrypt archives on the client using AES-GCM and a separate recovery key before uploading through temporary S3 URLs. Keep recovery keys off the server.

## Workout logging model

- **Exercise library:** Provide a chooser containing shared, read-only preset exercises and the user's private custom exercises. Users can create an exercise and reuse it in future workouts. Both kinds have stable exercise IDs; the API permits selection only of presets or custom exercises belonging to the current user. Introduce exercise identity and selection with blocks in milestone 2.
- **Exercise blocks:** Add an owner-scoped `workout_exercises` entity between workouts and sets. Each block has its own UUID, selected exercise ID, position, revision, and creation time. Exercise names live in the library; sets reference their block through an owner-scoped foreign key. Selecting a library exercise creates a new block. The same exercise may appear in multiple independent blocks within a workout; matching names or exercise IDs never automatically merge them.
- **Exercise names:** Display the current library name in every workout. Renaming a custom exercise updates its label in current and past workouts while keeping the same exercise ID, sets, and progress history. Store no separate name snapshot on blocks or sets. Use rename for corrections or relabeling the same exercise; create a new exercise for a different movement or variation. Presets remain read-only.
- **Ordering:** Exercise blocks can be reordered, carrying their sets with them. New blocks and sets append to the end of their respective lists. Sets cannot be moved within or between blocks; their stored positions stay fixed. Editing values or completion never changes order, including in old workouts. Deleting a set leaves a gap in stored positions; display set numbers as 1..n by sorted index. Use IDs as deterministic tie-breakers for equal positions.
- **Adding sets:** Prefill kg/reps from the preceding set in the same block, copying each field independently. With no preceding set, both fields start blank; copied blanks remain blank. The new set always starts unchecked, even if its source is completed. Copies remain independently editable.
- **Completion:** Unchecked sets may have blank or partially filled fields, represented by nullable kg/reps. Supplied values must be valid: positive integer reps and finite, nonnegative kg; `0 kg` differs from blank. Checking requires valid values in both fields; enforce this in the API and database. Filling fields does not automatically check a set. Unchecking preserves kg/reps.
- **Time:** Client-recorded `completed_at` stores the checkmark tap time; null means unchecked. Preserve it across ordinary value edits and retries, clear it when unchecking, and record a new tap time when checking again. Keep server `created_at` for ingestion time. History and charts use the parent workout's date (`started_at`), including when checking a set in an old workout. This replaces the planned per-set `performed_at` requirement for this scope.
- **Existing data:** Preserve the current set sequence (`created_at, id`) by migrating each contiguous run of exactly matching exercise names into a block. Do not merge nonadjacent occurrences or differently spelled names, or automatically assign presets based on names alone. Give legacy exercises stable user-owned IDs while preserving existing names and independent blocks. Original block boundaries cannot be recovered reliably from flat data; adjacent same-name blocks cannot be distinguished. Mark existing sets completed with `completed_at = created_at`, an approximate legacy completion time.
- **Exercise history:** Defer history, charts, and the "Previous" column to milestone 5, using the stable exercise IDs established in milestone 2. "Previous" will show values from a previous workout; prefilling from the preceding set in the current block belongs in milestone 2.
- **Deletion conflicts:** Reject block deletion if the block or any of its sets changed since the client last reviewed them, including added or removed sets and value/completion edits. Show the latest contents and require a fresh delete action after review; preserve pending local edits and never automatically retry deletion with a newer revision. Check for changes and remove the block/sets in one transaction.

Reorder blocks in one transaction using the complete, duplicate-free block-ID list and the workout's expected revision. Reject stale revisions or a list that does not exactly match the workout's current blocks. Adding, deleting, or reordering blocks must invalidate stale order revisions. Record the complete resulting order and affected revisions in one change-log event so future offline clients can reproduce the change. Specify the detailed workout/block revision changes and event payload before implementation.

The server assigns set positions when appending; create requests carry no position. Decide whether deleted positions may be reused before implementation: `max(position) + 1` reuses the last position after its set is deleted, while strict non-reuse requires a persistent counter. Surviving sets keep their stored positions in either case.

## Delivery milestones

Work through each milestone as small learning exercises. Finish and verify the current exercise before starting the next architectural change.

1. **Finish editing a workout set.** Connect the existing edit form to the PUT helper, replace the saved row, and preserve drafts on failure. Add explicit stale-revision recovery that retains the draft while the user reloads the latest record. Verify persistence, failed saves, and stale edits before moving on.
2. **Establish exercise identity, blocks, and ordering.** Accept client-generated UUIDs for workouts, custom exercises, blocks, and sets; generate each ID once and retain it across create retries. Presets retain stable seeded IDs. Keep database-generated UUIDs as a fallback. A successful retry must not create another record or change-log event, and reusing an ID must never overwrite a record or disclose another user's data. Define how different create payloads using the same ID are rejected. Start with ID/retry handling, then introduce preset/custom exercise selection, custom exercise renaming, blocks, explicit block/set positions, and the migration above. Next add nullable fields, completion validation/checkmarks, and add-set prefilling; then implement block reordering. Preserve owner-scoped relationships and transactional change tracking throughout, including custom exercise changes.
3. **Complete workout logging.** Add set/block deletion and workout rename/delete. Deleting a block removes its sets in the same transaction and must reject changes the client has not reviewed, as described above. Record deletions transactionally so future offline clients can learn about every removed record, including blocks and sets removed with a workout. Define how deletion history is retained and how an old client recovers if its download cursor expires. Verify ownership and change-log behavior.
4. **Make workout logging ready for daily use.** Schedule database backups and demonstrate restoration into a separate database. Add Clerk authentication and verify account isolation. Choose and document remote access, TLS, and server configuration before exposing the app beyond local development. Authentication and tested recovery are prerequisites for deployment with real records.
5. **Add exercise history and workout dashboards.** Use the preset/custom exercise identities established in milestone 2 to add exercise history, the "Previous" column, and progress views. Keep repeated blocks independent when they link to the same exercise. Adopt TanStack Query as server-state handling grows and TanStack Router when adding screens. Extract Go packages when responsibilities become clear. Define analytics using shared inputs and expected outputs that both server and mobile implementations can test.
6. **Add nutrition in stages.** Start with manual calories/macros and reusable meals. Select a food data source after checking coverage and privacy requirements, then add search, barcode lookup, and caching for mobile offline use.
7. **Build Android storage, then synchronization.** Implement SQLite logging and a durable pending-operation queue before connecting sync. Add upload deduplication, download cursors, deletion handling, and revision checks. Preserve unsynced edits during conflicts; start with explicit reload/reapply choices before considering a richer resolution interface. Verify offline edits and interrupted retries against the existing server contracts.
8. **Add portable archives and optional cloud backups.** Implement versioned per-user export/import and restoration first. Specify the client responsible for encryption, recovery-key handling, and AES-GCM nonce management before adding encrypted S3 uploads. These archives complement the server recovery backups established earlier.

Weight representation remains unchanged during the edit exercise. Before dashboards and mobile calculations, decide the supported precision and representation across PostgreSQL, Go, TypeScript, and SQLite; changing only the PostgreSQL column does not settle calculation behavior.

## Validation, delivery, and defaults

- Verify preset/custom exercise selection and visibility, custom rename propagation to past workouts, stable exercise IDs, repeated exercise blocks, atomic block reorders, fixed set order after edits/deletions, completion validation, independent add-set copies, migration order/completion preservation, and rejection of stale block deletions.
- Test offline edits, interrupted synchronization, duplicate retries, conflicting edits, deletions, migrations, and account isolation.
- Test import/export round trips and backup restoration, including corrupted archives and incorrect recovery keys.
- Run `vp check`, `vp test`, `vp build`, Playwright browser tests, mobile component tests, and Go tests as appropriate.
- Benchmark with 50,000 workout sets and 50,000 food entries. Initial targets: mobile saves under 100 ms and standard charts under one second after data retrieval. Measure server queries and network latency separately.
- Deliver the website first, establishing IDs and change tracking from the beginning. Add Android and synchronization afterward; the Expo architecture supports iOS later.
- Keep AWS backups optional. The web app requires access to your home server.
- Treat Samsung health integration as a later feasibility project, retaining your requirement to avoid Samsung or Google health-data uploads.
