# Implementer map

Status: **only the scaffold exists, no product behaviour.** Every package below is empty with a smoke test, except `cmd/portsim` (help, version and stubs that exit non-zero); `web/` holds only the design tokens and embeds nothing yet. This page maps the boundaries that [PROJECT-PLAN.md](PROJECT-PLAN.md) requires onto a module layout so that parallel tickets do not collide. It adds no requirement; the plan wins where they differ. The names are the orchestrator's proposal, fixed by the scaffold; a ticket that changes one corrects this page in the same pull request.

## Planned module map

| Boundary | Module | Owns | Parent ticket |
|---|---|---|---|
| Entry point | `cmd/portsim` | Subcommands `generate`, `run`, `compare`, `serve`, `verify` | #2 (`generate`), #5 (`run`, `compare`), #6 (`serve`), #7 (`verify`) |
| Manifest and generator | `internal/manifest` | Schema, validation, canonical form, hash, seeded generator, fixtures | #2 |
| Trace | `internal/trace` | Event and trace shapes, canonical bytes, hash, reopening and snapshots | #2 (shapes), #6 (reopening) |
| Engine | `internal/engine` | Priority queue, phases, berths, cranes, caps, invariants | #3 |
| Yard | `internal/yard` | Stacks, locations, locks, pickup preflight, handler | #4 |
| Policies | `internal/policy` | Baseline placement, pickup-aware placement | #4 (baseline), #5 (pickup-aware) |
| Reports | `internal/report` | Per-run and paired metrics | #5 |
| HTTP | `internal/httpapi` | Routes, guards, static assets, trace files | #6 |
| Browser assets | `web/ (embedded)` | Replay, inspection, comparison | #6 |
| Verification | `internal/verify` | Fixtures, two-day seeds, repeatability, evidence | #7 |

Dependency direction: `httpapi` and `cmd/portsim` depend on `engine`, `report`, `manifest` and `trace`, `httpapi` also on `web`, and `cmd/portsim` also on `httpapi` and `verify`; `verify` uses the packages it checks; `engine` depends on `yard`, `policy`, `manifest` and `trace`; `yard` and `policy` on `manifest`, `yard` also on `trace`. `manifest` imports nothing from the module.

## Freeze points

Parallel tickets rely on these shapes. A ticket that changes a frozen shape updates this page and every dependent ticket in the same pull request.

- **Manifest, terminal configuration and event/trace shapes** (schema version 1): ticket #2 writes them before #3 and #4 start in parallel.
- **Yard and dispatch interface** used by the engine: agreed between #3 and #4, joined in #5.
- **Report JSON** (#5) and **replay snapshot JSON** (#6).

## What is verified where

Go is verified in CI only (see [CLAUDE.md](../CLAUDE.md)), on Linux and Windows. Repeatability is checked between fresh processes. Manual browser steps are owner checks in Chrome.
