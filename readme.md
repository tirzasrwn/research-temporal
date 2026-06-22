# reseach-temporal

## Temporal cluster from docker compose

- start local temporal server

```sh
git clone https://github.com/temporalio/samples-server.git
cd samples-server/compose
docker compose up
```

- install temporal cli

[instalation link](https://docs.temporal.io/cli/setup-cli).

- run the worker

```sh
go run .
```

- run the workflow

```sh
temporal workflow execute --type GreetSomeone --task-queue greeting-tasks --workflow-id my-first-workflow --input '"Test"'
```

## Temporal cluster from temporal cli

- install temporal cli

[instalation link](https://docs.temporal.io/cli/setup-cli).

- start the temporal cluster

```sh
temporal server start-dev --db-filename temporal.db
```

- run the worker

```sh
go run .
```

- run the workflow

```sh
temporal workflow execute --type GreetSomeone --task-queue greeting-tasks --workflow-id my-first-workflow --input '"Test"'
```
