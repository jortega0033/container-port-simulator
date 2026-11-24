# Owner decisions

A dated record of the decisions made for this repository, so the rules do not depend on anyone's memory of a conversation. Changing any entry needs the owner. An agent that finds a decision blocking a ticket stops and says so.

The owner made decisions 1 to 10 and 12 on 2026-11-22 for all repositories that get the SDLC kit, and asked on 2026-11-23 to apply the kit here (jortega0033/sdlc-kit v1.1.0). Each row says who decided.

| # | Decision | Source | Date | Implemented in |
|---|---|---|---|---|
| 1 | The instructions, tooling and playbook live on the default branch (`main`). | decided by the owner | 2026-11-22 | `CLAUDE.md`, `.claude/`, `.github/` |
| 2 | Persona agents run automatically during ticket work and carry instructions tailored to this project. | decided by the owner | 2026-11-22 | `.claude/agents/`, `.claude/agency-agents/README.md`, `CLAUDE.md` "Agent roles" |
| 3 | Every open ticket is about app implementation; tooling is installed, not ticketed. | decided by the owner | 2026-11-22 | `.claude/`, `.github/`, `docs/` |
| 4 | Agents may take a ticket up to a ready pull request. Only the owner merges and releases. | decided by the owner | 2026-11-22 | `CLAUDE.md`, `.claude/skills/next-ticket/SKILL.md`, `docs/RELEASING.md`, `.claude/settings.json` |
| 5 | No scheduler. Ticket work starts only on the owner's go signal (`/next-ticket`), because the repository is used to train AI models. | decided by the owner | 2026-11-22 | `.claude/skills/next-ticket/SKILL.md`, `LOOP.md` |
| 6 | A review gate hook blocks opening or readying a pull request until the gate roles have run after the last code edit. | decided by the owner | 2026-11-22 | `.claude/hooks/role_gate.py`, `.claude/settings.json` |
| 7 | Loop Engineering stays at L1, report-only and unscheduled. | decided by the owner | 2026-11-22 | `LOOP.md`, `loop-constraints.md`, `loop-budget.md` |
| 8 | Continuous integration runs on Linux and Windows; the default branch requires a pull request and green checks. Go tests run on Linux; the gate hook also runs on Windows. | decided by the owner | 2026-11-22 | `.github/workflows/ci.yml`; branch protection is a repository setting |
| 9 | Permission rules: documented commands are allowed; force pushes, pushes to the default branch, merges, releases and reads of secrets are denied; edits to CI, agent configuration and Loop files ask first. | decided by the owner | 2026-11-22 | `.claude/settings.json` |
| 10 | Open decisions on tickets are settled with the ticket's stated default; where a ticket states none, the conservative option is taken and flagged in the issue. | decided by the owner | 2026-11-22 | The list below; the issues themselves |
| 11 | Go is not installed on the owner's PC and nothing is installed globally. Go is verified in CI only, like Rust and Go in the owner's other repositories. The module has no external dependencies. | recommended by the orchestrator, following the owner's 2026-11-22 rule for Rust | 2026-11-23 | `CLAUDE.md` "Commands", `.github/workflows/ci.yml` |
| 12 | The ten UI rules in `docs/ui-quality.md` are binding, and the `ui-ux-pro-max` and `unslop` skills are installed in the repository and used for UI work and for copy. | decided by the owner | 2026-11-22 | `docs/ui-quality.md`, `.claude/skills/` |
| 13 | Visual system: the Pipenzo Flow look tokens. | decided by the owner | 2026-11-23 | `web/tokens.css`, `docs/design-tokens.md` |
| 14 | Token rules 14-D1 to 14-D4. | 14-D1 use of yellow decided by the owner; how it is drawn, and 14-D2 to 14-D4, recommended by the orchestrator and adopted under the owner's delegation of 2026-11-23 | 2026-11-23 | `docs/design-tokens.md`; applied in UI tickets |

## Open ticket decisions settled by default (entry 10)

The plan leaves these open. Each is a recommendation by the orchestrator, adopted under entry 10 on 2026-11-23, and holds unless the owner changes it in the issue. None blocks a ticket.

| Issue | Decision | Default |
|---|---|---|
| #2 | Go version | `go 1.26` in `go.mod`; CI installs that version. |
| #2 | Module path and layout | `github.com/jortega0033/container-port-simulator`, layout in `docs/architecture.md`. |
| #2 | Generator | SplitMix64 with a documented draw order and integer bounds, version `gen-v1`. The specification is in the repository. |
| #2 | Caps | Manifest 2 MiB, 200 ships, 5,000 containers, 10 berths, 20 cranes, 200 stacks, 20,160 simulated minutes, 200,000 events, 64 MiB trace. A cap gives an incomplete run or a validation error, never a silent cut. |
| #3 | Phase numbers | 0 completions, 1 arrivals, 2 gate dispatch, 3 berth dispatch, 4 crane dispatch. |
| #5 | Percentile | Nearest rank: the value at position ceil(0.95 n) of the sorted values, with n shown. |
| #6 | Address and traces | `127.0.0.1:9640`, tests use port 0. Saved traces go to `--data-dir`, default `./traces` (git-ignored), under server-made names. |
| #7 | Two-day seeds | Seeds 1 to 5 of generator `gen-v1`, frozen; all results are reported, incomplete ones too. |
| #7 | Verification command | `./portsim verify` prints one line per check and exits non-zero on the first failure; CI runs it on Linux. Timings are recorded and never asserted. |
| all | Owner checks | Manual browser steps run in Chrome by the owner. |
| none | Licence | MIT, added by the owner on 2026-11-23 (`LICENSE`). No ticket changes it. |

## Program decisions (owner, 2026-11-24)

The owner accepted the goal and each suggestion below. They apply from the phase they name; the rules of `CLAUDE.md` hold for phase 1 as written.

| Decision | Implemented in |
|---|---|
| The repository has a program goal above phase 1: a full container port simulator and optimiser, in six phases. | Program issue #45 |
| Published benchmark instances need a licence check. Suggested: fetch at a pinned version and run from outside the repository. | Program issue #45 |
| Suggested: no packages outside the Go standard library, the solvers included. | Program issue #45 |
