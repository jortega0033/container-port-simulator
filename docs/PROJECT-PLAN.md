# Container Port Simulator

## Product requirements

A scheduling student or simulation developer can load an original synthetic manifest, run a container terminal under two policies, inspect a queue or container movement in a browser replay, and see how changed priorities affected service. The operator can pause or step through saved runs and inspect an individual ship, berth, crane, gate, or stack. Each movement links to the event and predecessor that caused it.

The first release models inbound containers only. A ship arrives, waits for a compatible berth, discharges through one of that berth's cranes, and leaves when its final container reaches the yard. A truck arrives to collect one container. A gate and its handling equipment can remove obstructing containers into other stacks before a pickup. An arriving truck can wait for discharge, a free gate, or a legal move. Reports show those reasons.

### First release

- Open the built-in hand-checkable scenario, start a run from a local Chrome page, and step through berth assignments, discharge, yard movements, and truck pickups.
- Import an original synthetic JSON manifest, reject invalid IDs and incomplete pickup references before a run starts, and run two policies on exactly the same input.
- Inspect per-container history, crane owner, queue reason, stack occupancy, and side-by-side metrics. Save a complete replay record and validate it on reload.
- Show that the reference run can be repeated byte-for-byte and that a different policy changes fixture reshuffles by the expected amount. A resource violation or truncated run is clearly marked as incomplete.
- Deliver a built-in two-day scenario with at least 20 calls, mixed loads and staged truck arrivals, plus a smaller exact verification fixture. Time and memory caps are explicit. Record observed runtime and outcomes when implemented.

The first release excludes export cargo, transshipment, vehicle geometry, tidal windows, labor shifts and real dispatch integrations. The terminal has synthetic rules and resource units. It can be used to compare outcomes within the model.

### Later milestones

1. Extend the workload generator to seasonality and peaks. Run an annual workload of 3,000 synthetic ship calls and verify exact replay of the recorded trace and checksum.
2. Add berth and crane scheduling beyond first-come-first-served. Keep zone exclusions and resource conservation in every scenario.
3. Freeze 20 annual workload seeds and compare tuned scheduling and yard policies against first-come-first-served. Work towards a 30% reduction for each chosen ship-wait and reshuffle metric. Publish the actual paired results and seed coverage. No assumed performance improvement is part of the first release.

## Technical design

Use one Go standard-library executable for a CLI, local HTTP API and bundled HTML, CSS and JavaScript. Serve the UI on loopback only. Go components are the manifest validator and generator, event engine, policies, trace/replay validation, and CLI and HTTP adapters. The CLI and browser run the same engine. A fresh workspace can build and test with the already installed Go toolchain, `GOPROXY=off` and `GOTOOLCHAIN=local`. Chrome uses bundled assets. Linux files and a browser tab are sufficient. No installed external database or simulation service is needed.

### Manifest and stable event ordering

Use a versioned JSON input with integer minutes, stable printable ASCII IDs, explicit ship arrivals, container ownership, booked pickup time, actual truck arrivals and terminal configuration. Cargo types are homogeneous in the first release. Berths specify compatible ship sizes and a crane pool. Every crane belongs to exactly one berth and operates one job at a time. A yard block specifies stacks in a fixed order, capacities and accepting berths. The terminal config defines positive service durations, gate count and one shared yard handler. Reject a container without exactly one declared owner and one truck. A run may include truck arrivals before its ship is discharged. Validate all cross-references before allocating resources or inserting events.

Generate workloads from a versioned in-repo PRNG specification with integer arithmetic and a fixed draw order. Freeze a manifest before policy execution. No policy uses the generator, clock, filesystem ordering, or actual future arrivals. Canonical form sorts key collections by ID, includes configuration and booked pickup times, and excludes wall-clock metadata. Record the manifest SHA-256 and the generator version/seed where used.

Use a priority queue keyed by minute, phase, entity ID and monotonic enqueue sequence. At a minute, apply completions in stable order, then admit pending arrivals in stable order, then dispatch jobs in a fixed resource order. A completion releases resources before new jobs can start. Dispatch gates first, then berth assignments, then crane jobs. A berth assignment can start compatible cranes that minute. At any tie use the creation sequence after phase and ID. Jobs take at least a minute, so new completions follow the dispatch minute. Process resource and output arrays in explicit order. Reject an event in the past. Check bounds on time, job count and trace size independently. A drained event queue with unfinished work is a blocked run with a reason, never a successful run. Stop on a violated state invariant and identify the first failing event.

### Constrained yard and dispatch

Assign compatible berths to waiting ships by arrival time and ID. A ship keeps its berth while discharging. A crane claims one cargo item and reserves a specific destination stack slot before launch. Each stack admits at most one pending operation, including a pending crane discharge. Releases, berth departures and arrivals at the same minute follow the phase order. Compute each container location as one of ship, in-flight crane, stack, or collected, with explicit locking during yard work. At every completed simulated operation verify container conservation, capacity, ownership, legal stacks and no overlap for a crane, handler, gate, or stack lock. A truck can never collect the same container twice.

