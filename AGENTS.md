# AGENTS.md

> This file comes from a project template. Values in `<angle brackets>` are placeholders. If a placeholder you need is still unfilled, do not guess: ask me, or offer to fill it in from what is actually in the repo. Once confirmed, replace it in this file.

## 1. Think Before Coding
**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.
- Look for similar implementations in the codebase before writing new code.

## 2. Simplicity First
**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.
- Reuse existing code and the standard library - do not reinvent what already exists.
- Never add a new dependency (Go module, npm, or any other package) without asking me first.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes
**Touch only what you must. Clean up only your own mess.**

- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.
- Remove imports/variables/functions that YOUR changes made unused. Leave pre-existing dead code alone unless asked.

The test: every changed line should trace directly to the request.

**Comments - keep them minimal:**
- At most one short doc comment on a package or exported type (Go: `// Name ...`, TS: `/** ... */`) saying what it is for.
- No doc comments on every function/field, no comments that repeat what the code says, no commented-out code.
- An inline comment only when the *why* isn't obvious from the code.

**Keep things unexported by default:**
- Don't export identifiers, add interfaces, generics, struct tags or options "just in case".
- Export only what another package actually uses. Match what the surrounding code already does.

**Language:** everything written into the project is in English by default - code, identifiers, comments, UI text, error messages, docs, commit messages and PR texts - unless Section 5 says otherwise for this project.

## 4. Goal-Driven Execution
**Define success criteria. Loop until verified.**

