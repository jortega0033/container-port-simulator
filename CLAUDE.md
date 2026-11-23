# CLAUDE.md

Working instructions for this repository. Read this file first, then the documents it links to. It does not repeat them.

## What this is

A local simulator for a scheduling student: load a synthetic manifest, run an inbound container terminal under two policies, replay berth, crane, yard and truck movements in Chrome, and compare waits and reshuffles on the same input. Go standard library only, `127.0.0.1` only, deterministic traces with checksums.

**Only the scaffold exists, with no product behaviour.** The repository holds the plan, the delivery tooling, the Go scaffold (module, packages, smoke tests, CI) and the backlog in epic #1. Commands marked planned below do not work yet; `portsim generate`, `run`, `compare`, `serve` and `verify` exit non-zero with "not implemented yet". Do not describe a planned thing as working, and do not claim a passing baseline: there is none.

## Sources of truth

| Topic | Document |
|---|---|
| Product scope, design, fixtures with expected values | [docs/PROJECT-PLAN.md](docs/PROJECT-PLAN.md) |
| Planned module map, freeze points, planned versus implemented | [docs/architecture.md](docs/architecture.md) |
| Backlog | Epic issue #1 and its child tickets on GitHub. The epic's table is authoritative for order, dependencies and estimates; tickets state their own owner-decision defaults. |
| Loop setup and its limits | [LOOP.md](LOOP.md), [loop-constraints.md](loop-constraints.md) |
| Project subagents, origin and local changes | [.claude/agency-agents/README.md](.claude/agency-agents/README.md) |
| Ticket workflow and the step at which each role runs | [.claude/skills/implement-ticket/SKILL.md](.claude/skills/implement-ticket/SKILL.md), section "Agent roles" below |
| How a ticket goes from the queue to a ready pull request, who may do what | "Delivery workflow" below |
| UI and copy rules, design tokens | [docs/ui-quality.md](docs/ui-quality.md), [docs/design-tokens.md](docs/design-tokens.md) |
| Decisions the owner made, with dates | [docs/owner-decisions.md](docs/owner-decisions.md) |

The plan wins for product design; the epic and tickets win for order and scope. A ticket that changes a contract fixes the plan in the same pull request.

## Stack and constraints

- Go module, **no external dependencies**: standard library only, in the application and in the tests.
- One executable: command line and loopback server run the same engine; browser assets are embedded.
- Deterministic: same manifest, policy and version give the same canonical trace bytes in a fresh process.
- Offline: same-origin assets only, original synthetic inputs only.
- Inbound containers only, one handler, fixed crane zones, two policies.
- Export cargo, transshipment, vehicle geometry, tides, shifts, moving cranes, annual workloads and performance claims are later work.
- A runtime version change, a dependency, an external request or a non-loopback listener needs the owner's approval first, not a pull request.

## Commands

Go is not installed on the owner's PC and nothing is installed globally. **Go checks run in CI only.** Run `go version` first: if it is missing, do not install it, and do not claim you ran any command below.

| Purpose | Command | Status |
|---|---|---|
| Build | `go build -o portsim ./cmd/portsim` | real (scaffold), verified in CI only |
| Format check (empty output is good) | `gofmt -l .` | real (scaffold), verified in CI only |
| Static check | `go vet ./...` | real (scaffold), verified in CI only |
| Tests | `go test ./...` | real (scaffold), verified in CI only |
| Generate a manifest | `./portsim generate --seed 1` | planned: #2 |
| Run a policy | `./portsim run --policy baseline <manifest.json>` | planned: #5 |
| Compare two policies | `./portsim compare <manifest.json>` | planned: #5 |
| Run the app | `./portsim serve` | planned: #6 |
| Release checks | `./portsim verify` | planned: #7 |

## Architecture boundaries

The planned module map is in [docs/architecture.md](docs/architecture.md). Invariants every change must keep (the plan is the source):

