# research-temporal

Minimal [Temporal](https://temporal.io) project demonstrating workflows with error-prone, non-deterministic activities, and a **cross-language parent/child chain**: Go's `GreetSomeone` workflow runs first, then starts TypeScript's `fulfillGreeting` workflow as a child and pipes its own result into it as the child's input. The child's result is what the Go workflow ultimately returns.

| Language | Path | Workflow | Activities | Task queue |
|---|---|---|---|---|
| Go | `go/` | `GreetSomeone` (parent) | `SendNotification`, `ProcessPayment` | `greeting-tasks` |
| TypeScript | `ts/` | `fulfillGreeting` (child) | `sendGreeting` | `fulfillment-tasks` |

## How the chain works

The data flows inside the workflow itself, not in a separate client:

1. `GreetSomeone` runs `SendNotification` and `ProcessPayment`.
2. It builds a receipt string `Hello <name>! ntf-<id> | rcpt-<id>`.
3. It calls `workflow.ExecuteChildWorkflow` with `TaskQueue: "fulfillment-tasks"` and the receipt as the only argument — Temporal routes the child to the TS worker, regardless of the language mismatch.
4. The TS worker runs `fulfillGreeting`, which validates the format and calls `sendGreeting`, returning `Hello <name>! ntf-<id> | rcpt-<id> -> dlv-<id>`.
5. `GreetSomeone` returns the child's result.

A single CLI invocation starts the whole chain. The Temporal UI shows the child workflow nested under the parent run.

## Project structure

```
.
├── go/                        # Go worker
│   ├── main.go                # Worker registration (greeting-tasks)
│   ├── go.mod
│   ├── go.sum
│   └── hello/
│       ├── greeting.go        # GreetSomeone workflow (parent)
│       └── activity.go        # SendNotification, ProcessPayment
├── ts/                        # TypeScript worker
│   ├── package.json
│   ├── tsconfig.json
│   └── src/
│       ├── worker.ts          # Worker registration (fulfillment-tasks)
│       ├── workflows.ts       # fulfillGreeting workflow (child)
│       └── activities.ts      # sendGreeting
└── readme.md
```

## Prerequisites

- [Go](https://go.dev/dl/) 1.25+
- [Node.js](https://nodejs.org/) 20+
- [Temporal CLI](https://docs.temporal.io/cli/setup-cli)

## Quick start

Start the dev server and both workers in three terminals:

```sh
# terminal 1
temporal server start-dev --db-filename temporal.db

# terminal 2
cd go && go run .

# terminal 3
cd ts && npm install && npm run dev
```

Run the chain with a single CLI command — the Go worker handles the first two activities, then starts the TS child, then awaits and returns its result:

```sh
temporal workflow execute \
  --type GreetSomeone \
  --task-queue greeting-tasks \
  --workflow-id my-workflow \
  --input '"Test"'
```

Expected output:

```
Hello Test! ntf-84321 | rcpt-USD-592013 -> dlv-7f2c1d9a
```

## What each side does

### Go: `GreetSomeone` on `greeting-tasks` (parent)

Three-stage orchestrator:

- **`SendNotification`** — simulated push/email service, 500–2000ms latency, ~40% transient failure rate.
- **`ProcessPayment`** — simulated payment gateway, 2s latency, rejects `amount <= 0` and `amount > 5000`, ~25% transient timeout rate.
- **`ExecuteChildWorkflow("fulfillGreeting", ...)`** — starts the TS child on `fulfillment-tasks` with the receipt string as input. The parent blocks until the child completes, then returns the child's result.

10s `startToCloseTimeout`, 5 retries with exponential backoff on each activity.

### TypeScript: `fulfillGreeting` on `fulfillment-tasks` (child)

Single-activity workflow that consumes the receipt from the parent.

- **Format check (in workflow, deterministic):** the receipt must start with `Hello ` and contain ` | `. Otherwise the workflow throws `ApplicationFailure.nonRetryable` of type `InvalidGreetingFormat` and the parent fails fast.
- **`sendGreeting`** — simulated delivery service, 500–1500ms latency, ~20% transient failure rate. Returns a delivery ID.

10s `startToCloseTimeout`, same 5-attempt exponential backoff.

The format check lives in the workflow (deterministic, non-retryable) and the simulated failure lives in the activity (non-deterministic, retryable) — same split as the Go side.

## Workflow vs Activity

Temporal draws a hard line between **workflow** and **activity** code:

| | Workflow | Activity |
|---|---|---|
| **Deterministic** | Must be fully deterministic — same input always produces same result, no matter how many times replayed | **Non-deterministic by design** — can fail, sleep, call external APIs, use random numbers, read wall-clock time |
| **Replayed** | Yes — Temporal replays workflow code from event history to recover state | No — only the result (or error) is recorded; the code itself is never replayed |
| **Restrictions** | No `time.Sleep`, `rand`, `time.Now`, network calls, or goroutines without SDK wrappers | No restrictions — use any language library, any I/O, any blocking call |
| **Retries** | N/A (drives the orchestration) | Can be retried with backoff policies on failure |
| **Idempotency** | Not relevant | Should be idempotent if retried |

In both study cases the workflow file is small and orchestration-only; all the interesting non-deterministic work lives in the activities. Workflows can call other workflows (the parent/child relationship in this project) — the same determinism rules apply, and the child is itself a workflow, not an activity.

## Alternative: Docker Compose

If you prefer running Temporal via Docker:

```sh
git clone https://github.com/temporalio/samples-server.git
cd samples-server/compose
docker compose up
```

Then follow the same worker and workflow steps above.