- "Add validation" -> "Write tests for invalid inputs, then make them pass"
- "Fix the bug" -> "Write a test that reproduces it, then make it pass"
- "Refactor X" -> "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```
1. [Step] -> verify: [check]
2. [Step] -> verify: [check]
```
Run the relevant build/test/lint commands (Section 6) before saying a step is done.

## 5. Project
<!-- Fill in per project. Delete what doesn't apply. -->
- **Name:** <project-name>
- **Purpose:** <1-2 sentences>
- **Layout:**
  - `/backend` - <short description>
  - `/frontend` - <short description>
  - `/mobile` - <short description>
- **Specs:** if an `openspec/` folder exists, follow the active change proposal and do not implement anything outside it.

## 6. Stack & Commands
<!-- Keep only the layers this project has. -->

### Backend (Go)
- Go 1.25 (see `go.mod`). Standard library first: `log/slog` for logging, `net/http` with Go 1.22+ method/pattern routing (`mux.HandleFunc("GET /health", ...)`).
- `context.Context` is the first argument of every function that does I/O (database, network, messaging). Never store a context in a struct.
- Build: `go build ./...`
- Vet: `go vet ./...`
- Test: `go test ./...`
- Run: `go run ./cmd/server`
- Generate data access code: `sqlc generate` (after changing SQL queries or schema).
- Never apply migrations to a database - I apply them (Section 9).

### Frontend (<framework>)
- <framework, language, state management>
- Install: `<npm|pnpm> install`
- Dev: `<cmd>`
- Build: `<cmd>`
- Lint / typecheck: `<cmd>`
- Test: `<cmd>`
- New logic (stores, services, utils, composables/hooks) gets unit tests in the same task. A task is not done until they exist and pass.

### Mobile (<framework>)
- <framework, language>
- Run: `<cmd>`
- Test: `<cmd>`
- New logic (stores, services, utils, hooks) gets unit tests in the same task. A task is not done until they exist and pass.

### UI (Frontend, Mobile & server-rendered backend views) - HARD rules
These rules apply to every UI, including any server-rendered views in the backend.
- Use only my design system (based on Figma components): its components and tokens (colors, typography, spacing, radii, etc.). No other UI libraries, no hardcoded colors/sizes, no ad-hoc styles outside the tokens.
- Exception: the accent color is project-specific. It is overridden in exactly one place (`<accent-token-location>`), for both light and dark mode. Everything else uses the accent token - never hardcode the accent value in components.
- The design for a feature is made in Figma first and imported into `<design-folder>`. Before implementing any UI, find the matching design there and build exactly from it.
- If the design is missing, incomplete, or doesn't cover a state you need (empty, loading, error, mobile view, etc.) - stop and ask me. Never invent the design yourself.
- If the design needs a component the design system doesn't have yet - stop and ask. Don't create a one-off replacement.

## 7. Backend Architecture (layered)
The package boundaries below are strict. Never add an import that breaks them. If a task seems to need one, stop and ask.

### Layout and dependencies
The layers mirror the .NET template. There is no separate DTO package: persistence types live in `dataaccess` and web DTOs in `web/api`, each next to its mappers.

| Package | Contains | May import |
|---|---|---|
| `cmd/server` | `main` - loads config, builds dependencies, wires and starts everything (like `Program.cs`). No logic. | Anything |
| `internal/domain` | Core entities with their validation rules; sentinel errors; enums used across the app | Standard library only |
| `internal/contracts` | Interfaces only: services, repositories, external adapters. Provider-agnostic names and types | `domain` |
| `internal/application` | Business logic (services); implements the service interfaces in `contracts` | `domain`, `contracts` |
| `internal/dataaccess` | PostgreSQL repositories and sqlc-generated rows; maps rows <-> domain; implements the repository interfaces in `contracts` | `domain`, `contracts` |
| `internal/external/<concern>` | Concrete third-party adapters, one package per concern (`cache`, `rabbitmq`, ...); implements the external interfaces in `contracts` | `domain`, `contracts` |
| `internal/web/api` | HTTP handlers, request/response DTOs and their mappers (DTO <-> domain), router and server lifecycle | `domain`, `contracts` |
| `internal/config` | `.env` loader and the `Config` struct, loaded once at startup | Standard library only |
| `internal/platform/<concern>` | Cross-cutting infrastructure: logging, error reporting | Standard library and `config` |
| `tests/` | Integration tests (optional) | Anything |

### Naming
- Domain entity: no suffix - `AppUser`
- Web DTO: `...Dto` - `AppUserDto`
- Persistence rows: the types sqlc generates in `dataaccess`. Generated files are never edited by hand.

### Rules
- **Dependency direction:** `domain` depends on nothing; `contracts` on `domain`; `application`, `dataaccess` and `external` on `domain` + `contracts`; `web` on `domain` + `contracts`. Only `cmd/server` imports `dataaccess`/`external`/`application` to wire them. Never import across these boundaries otherwise.
- **Interfaces live in `contracts`**, never next to their implementation or in the consumer. Contract interfaces are required even when there is only one implementation - Section 2's "no abstractions for single-use code" does not apply to layer boundaries.
- **Domain entities are the only currency between layers.** Services and repositories take and return domain entities only.
- **Mapping happens only at the edges:** `dataaccess` maps rows <-> domain, `web/api` maps DTOs <-> domain. Map before calling a service, map after it returns. No ad-hoc mapping inside handlers, services or repositories - use the mapper functions.
- **Web returns only `...Dto` data.** Never expose domain entities or rows from `web`.
- **Contracts are provider-agnostic.** No third-party types in signatures, generic names only: `Cache.Set`, not `Redis.SaveToRedis`. Provider-specific names and types (Postgres, RabbitMQ, Redis, ...) exist only in `dataaccess`, `external` and `platform`.
- **Business validation rules live on domain entities**, not in DTOs, handlers or repositories.
- **Wiring lives only in `cmd/server`.** Keep `main` short: config, constructors, start.

### Errors
- Return `error`; never panic for expected failures. Panics only for programmer errors that cannot happen at runtime.
- Expected failures are sentinel errors (`var ErrNotFound = errors.New("...")`) declared in `domain`, so every layer can check them with `errors.Is`. Wrap with `fmt.Errorf("...: %w", err)` to add context.
- No custom error type hierarchies unless asked.
- `web/api` maps sentinel errors to status codes, `external/rabbitmq` to message dispositions; `dataaccess` maps driver errors to sentinel errors.
- Never swallow an error: handle it, return it, or report it (Section 8).

### Tests
- Standard `testing` package only, table-driven, in `_test.go` files next to the code. Integration tests may go in `tests/`.
- HTTP handlers are tested with `net/http/httptest`.
- Anything new that contains logic (domain validation, services, mappers, repositories, external adapters, config) gets tests in the same task.
- A task is not done until the new tests exist and `go build ./...`, `go vet ./...` and `go test ./...` pass.

### API
- Routes are URL-versioned and kebab-case plural: `/api/v1/<resources>`.
- Every API route uses the shared rate limiter.
- `PUT` and `DELETE` use optimistic concurrency: current concurrency token in the `If-Match` header.

## 8. Platform Integration
<!-- Only if the project uses the shared hubs. -->
### Contracts - HARD rules
- The JWT contract (`<jwt-contract-location>`) and the event contract (`<event-contract-location>`) are the source of truth. Read the relevant one before any work on auth/claims or events.
- Implement them exactly as written: same field names, types, casing, values, routing keys and structure. No renamed, added, removed, merged or "improved" fields, no local variants or workarounds in the app.
- If a contract doesn't fit the task, seems wrong, is ambiguous, or doesn't cover what's needed - stop and tell me before implementing anything that depends on it. Explain what doesn't fit and why. Never change, extend or reinterpret a contract yourself.

### Hubs
- Auth goes through Keycloak (realm: `<app-slug>`). Never implement custom auth.
- **API authentication uses JWT bearer tokens issued by Keycloak only.** No cookies, API keys, basic auth or any other scheme on API endpoints.
- **Server-rendered views** (if any) use Keycloak OpenID Connect login with a cookie session - the only exception. Never use the cookie scheme on API endpoints or the bearer scheme on views.
- Async communication goes through RabbitMQ using the event envelope from the event contract. Do not define a new envelope or change exchange/queue topology in code - topology is set up manually.
- **Errors:** see Section 7. Errors are reported to error-hub (GlitchTip), never only to the console. Don't `recover()` unless the code can actually handle the panic - never swallow errors.
- Email is sent via email-hub commands (`<app-slug>.<action>`), never directly from the app.
- Config comes from environment variables (`APP_SLUG` etc.) - see Section 10. Never hardcode secrets; never commit `.env`.

## 9. Version Control - HARD Rules
- Never run `git commit`, `git push`, `git merge`, `git pull`, `git rebase` or `git reset`.
- Creating a new branch (`git switch -c`) is allowed when starting a task on `main`/`develop`. Name: `<feature|fix|chore>/<2-4-keywords>`, lowercase, hyphen-separated. If it already exists, stop and ask.
- Never run SQL or migrations against a real database. Write migrations; I apply them.
- Never put `[DEPLOY]` in a suggested commit message - it triggers a production deploy. I add it myself.
- Commit messages and PR texts go in the chat only, never into files.

## 10. Environment Variables - HARD Rules
- Environment variables live only in `.env`. Never put config values anywhere else (no hardcoded defaults in code, no inline values in compose files or workflows).
- All variables are read in `internal/config` only. Other packages receive values from `Config`, never call `os.Getenv` themselves.
- A missing or invalid required value fails at startup with an error naming the variable.
- `.env.example` is a copy of `.env` with the same keys in the same order and all values empty.
- `.env` and `.env.example` must ALWAYS be in sync. Every added, renamed or removed variable is changed in both files in the same step.
- Format is exactly `UPPER_SNAKE_CASE=value` - no spaces around `=`, no other naming style. In `.env.example`: `UPPER_SNAKE_CASE=`
- Never invent secret values. A new secret goes into `.env` with an empty value and is listed under Manual steps for me to fill in.
- Every env change must also be applied everywhere the variable is referenced, ALWAYS including:
  - `.github/workflows/*` (e.g. `.env` generation from GitHub secrets)
  - Docker files: `Dockerfile*`, `docker-compose*.yml` / `compose*.yml`, and any other Docker-related folders or files
- Before finishing an env change, search the whole repo for the variable name to verify no reference was missed.

## 11. Change Summary (every task that changes files)
At the end of the response:
- **Branch** (if created)
- **Files changed** - path + 1 line on what changed
- **⚠️ Env variables changed** - if any variable was added, renamed or removed: list each variable, what happened to it, and every file where it was changed or is referenced (`.env`, `.env.example`, workflows, Docker files, code). If a GitHub secret must be added or renamed, say so explicitly.
- **Manual steps** - migrations, env vars, config changes I must do by hand
- **Suggested commit message** - in a code block, one per logical change, past tense, plain language
