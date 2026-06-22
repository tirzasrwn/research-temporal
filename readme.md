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