A pickup preflights an ordered path to remove all blockers and then hand off its target. Reserve its gate and handler and lock all stacks it uses until the job completes. Each blocker moves from the top of the target stack to the first admissible non-target stack with room. A blocker move costs one configured move duration. The handoff has its own positive duration. Show each completed move in the trace rather than crediting the entire sequence at the start. Avoid mutating yard contents when preflight fails. A blocked truck keeps its original queue position and reason, while dispatch can choose the first feasible truck. Record skips. Fail with outstanding IDs when no event can make further progress. Ship discharge and pickup must agree on reservations and stack limits.

Implement two offline policies that share event timing, berth assignment and pickup dispatch:

- Baseline chooses first admissible stack by block and stack ID. Trucks are considered by arrival minute and ID. Ships are considered by arrival minute and ID within compatible berths.
- Pickup-aware placement scores each feasible stack by how many existing boxes it would bury ahead of their earlier booked pickup rank, then by stack height and stable stack ID. Rank is booked minute then truck ID. It receives booked times only. Realized arrivals are input events revealed to the engine as their minutes occur.

The first release assigns each crane its own fixed operating zone within its berth. This prevents concurrent crane ownership of an operating zone by construction. The event and lock checks verify that property at runtime. Moving cranes between berths and testing geometric collision belong to a later model.

### Replay and reporting

Log structured records with simulation minute, ordinal, event type, entity IDs, affected resource IDs and before/after location or queue fields. Do not serialize map iteration order. Set a versioned codec, canonical field order, integer durations and a SHA-256 over deterministic canonical trace bytes. Persist manifest hash, policy and trace hash. Wall-clock computation time can be logged separately without affecting equality. Trace size limits and invariant violations must identify incomplete runs. Save a viewable completed trace locally as a versioned file. Import validates limits, checksum and the supported version before changing the loaded replay. A scrubber advances only to verified snapshots and renders occupancy from an initial state plus event deltas.

In paired mode, require one immutable manifest hash, the same resource and timing rules, the same repetition seeds and an explicit policy ID. Compute per-ship berth wait from arrival to assignment, per-truck gate wait from actual arrival to job start, and truck turnaround from arrival to handoff. Count actual blocker movements per completed calendar run, excluding discharges and handoffs. Include completed and unfinished counts. Show totals, means, nearest-rank p95 and the exact denominator. Mark comparisons with incomplete work and show outstanding work separately. A percentage improvement is undefined for a zero baseline. Runtime is measured on the target machine and reported separately from simulated elapsed minutes.

### Independent verification fixtures

Fixture A starts with berth B1, crane C1, gate G1, one yard handler, and empty stacks S1/S2/S3 each of height two, ordered by ID. Discharge takes two minutes, blocker movement and handoff take one minute each. Ship A arrives at minute 0 with P then Q, and ship B arrives at minute 1 with R. All use B1. Truck T1 for P arrives and is booked at minute 7. Trucks T2 for Q and T3 for R arrive and are booked at minute 10. Only one simulation action for a resource can be active at once.

- Both policies assign ship A at 0, complete P at 2 and Q at 4, depart A at 4, assign ship B at 4 and complete R at 6. Ship waits are 0 and 3. Same-minute completion must free B1 before assignment.
- Baseline lands P under Q in S1 and R in S2. At 7 T1 preflights a move of Q onto R. That move completes at 8 and P leaves at 9. T2 collects Q at 11 and T3 collects R at 12. T2 starts at 10. T3 starts at 11 and waits one minute. Turnarounds are 2, 1 and 2 minutes. One blocker moves.
- Pickup-aware placement lands P in S1, Q in S2 and R in S3. At R placement, S2 would bury Q because T2 ranks ahead of T3 at the same booked minute. The next feasible empty stack wins. P leaves at 8, Q at 11 and R at 12. Gate waits stay 0, 0 and 1 minutes, turnarounds are 1, 1 and 2 minutes, and no blockers move. Both traces end with no owned crane, no occupied berth, and empty stacks.
- Fixture B uses one stack of height two, P underneath Q, no alternate stacks, and T1 scheduled to arrive ahead of T2. Before T2 appears, T1 is blocked and its original queue position is preserved. After T2 arrives it is feasible, and can collect Q first so T1 then collects P. A separate structurally valid variant defines a ship that is too large for every configured berth. Once that ship and its truck arrive and no job is active, report a blocked terminal state with the outstanding IDs and their exact wait reasons.
- Fixture C admits two compatible ships at one minute to different berths with distinct cranes and stacks. Assert distinct ownership at launch. Force both to request a single stack in a second variant, and verify that only one pending move can own that stack.

Write expected traces and invariants independently of the policy implementation. Fixture checks need event times and operations, not wall-clock speed. Compare a trace with a separately generated expected sequence and independently recalculate reports from events. Repeat runs from fresh processes and verify exact equality.