- Manifest, engine, yard, policies, trace, report and adapters stay separate. Command line and browser run the same engine.
- Events are ordered by minute, phase, entity ID and enqueue sequence. In a minute: completions, then arrivals, then dispatch of gates, berths and cranes. No map iteration order, clock or filesystem order reaches a result.
- Time is integer minutes and every job takes at least one. An event in the past, a broken invariant or a cap stops the run at the first failing event and marks it incomplete.
- A manifest is validated completely (stable ASCII IDs, references, times, workload caps) before any resource or event exists. A failed validation leaves no trace and no engine change.
- A drained queue with unfinished work is a blocked run with outstanding IDs and reasons, never a successful one.
- After every completed operation: container conservation, capacity, ownership, legal stacks, and no overlap for a crane, handler, gate or stack lock. Each stack admits one pending operation.
- A pickup plans its whole blocker path first, then takes gate, handler and stack locks. A failed preflight changes nothing. A blocked truck keeps its queue position and reason, and skips are recorded.
- Policies share timing, berth assignment and pickup dispatch. Pickup-aware placement sees booked ranks only, never actual arrivals, and no policy sees the generator or a clock.
- Traces use a versioned codec, canonical field order and SHA-256; computation time stays outside. A reopened trace is checked (limits, version, checksum) before the loaded replay changes.
- Paired runs need one manifest hash. Report totals, means, nearest-rank p95 and the denominator; incomplete work separately; a change from a zero baseline is undefined.
- Expected traces and the report reducer used in tests share no code with the engine, the policies or the reporter.
- The server binds loopback, checks host and origin on writes, caps request bodies, and builds no path from an HTTP value.
- IDs and reasons are text nodes. The scrubber moves only to verified snapshots. Same-origin assets only.
- Step, pause, scrub and filters work by keyboard. Occupancy, ownership, queue reasons and run status are text, not colour alone.

## Git workflow

- `main` is the default branch. Work on one branch per ticket: `<issue-number>-<short-slug>`. Never commit to `main`.
- Small commits with imperative subjects and a type prefix (`feat:`, `fix:`, `docs:`, `test:`, `chore:`).
- Order for every ticket: local evidence, Code Reviewer, Reality Checker, then a draft pull request that says `Part of epic #1` and `Closes #<ticket>`, then CI on that draft, then `gh pr ready` only when CI is green.
- Do not push or open a pull request unless the owner asked in the current session; starting `/next-ticket` is that request for one ticket. Only the owner merges and releases. Do not push to `main`, force-push, or merge your own pull request.
- Never commit saved traces (`traces/`), `.env*`, or build output.
- Do not change `.claude/agents/`, `LOOP.md`, `loop-constraints.md` or `loop-budget.md` as part of a feature ticket.
- Open owner decisions on a ticket are settled with the ticket's stated default. Where a ticket states none, take the conservative option and say so in the pull request.

## Evidence before a ticket is called done

Paste real output into the pull request; never summarise a run you did not make.

1. For each acceptance criterion of the issue: the name of the test that proves it, or the exact step and what was observed.
2. Go results come from CI, and green CI on the current diff is a condition for `gh pr ready`: the link to the CI run on the draft pull request with the Go version and the results of `gofmt -l .`, `go vet ./...` and `go test ./...`. If Go is not available locally, say so.
3. For an engine, yard or policy change: the fixture with expected and observed minutes, and the trace hash from two fresh processes.
4. For a report change: the independent reducer's numbers next to the report's.
5. For a UI change: the browser and version, the steps and what was observed, and the manual steps left as owner checks.
6. `git status` shows no stray files, and no dependency, external URL or non-loopback listener was added.
7. Anything not verified, said plainly.

## Agent roles

Project subagents are in `.claude/agents/`: eight Agency Agents roles and `loop-verifier`, which belongs to the Loop setup. Each Agency Agents role file starts with a section "In this repository", written for this project. That section and this file override the generic text below it. Sources, licence, local changes and the name collision with a user-level install are in [.claude/agency-agents/README.md](.claude/agency-agents/README.md).

When asked to implement, continue or finish a GitHub issue, follow the playbook [`.claude/skills/implement-ticket/SKILL.md`](.claude/skills/implement-ticket/SKILL.md). Claude Code loads it on its own for ticket work; `/implement-ticket <issue number>` starts it by hand. Delegate to the roles without waiting to be asked: each role's `description` ends with a "Use proactively ..." sentence, which is what makes Claude Code delegate. The playbook gives the exact trigger of each role.

