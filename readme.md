# research-temporal

A minimal Go [Temporal](https://temporal.io) project demonstrating workflows and **error-prone, non-deterministic activities**.

## Project structure

| File | Description |
|---|---|
| `hello/greeting.go` | `GreetSomeone` workflow — calls `SendNotification` then `ProcessPayment` with retries |
| `hello/activity.go` | Two activities: simulated notification service and payment gateway (random failures, variable latency) |
| `main.go` | Worker registration |

## Prerequisites

- [Go](https://go.dev/dl/) 1.25+
- [Temporal CLI](https://docs.temporal.io/cli/setup-cli) (for running workflows)

## Quick start

Start a local Temporal dev server:

```sh
temporal server start-dev --db-filename temporal.db
```

In another terminal, start the worker:

```sh
go run .
```

Execute the workflow:

```sh
temporal workflow execute \
  --type GreetSomeone \
  --task-queue greeting-tasks \
  --workflow-id my-first-workflow \
  --input '"Test"'
```

On success you'll see output like:

```
Hello Test! ntf-84321 | rcpt-USD-592013
```

If an activity fails, Temporal automatically retries (up to 5 times with exponential backoff) — the workflow will complete once all retries succeed or return the final error.

## Workflow vs Activity

Temporal draws a hard line between **workflow** and **activity** code:

| | Workflow (`greeting.go`) | Activity (`activity.go`) |
|---|---|---|
| **Deterministic** | Must be fully deterministic — same input always produces same result, no matter how many times replayed | **Non-deterministic by design** — can fail, sleep, call external APIs, use random numbers, read wall-clock time |
| **Replayed** | Yes — Temporal replays workflow code from the event history to recover state | No — only the result (or error) is recorded; the code itself is never replayed |
| **Restrictions** | No `time.Sleep`, `rand`, `time.Now`, network calls, or goroutines without SDK wrappers | No restrictions — use any Go library, any I/O, any blocking call |
| **Retries** | N/A (drives the orchestration) | Can be retried with backoff policies on failure |
| **Idempotency** | Not relevant | Should be idempotent if retried |

In this example, `GreetSomeone` (workflow) orchestrates the two calls and handles retries. `SendNotification` and `ProcessPayment` (activities) contain the actual non-deterministic logic — simulated failures, random latency, and random outputs.

## Activities (error-prone & non-deterministic)

- **`SendNotification`** — simulates an external push/email service. Random latency (500–2000ms) and ~40% transient failure rate.
- **`ProcessPayment`** — simulates a payment gateway. Validates amount bounds and rejects >5000; ~25% transient timeout rate. Returns a random receipt ID.

Both use `time.Sleep`, `rand`, and `time.Now()` — operations that are **non-deterministic** (and therefore belong in activities, not workflows).

## Alternative: Docker Compose

If you prefer running Temporal via Docker:

```sh
git clone https://github.com/temporalio/samples-server.git
cd samples-server/compose
docker compose up
```

Then follow the same worker and workflow steps above.