## Epic · Replay and compare an offline container terminal

Build and verify a standalone terminal simulator with an original two-day workload, a small exact trace, two constrained yard policies, and an inspectable offline browser workflow. The six tickets total **19–31 estimated engineering hours** including tests. Run CP-01 first. Develop CP-02 and CP-03 alongside each other against the shared engine and trace contracts. CP-04 integrates the policies and comparison. CP-05 adds the browser workflow using those contracts. CP-06 verifies the complete release.

### CP-01 · Define manifests and deterministic fixtures

**Estimate** 3–5 hours. **Dependencies** None.

Define versioned manifest, terminal configuration and event/trace shapes. Include the original hand-worked fixture plus a deterministic scenario generator and two-day input.

**Acceptance checks**

- The exact fixture has stable IDs and intended order, with manifest bytes and SHA-256 identical in two fresh runs from the same seed. A different seed changes the generated manifest reproducibly.
- Reject missing truck references, duplicate IDs, unknown resource IDs, invalid times and excess workload. Failed validation produces no trace or engine mutation.
- Document how generator parameters affect ships, trucks and load, and include an oracle written from the fixture's expected operations.

### CP-02 · Build the event queue and resource transitions

**Estimate** 4–6 hours. **Dependencies** CP-01.

Implement deterministic event order, berth assignment and crane discharge with fixed resource capacities and in-progress reservations. Expose the event contract to independent consumers.

**Acceptance checks**

- Fixture A shows ship assignment, waits, completion and departure at the exact times. Fixture C proves distinct ownership and observed lock enforcement.
- Completions and arrivals at the same minute produce the specified phase ordering. A resource completion releases ownership only once.
- Volume and event caps preserve an honest incomplete result. Impossible transfers report what remains, rather than looping indefinitely.

### CP-03 · Build yard placement and verified pickup transitions

**Estimate** 4–6 hours. **Dependencies** CP-01.

Implement yard capacity and location tracking, the baseline placement rule, and gate job preflight with one exclusive yard handler. Start against documented engine/trace contracts and synthetic events.

**Acceptance checks**

- A pending discharge reserves its slot, an uncommitted pickup leaves yards unchanged, and failed preflight does not leak a gate or lock.
- Fixture A's baseline trace shows the blocker move and three exact handoffs, with no duplicate cargo assignment or collection.
- Fixture B shows queue reason, a skipped blocked job and subsequent progress after another collection. The oversized-ship variant reports exact outstanding IDs and stable reasons for the lack of a compatible berth and an unavailable container.

### CP-04 · Compare pickup-aware placement against baseline

**Estimate** 3–5 hours. **Dependencies** CP-02 and CP-03.

Integrate both dispatch components, implement the booked-time scoring policy and record per-run and paired metrics.

**Acceptance checks**

- Both variants of Fixture A use the same manifest hash, terminal limits and event durations. Verify baseline result one blocker move and pickup-aware result zero, with all exact waits and turnaround figures.
- Include reports for a completed paired run, a blocked variant, a zero baseline and repeated seed runs. An independent reducer reproduces counts and quantiles.
- Actual arrivals later than bookings change execution when they arrive. The policy receives only booked ranks.

### CP-05 · Build local replay and diagnostics

**Estimate** 3–5 hours. **Dependencies** CP-04.

Build the loopback HTTP adapter and Chrome view. Display stacks, berth/crane ownership, queues and metrics. Include single-step, scrub, filters and trace reopening.

**Acceptance checks**

- In Chrome open Fixture A, step to minute 8, inspect one identified blocker move and see which truck triggered it. Pause, scrub backward and return to the same event with identical state.
- Compare both policies in the browser. Reopen the saved trace after restarting the server and verify the trace hash and inspected event.
- Broken input, unsupported trace versions and blocked runs receive specific feedback. Inspect the browser network log to confirm the UI loads from loopback only. Use keyboard-only controls for the whole fixture walkthrough.

### CP-06 · Verify repeatability and publish the release evidence

**Estimate** 2–4 hours. **Dependencies** CP-05.

Document offline launch and verification. Run the deterministic fixtures, two-day workload, invariant suite and browser acceptance walkthrough using only bundled assets and original inputs.

**Acceptance checks**

- From a clean checkout, record exact commands, environment, test results, paired workload statistics, wall-clock runtime and trace sizes.
- Repeat both policies in separate fresh processes and confirm exact equality of each policy's canonical trace bytes and checksum without elapsed runtime noise.
- Run at least five frozen two-day seeds, report all results including incomplete runs, and check state invariants on complete and blocked traces. Record limits and evidence gaps for an annual experiment.

## Scope decisions and risks

Yard job locks trade modeled throughput for deterministic, verifiable resource ownership. One handler limits shared simultaneous work. The initial comparison measures the declared synthetic system. Annual scale and geometric crane checks need their own workload and model validation.
