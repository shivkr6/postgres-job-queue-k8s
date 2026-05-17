# Postgres Job Queue Coding Milestones

## Milestone 1: Project Skeleton

Goal: create the Go project structure and prove the CLI can compile and connect to Postgres.

Files:

```text
postgres-job-queue/
  go.mod
  docker-compose.yml
  migrations/001_create_jobs.sql
  cmd/queue/main.go
  internal/queue/queue.go
  internal/queue/worker.go
```

Done when:

```bash
docker compose up -d
go run ./cmd/queue
```

The CLI should print that the Postgres connection is OK.

## Milestone 2: Migration

Goal: create the database objects needed by the queue.

Implement:

```bash
queue migrate
```

Create:

```text
job_state enum
jobs table
indexes for claimable jobs
indexes for stuck running jobs
```

Done when `queue migrate` creates the schema successfully.

## Milestone 3: Enqueue and Seed

Goal: put jobs into the table.

Implement:

```bash
queue enqueue '{"type":"email","to":"a@example.com"}'
queue seed --count=100
```

Functions:

```go
Enqueue(ctx, payload)
Seed(ctx, count)
```

Done when jobs are inserted with `state = pending`.

## Milestone 4: Stats

Goal: inspect the queue state.

Implement:

```bash
queue stats
```

Function:

```go
Stats(ctx)
```

Example output:

```text
pending: 100
running: 0
done: 0
failed: 0
dead: 0
```

Done when the CLI can show counts for every job state.

## Milestone 5: Claim

Goal: safely let one worker take one available job.

Implement:

```go
Claim(ctx, workerID)
```

The claim query must use:

```sql
FOR UPDATE SKIP LOCKED
```

A claimed job changes:

```text
pending -> running
```

Done when one worker can claim a job and no other worker can claim that same job while it is running.

## Milestone 6: Complete and Fail

Goal: finish jobs after processing.

Implement:

```go
Complete(ctx, jobID, workerID)
Fail(ctx, jobID, workerID, err)
```

Successful job:

```text
running -> done
```

Failed job with attempts left:

```text
running -> failed
```

Failed job with no attempts left:

```text
running -> dead
```

Done when workers can mark claimed jobs as done, failed, or dead.

## Milestone 7: Worker Loop

Goal: run workers continuously.

Implement:

```bash
queue work --workers=5
```

Each worker loop:

```text
claim job
if no job, sleep briefly
process fake job
complete or fail
repeat
```

Done when 5 worker goroutines process jobs concurrently.

## Milestone 8: Retry Scheduling

Goal: retry failed jobs later instead of immediately.

Failed jobs get:

```text
state = failed
run_at = now + delay
```

The claim query should treat these as claimable:

```text
pending jobs where run_at <= now
failed jobs where run_at <= now
```

Done when failed jobs retry after a delay.

## Milestone 9: Stuck Job Recovery

Goal: recover jobs abandoned by crashed workers.

Implement:

```go
RecoverStuck(ctx, timeout)
```

CLI command:

```bash
queue recover
```

Integrate recovery into:

```bash
queue work --workers=5
```

The worker process should start:

```text
5 workers
1 recovery sweeper
```

Done when old `running` jobs are moved back into retry flow or marked dead.

## Milestone 10: Verification

Goal: prove the queue works correctly.

Run:

```bash
queue migrate
queue seed --count=100
queue work --workers=5
queue stats
```

Verify:

```text
5 workers run at the same time
no two workers process the same job at the same time
failed jobs retry
jobs eventually become done or dead
stuck jobs are recovered
```
