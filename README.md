# Postgres Job Queue

A small Go project for learning how to build a PostgreSQL-backed job queue from first principles.

The project currently implements:

- Milestone 1: Go CLI skeleton and PostgreSQL connection check
- Milestone 2: database migration for the queue schema
- Milestone 3: enqueue and seed commands for inserting pending jobs
- Milestone 4: stats command for inspecting queue state counts

Future milestones will add job claiming, workers, retries, stuck-job recovery, and Kubernetes deployment.

For the local Kubernetes learning setup, see [`k8s/README.md`](k8s/README.md).

## Requirements

- Go 1.25+
- Docker
- kind and kubectl for the Kubernetes milestones

## Database Target

The current learning path verifies the app against PostgreSQL inside the local kind cluster.

Inside Kubernetes, the database URL uses the Postgres Service name:

```text
postgres://queue:queue@postgres:5432/queue?sslmode=disable
```

When running the CLI from the host machine instead of from a Kubernetes Pod, the default app database URL is:

```text
postgres://queue:queue@localhost:5432/queue?sslmode=disable
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

## Enqueue Jobs

Insert one job:

```bash
queue enqueue '{"type":"email","to":"a@example.com"}'
```

Insert many learning jobs:

```bash
queue seed --count=100
```

Both commands insert rows into `jobs` and rely on the database defaults to set `state = pending`.

## Inspect Queue Stats

Show how many jobs are in each queue state:

```bash
queue stats
```

Example output:

```text
pending: 100
running: 0
done: 0
failed: 0
dead: 0
```

The stats output includes every allowed job state, even when the count is zero.

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

In the Kubernetes learning cluster, inspect the table through the Postgres StatefulSet Pod:

```bash
kubectl exec -n postgres-job-queue pod/postgres-0 -- psql -U queue -d queue -c "\d jobs"
```

## Project Structure

```text
.
|-- .dockerignore                  # Docker build context excludes
|-- charts/                         # Helm charts
|-- cmd/queue/main.go              # CLI entrypoint
|-- Dockerfile                     # queue CLI container image
|-- docker-compose.yml             # Optional local PostgreSQL service
|-- internal/queue/queue.go         # Database wrapper
|-- internal/queue/worker.go        # Worker placeholder
|-- k8s/                            # kind cluster config and Kubernetes notes
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
go run ./cmd/queue enqueue '{"type":"email","to":"a@example.com"}'
go run ./cmd/queue seed --count=100
go run ./cmd/queue stats
```

## Next Milestones

Milestones 1 through 4 are complete. Kubernetes Milestones K1 through K8 are also complete.

The recommended next path is:

```text
App Milestone 5: claim one available job
K9: One-Off CLI Jobs
```

K7 runs Postgres as a StatefulSet, which gives it a stable Pod identity and a per-pod PVC.

K8 packages the current Kubernetes resources as a Helm chart in `charts/postgres-job-queue`.

Next, app Milestone 5 will add the job-claiming query that safely moves one available job into `running`.

## License

MIT License. See `LICENSE`.
