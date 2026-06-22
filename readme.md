# research-temporal

A minimal Go [Temporal](https://temporal.io) hello-world project with a `GreetSomeone` workflow.

## Prerequisites

- [Go](https://go.dev/dl/) 1.25+
- [Temporal CLI](https://docs.temporal.io/cli/setup-cli) (for running workflows)

## Quick start (recommended)

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

You should see output like: `Hello Test!`

## Alternative: Docker Compose

If you prefer running Temporal via Docker:

```sh
git clone https://github.com/temporalio/samples-server.git
cd samples-server/compose
docker compose up
```

Then follow the same worker and workflow steps above.
