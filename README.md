# Container Port Simulator

A seeded container terminal simulator with scheduling policies, resource constraints and a browser replay.

Only the scaffold exists; the implementation is planned. The scope, architecture and acceptance criteria are in [the project plan](docs/PROJECT-PLAN.md). The plan includes a reproducible verification strategy for the future build.

- [Product requirements](docs/PROJECT-PLAN.md#product-requirements)
- [Technical design](docs/PROJECT-PLAN.md#technical-design)
- [Implementation backlog and epic](https://github.com/jortega0033/container-port-simulator/issues)

Estimated implementation and verification effort is 19–31 engineering hours.

## Program goal

Build a full container port simulator and optimiser: vessels, quay, cranes, vehicles, yard, gate and rail in one deterministic engine, with planners that compete against proven bounds, calibrated uncertainty and an inspectable replay of every decision.

The program is planned ([#45](https://github.com/jortega0033/container-port-simulator/issues/45)). Phase 1 is the current epic [#1](https://github.com/jortega0033/container-port-simulator/issues/1) and the only phase with tickets.

**Score.** Average gap between the planners' results and the proven optimum or best bound on a fixed, pinned set of berth, crane and stacking benchmark instances.

| Phase | Content | Hours (rough) | Issue |
|---|---|---|---|
| 1 | The current epic and its tickets | 80 to 120 | [#1](https://github.com/jortega0033/container-port-simulator/issues/1) |
| 2 | Whole terminal | 130 to 220 | [#46](https://github.com/jortega0033/container-port-simulator/issues/46) |
| 3 | Optimisers | 180 to 300 | [#47](https://github.com/jortega0033/container-port-simulator/issues/47) |
| 4 | Traffic | 120 to 200 | [#48](https://github.com/jortega0033/container-port-simulator/issues/48) |
| 5 | Uncertainty | 100 to 170 | [#49](https://github.com/jortega0033/container-port-simulator/issues/49) |
| 6 | Workbench | 120 to 200 | [#50](https://github.com/jortega0033/container-port-simulator/issues/50) |

About 730 to 1210 engineer hours in total: rough estimates, not measurements. The numbers in the score and the finish lines are targets until a run shows them. Only phase 1 has tickets to work; the sub-issues of a later phase are written when it starts. No single session is expected to reach the end.

## License

[MIT](LICENSE) © 2026 Jake Ortega.
