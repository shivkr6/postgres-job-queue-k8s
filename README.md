# Postgres Job Queue

A small Go project for learning how to build a PostgreSQL-backed job queue from first principles.

The project currently implements:

- Milestone 1: Go CLI skeleton and PostgreSQL connection check
- Milestone 2: database migration for the queue schema

Future milestones will add enqueueing, stats, job claiming, workers, retries, stuck-job recovery, and Kubernetes deployment.

For the local Kubernetes learning setup, see [`k8s/README.md`](k8s/README.md).

## Requirements

- Go 1.25+
- Docker
- Docker Compose
- kind and kubectl for the Kubernetes milestones

## Start Postgres

```bash
docker compose up -d
```

This starts a local PostgreSQL container using:

```text
user: queue
password: queue
database: queue
host port: 5432
```

The default app database URL is:

```text
postgres://queue:queue@localhost:5432/queue?sslmode=disable
```

Inside Kubernetes, the database URL will use the Postgres Service name instead of `localhost`:

```text
postgres://queue:queue@postgres:5432/queue?sslmode=disable
```

You can override it with:

```bash
DATABASE_URL="postgres://user:password@host:5432/dbname?sslmode=disable" go run ./cmd/queue
```

## Check The CLI

```bash
go run ./cmd/queue
```

Expected output:

```text
queue CLI is ready
postgres connection OK
```

## Run Migrations

```bash
go run ./cmd/queue migrate
```

Expected output:

```text
migration complete
```

The migration creates:

- `job_state` enum
- `jobs` table
- `jobs_claimable_idx` for finding runnable jobs
- `jobs_stuck_running_idx` for finding old running jobs during recovery

## Queue Schema

The `jobs` table stores both the job payload and the queue bookkeeping data.

Important columns:

- `id`: unique job ID
- `payload`: JSON job data
- `state`: current job state
- `attempts`: number of times the job has been tried
- `max_attempts`: maximum tries before the job is considered dead
- `run_at`: earliest time this job may run
- `locked_by`: worker that claimed the job
- `locked_at`: time the worker claimed the job
- `last_error`: latest failure message
- `created_at`: insert time
- `updated_at`: latest update time

Allowed job states:

```text
pending
running
done
failed
dead
```

## Inspect The Database

After running migrations, inspect the table with:

```bash
docker exec "postgres-job-queue" psql -U queue -d queue -c "\d jobs"
```

## Project Structure

```text
.
|-- .dockerignore                  # Docker build context excludes
|-- cmd/queue/main.go              # CLI entrypoint
|-- Dockerfile                     # queue CLI container image
|-- docker-compose.yml             # Local PostgreSQL service
|-- internal/queue/queue.go         # Database wrapper
|-- internal/queue/worker.go        # Worker placeholder
|-- k8s/                            # Kubernetes manifests and notes
|-- migrations/001_create_jobs.sql  # Queue schema migration
|-- migrations/migrations.go        # Embedded migration SQL
|-- milestones.md                   # App learning milestones
`-- k8s-milestones.md              # Kubernetes learning milestones
```

## Development Checks

Run package build/test verification:

```bash
go test ./...
```

## Current CLI

```bash
go run ./cmd/queue
go run ./cmd/queue migrate
```

## Next Milestones

Milestone 1 and Milestone 2 are complete. Kubernetes Milestones K1, K2, K3, K4, and K5 are also complete.

The recommended next path is:

```text
K6: Migration Job
K7: Postgres StatefulSet
K8: One-Off CLI Jobs
```

This lets Kubernetes run the app functionality that already exists: `queue migrate`.

After K6, return to app Milestone 3, which will add commands to insert jobs:

```bash
queue enqueue '{"type":"email","to":"a@example.com"}'
queue seed --count=100
```

Those commands will insert rows into `jobs` with `state = pending`.

## License

MIT License. See `LICENSE`.
