<!--VITE PLUS START-->

# Using Vite+, the Unified Toolchain for the Web

This project is using Vite+, a unified toolchain built on top of Vite, Rolldown, Vitest, tsdown, Oxlint, Oxfmt, and Vite Task. Vite+ wraps runtime management, package management, and frontend tooling in a single global CLI called `vp`. Vite+ is distinct from Vite, and it invokes Vite through `vp dev` and `vp build`. Run `vp help` to print a list of commands and `vp <command> --help` for information about a specific command.

Docs are local at `node_modules/vite-plus/docs` or online at https://viteplus.dev/guide/.

## Built-in Commands vs Scripts

`vp <name>` runs a built-in command. `vp run <name>` runs a `package.json` script or a `vite.config.ts` task. Scripts cannot overwrite built-ins, so `vp dev` and `vp run dev` may do different things. Check `package.json` and `vite.config.ts` first, and run `vp run <name>` when the project defines a script or task with that name.

## Running Scripts

Run TypeScript and JavaScript files with `vpx <file>`, for example `vpx ./scripts/seed.ts --dry-run`. It handles enums, decorators, JSX, and tsconfig `paths`, so do not add `tsx` or `ts-node` to run scripts. `vpx --watch <file>` reruns the file when a file it uses changes, in place of `tsx watch`. `vp node <file>` runs plain Node.js.

## Tool Versions

Run `vp toolchain` to show versions and relationships in the active Vite+
release. Add a tool name to select part of the graph. For example, run
`vp toolchain vite`. Use `--global` to ignore the local `vite-plus` package. Use
`vp why <package>` to show the package-manager dependency graph.

## Review Checklist

- [ ] Run `vp install` after pulling remote changes and before getting started.
- [ ] Run `vp check` and `vp test` to format, lint, type check and test changes.
- [ ] Check if there are `vite.config.ts` tasks or `package.json` scripts necessary for validation, run via `vp run <script>`.
- [ ] If setup, runtime, or package-manager behavior looks wrong, run `vp env doctor` and include its output when asking for help.

<!--VITE PLUS END-->

## Guided learning

- The user is learning through this project, especially Go. Prefer small, explained exercises and reviews of their attempts. Provide complete examples when asked, and honor explicit implementation requests.
- Keep explanations concise and give one exercise at a time. Use TypeScript comparisons when useful; expand when asked.
- When resuming a learning session, read `docs/learning.md` for the current exercise and checkpoint. Use `PLAN.md` for the architecture and roadmap.
- Keep `docs/learning.md` as a rolling checkpoint: target 300 words, maximum 400. Replace outdated entries in place; do not append session history or create a separate history archive.
- Record only current progress, the active exercise or next small step, essential constraints/blockers, useful startup commands, and the latest relevant verification. Briefly distinguish code review, checks actually run, and user-reported results once.
- Keep architecture in `PLAN.md` and implementation details in code. Omit copied code, exhaustive test cases, resolved failures, commit/dirty-file inventories, repeated caveats, and conversational recaps.
- Update the checkpoint when progress or the next step changes. For documentation-only updates, check the diff and word limit; avoid unrelated installs, builds, or tests.