| When | Roles | Access |
|---|---|---|
| Before coding | Software Architect | `Read, Grep, Glob` |
| While coding | Backend Architect, Frontend Developer | can edit |
| After coding, when the trigger matches | Application Security Engineer, Accessibility Auditor | `Read, Grep, Glob` |
| After coding, when the trigger matches | API Tester | can edit |
| Gate on every ticket | Code Reviewer on the full diff, until no material finding is left | `Read, Grep, Glob` |
| Acceptance gate on every ticket, after Code Reviewer | Reality Checker on each acceptance criterion and the recorded evidence, and that CI is configured to run the checks (no CI result is needed yet) | `Read, Grep, Glob` |

- The two gates run on every ticket, whatever its size. A ticket is not called done before both have passed.
- Roles limited to `Read, Grep, Glob` cannot run commands. Give them the diff, the acceptance criteria and whatever command output exists.
- The session or role that wrote a change does not count as its reviewer.
- If a role is missing or cannot be delegated to, the session does the before-coding, while-coding and specialist steps itself and says so. It cannot pass the two gates on its own change: do not push or open the pull request. Stop, report that the gates could not run, and leave the run for the owner. This overrides the playbook's sentence "do that step yourself" for the two gates.
- In the pull request description, list under "Agent roles run" each role that ran, what it was asked and its verdict, and any step done without a role.

## Delivery workflow

How a ticket reaches the owner. Each piece is described in its own file; this part only links them.

- **Start:** the owner runs `/next-ticket` ([`.claude/skills/next-ticket/SKILL.md`](.claude/skills/next-ticket/SKILL.md)); it picks one leaf ticket (never a `tracking` parent or a ticket with open sub-issues; the last sub-issue's pull request also closes the parent), creates the branch and worktree, and runs the playbook [`.claude/skills/implement-ticket/SKILL.md`](.claude/skills/implement-ticket/SKILL.md). Nothing is scheduled.
- **Gate hook:** [`.claude/hooks/role_gate.py`](.claude/hooks/role_gate.py) blocks `gh pr create` and `gh pr ready` until Code Reviewer and Reality Checker ran after the last code edit. `ROLE_GATE=off` in the environment Claude Code is started from switches it off for that session; only the owner does that.
- **Permissions:** [`.claude/settings.json`](.claude/settings.json). Agent configuration edits ask in the main checkout and in ticket worktrees; CI, `CLAUDE.md`, `AGENTS.md` and Loop-file edits ask only in the main checkout; in a ticket worktree they are limited by the playbook's scope rule and reviewed by the owner at merge (shell-made edits are not covered, so a ticket never changes agent configuration); force pushes, pushes to `main`, `gh pr merge`, releases and reads of `.env` are denied.
- **CI:** [`.github/workflows/ci.yml`](.github/workflows/ci.yml). Every job must be green on the current diff before a pull request is marked ready. The agent does not poll: it ends its turn with "Draft pull request #N is waiting for CI" and reads `gh pr checks` once when woken. Reviewers read the diff and evidence from `.claude/state/review/`. An acceptance step that needs a person is listed under "Owner checks" in the pull request.
- **Pull request:** [`.github/pull_request_template.md`](.github/pull_request_template.md). Tickets use [`.github/ISSUE_TEMPLATE/implementation-ticket.md`](.github/ISSUE_TEMPLATE/implementation-ticket.md).
- **Labels:** `status:ready`, `status:blocked`, `status:in-progress`, `needs-decision`, `tracking` (parent container, never worked; its sub-issues are the work items).
- **Changelog:** every ticket pull request adds one line under Unreleased in [CHANGELOG.md](CHANGELOG.md).
- **UI and copy:** [docs/ui-quality.md](docs/ui-quality.md); the `ui-ux-pro-max` and `unslop` skills are in `.claude/skills/`. Colours come from the shared design tokens in `web/tokens.css` ([docs/design-tokens.md](docs/design-tokens.md)); new UI adds no colour literals.
- **Decisions:** [docs/owner-decisions.md](docs/owner-decisions.md). Security reports: [SECURITY.md](SECURITY.md).
- **Merge and release:** only the owner merges and releases ([docs/RELEASING.md](docs/RELEASING.md)). An agent stops at a ready pull request.

## Loop

Issue triage runs at L1: report-only, started by hand, no scheduler, no GitHub writes. See [LOOP.md](LOOP.md).
