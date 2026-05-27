# Kubernetes Learning Milestones With kind

This track runs alongside `milestones.md`.

The app milestones teach the job queue itself. These Kubernetes milestones teach how to run that queue inside a local kind cluster.

Current source-of-truth note:

```text
K2 through K7 were originally learned with raw Kubernetes YAML.
After K8, the app Kubernetes resources live in the Helm chart under charts/postgres-job-queue.
k8s/ now keeps only the local kind cluster config and Kubernetes notes.
```

Recommended timing:

```text
App Milestone 1: stay local with Docker Compose
App Milestone 2: start Kubernetes
App Milestone 7: make Kubernetes the main runtime for workers
App Milestone 10: verify the whole system in Kubernetes
```

## Milestone K1: Local Cluster

Goal: create a local kind cluster and prove `kubectl` can talk to it.

Create:

```text
k8s/kind-config.yaml
```

Start simple:

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: queue
nodes:
  - role: control-plane
```

Create the cluster:

```bash
kind create cluster --config k8s/kind-config.yaml
kubectl cluster-info
kubectl get nodes
```

Done when:

```bash
kubectl get nodes
```

shows at least one `Ready` node.

Memory Box:

```text
kind runs Kubernetes nodes as Docker containers.
Your cluster is local, disposable, and good for learning.
```

## Milestone K2: Namespace

Goal: isolate the project resources in their own namespace.

Create:

```text
k8s/namespace.yaml
```

Resource:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: postgres-job-queue
```

Apply:

```bash
kubectl apply -f k8s/namespace.yaml
kubectl get namespaces
```

Done when the namespace exists.

## Milestone K3: Container Image

Goal: package the Go CLI as a Docker image Kubernetes can run.

Create:

```text
Dockerfile
.dockerignore
```

The image should contain the compiled `queue` binary.

Build locally:

```bash
docker build -t postgres-job-queue:dev .
```

Load the image into the kind cluster:

```bash
kind load docker-image postgres-job-queue:dev --name queue
```

Use this image in Kubernetes manifests:

```yaml
image: postgres-job-queue:dev
imagePullPolicy: IfNotPresent
```

Done when Kubernetes can use the image from kind's local node image store without pulling from a remote registry.

Memory Box:

```text
Your host Docker image store and the kind node image store are not the same thing.
After docker build, run kind load docker-image so pods can use the image.
```

## Milestone K4: Postgres In Kubernetes

Goal: run Postgres inside the kind cluster.

Create:

```text
k8s/postgres-secret.yaml
k8s/postgres-pvc.yaml
k8s/postgres-deployment.yaml
k8s/postgres-service.yaml
```

Resources:

```text
Secret: stores POSTGRES_USER and POSTGRES_PASSWORD
PersistentVolumeClaim: stores Postgres data
Deployment: runs the postgres container
Service: gives Postgres a stable DNS name
```

Use a Deployment first. The StatefulSet switch comes in K7 after the basic in-cluster database path works.

Use this service name for in-cluster connections:

```text
postgres
```

The in-cluster database URL should look like:

```text
postgres://queue:queue@postgres:5432/queue?sslmode=disable
```

Apply:

```bash
kubectl apply -n postgres-job-queue -f k8s/postgres-secret.yaml
kubectl apply -n postgres-job-queue -f k8s/postgres-pvc.yaml
kubectl apply -n postgres-job-queue -f k8s/postgres-deployment.yaml
kubectl apply -n postgres-job-queue -f k8s/postgres-service.yaml
```

Inspect:

```bash
kubectl get pods -n postgres-job-queue
kubectl get svc -n postgres-job-queue
```

Done when the Postgres pod is `Running` and the service exists.

Note for kind:

```text
The PersistentVolumeClaim uses storage inside the kind node container.
It survives pod restarts, but it is still local learning storage.
If you delete the kind cluster, the data is gone.
```

## Milestone K5: Database URL Configuration

Goal: give the app its database URL through Kubernetes configuration instead of hardcoding it.

Create:

```text
k8s/app-secret.yaml
```

Store:

```text
DATABASE_URL=postgres://queue:queue@postgres:5432/queue?sslmode=disable
```

Use it from pods as an environment variable:

```yaml
envFrom:
  - secretRef:
      name: queue-app
```

Done when a pod can receive `DATABASE_URL` from Kubernetes.

## Milestone K6: Migration Job

Prerequisite: app Milestone 2 is complete.

Goal: run `queue migrate` as a Kubernetes Job.

Create:

```text
k8s/migrate-job.yaml
```

The job should run:

```bash
queue migrate
```

Apply:

```bash
kubectl apply -n postgres-job-queue -f k8s/migrate-job.yaml
```

Inspect:

```bash
kubectl get jobs -n postgres-job-queue
kubectl get pods -n postgres-job-queue
kubectl logs -n postgres-job-queue job/queue-migrate
```

Done when the Job reaches `Complete` and the schema exists in Postgres.

Memory Box:

```text
Kubernetes Job = run this command until it succeeds once.
queue migrate is a perfect Kubernetes Job.
```

## Milestone K7: Postgres StatefulSet

Prerequisite: K4, K5, and K6 are complete.

Goal: replace the learning Postgres Deployment with a StatefulSet before adding real queue data and workers.

Why now:

```text
First prove the app can reach Postgres and run migrations.
Then switch Postgres to the database-shaped Kubernetes object.
Do this before seed/stats/workers so resetting local learning data is cheap.
```

Create:

```text
k8s/postgres-headless-service.yaml
k8s/postgres-statefulset.yaml
```

Keep:

```text
k8s/postgres-service.yaml
```

Resources:

```text
Headless Service: gives StatefulSet pods stable DNS identities
StatefulSet: keeps postgres-0 alive with a stable identity
volumeClaimTemplates: gives postgres-0 its own stable PVC
Service: keeps the app connection name as postgres
```

Use this split:

```text
postgres-headless: used by the StatefulSet for pod identity
postgres: used by the app as the database hostname
```

The app database URL should stay:

```text
postgres://queue:queue@postgres:5432/queue?sslmode=disable
```

Apply:

```bash
kubectl delete deployment/postgres -n postgres-job-queue
kubectl apply -n postgres-job-queue -f k8s/postgres-headless-service.yaml
kubectl apply -n postgres-job-queue -f k8s/postgres-statefulset.yaml
kubectl rollout status statefulset/postgres -n postgres-job-queue
```

Inspect:

```bash
kubectl get statefulsets -n postgres-job-queue
kubectl get pods -n postgres-job-queue
kubectl get pvc -n postgres-job-queue
```

If the switch creates a fresh local database, rerun K6 so the schema exists again.

Done when `postgres-0` is `Running`, a StatefulSet exists, and Postgres has a per-pod PVC.

Memory Box:

```text
Deployment asks: do I have enough copies?
StatefulSet asks: do I have the right named copies with their own storage?
```

## Milestone K8: Package The Current YAML With Helm

Prerequisite: K1 through K7 are complete.

Goal: learn Helm by converting the Kubernetes objects you already understand into a small chart.

Why this detour is useful now:

```text
Production Kubernetes often uses Helm to install, upgrade, roll back, and inspect releases.
Learning Helm now helps you understand how raw Kubernetes YAML becomes a managed release.
```

Important rule:

```text
Helm should not hide Kubernetes from you.
The first chart should render almost the same objects as the existing YAML.
```

Create:

```text
charts/postgres-job-queue/Chart.yaml
charts/postgres-job-queue/values.yaml
charts/postgres-job-queue/templates/
charts/postgres-job-queue/templates/_helpers.tpl
```

Move or copy the current resource shapes into `templates/`:

```text
postgres Secret
app Secret
postgres Service
postgres Headless Service
postgres StatefulSet
migration Job
```

Start with a small `values.yaml`:

```text
app image repository and tag
postgres image repository and tag
database name
database user
database password
database URL
storage size
```

Do not template everything on the first pass. Template only values that you expect to change between environments.

Render without installing:

```bash
helm template queue ./charts/postgres-job-queue
```

Compare the rendered manifests with the current YAML:

```bash
helm template queue ./charts/postgres-job-queue > /tmp/queue-rendered.yaml
kubectl apply --dry-run=client -f /tmp/queue-rendered.yaml
```

Install into the local kind cluster:

```bash
helm install queue ./charts/postgres-job-queue --namespace postgres-job-queue --create-namespace
```

Inspect the release:

```bash
helm list -n postgres-job-queue
helm status queue -n postgres-job-queue
helm get values queue -n postgres-job-queue
helm get manifest queue -n postgres-job-queue
```

Practice an upgrade:

```bash
helm upgrade queue ./charts/postgres-job-queue -n postgres-job-queue
helm history queue -n postgres-job-queue
```

Practice a rollback:

```bash
helm rollback queue 1 -n postgres-job-queue
```

Useful debugging drills:

```text
change the app image tag to a missing image and inspect ImagePullBackOff
change the database password and inspect failed startup or failed migration behavior
change the storage size and inspect what Kubernetes accepts or rejects
change a template label and inspect Service selector breakage
delete a Helm-managed object manually and run helm upgrade again
```

Done when:

```text
helm template renders valid Kubernetes manifests
helm install creates the working local cluster resources from the chart
helm get manifest lets you inspect the YAML that Helm installed
helm upgrade changes an installed release
helm rollback restores an earlier release revision
you can still explain every rendered Kubernetes object without relying on Helm magic
```

Memory Box:

```text
Helm is a release manager that renders Kubernetes YAML from templates and values.
In production debugging, always check both sides: Helm release state and Kubernetes object state.
```

## Milestone K9: One-Off CLI Jobs

Prerequisite: app Milestones 3 and 4 are complete.

Goal: run short CLI commands inside Kubernetes.

Create jobs for:

```bash
queue seed --count=100
queue stats
```

Possible chart templates:

```text
charts/postgres-job-queue/templates/seed-job.yaml
charts/postgres-job-queue/templates/stats-job.yaml
```

Install or upgrade the release:

```bash
helm upgrade --install queue ./charts/postgres-job-queue --namespace postgres-job-queue --create-namespace
kubectl logs -n postgres-job-queue job/queue-seed
```

For stats:

```bash
helm upgrade queue ./charts/postgres-job-queue -n postgres-job-queue
kubectl logs -n postgres-job-queue job/queue-stats
```

Done when jobs can be inserted and inspected from inside the cluster.

## Milestone K10: Worker Deployment

Prerequisite: app Milestone 7 is complete.

Goal: run workers continuously in the kind cluster.

Create:

```text
charts/postgres-job-queue/templates/worker-deployment.yaml
```

The container should run:

```bash
queue work --workers=1
```

Start with one worker process per pod.

Scale with Kubernetes replicas instead of only using the app's `--workers` flag.

Upgrade the release:

```bash
helm upgrade queue ./charts/postgres-job-queue -n postgres-job-queue
```

Scale:

```bash
kubectl scale deployment/queue-worker -n postgres-job-queue --replicas=5
```

Inspect:

```bash
kubectl get pods -n postgres-job-queue
kubectl logs -n postgres-job-queue deployment/queue-worker
```

Done when multiple worker pods process jobs concurrently.

Memory Box:

```text
Deployment = keep this long-running process alive.
Job = run this command to completion.
```

## Milestone K11: Kubernetes Scaling Test

Prerequisite: app Milestones 5, 6, and 7 are complete.

Goal: prove database locking still works when Kubernetes scales workers.

Run:

```bash
kubectl scale deployment/queue-worker -n postgres-job-queue --replicas=10
```

Verify:

```text
workers run in different pods
jobs are claimed once
no two pods process the same job at the same time
stats eventually show done/dead jobs increasing
```

Useful commands:

```bash
kubectl get pods -n postgres-job-queue -w
kubectl logs -n postgres-job-queue deployment/queue-worker --follow
helm upgrade queue ./charts/postgres-job-queue -n postgres-job-queue
kubectl logs -n postgres-job-queue job/queue-stats
```

Done when scaling replicas increases throughput without duplicate processing.

## Milestone K12: Health And Restart Behavior

Goal: learn how Kubernetes reacts when containers fail.

Add to the worker Deployment:

```text
resource requests and limits
readiness probe if the app exposes one
liveness probe only if the app can safely report stuck/unhealthy state
```

Test failure:

```bash
kubectl delete pod -n postgres-job-queue -l app=queue-worker
```

Verify:

```text
Kubernetes creates replacement pods
running jobs from killed pods are eventually recovered by app logic
workers continue processing after restart
```

Done when worker pods can be killed and the system recovers.

## Milestone K13: Stuck Job Recovery In Kubernetes

Prerequisite: app Milestone 9 is complete.

Goal: choose how recovery should run in the cluster.

Choose one design:

```text
Option A: recovery sweeper inside every worker process
Option B: one separate recovery Deployment
Option C: Kubernetes CronJob that runs queue recover periodically
```

Recommended for learning:

```text
Start with Option C: CronJob
```

Create:

```text
charts/postgres-job-queue/templates/recover-cronjob.yaml
```

The CronJob should run:

```bash
queue recover
```

Inspect:

```bash
kubectl get cronjobs -n postgres-job-queue
kubectl get jobs -n postgres-job-queue
kubectl logs -n postgres-job-queue job/<recover-job-name>
```

Done when abandoned `running` jobs are returned to the retry flow or marked dead.

## Milestone K14: Full Cluster Verification

Prerequisite: app Milestone 10 is complete.

Goal: prove the full queue works inside Kubernetes.

Run the whole flow inside the cluster:

```bash
helm upgrade --install queue ./charts/postgres-job-queue --namespace postgres-job-queue --create-namespace
kubectl scale deployment/queue-worker -n postgres-job-queue --replicas=5
kubectl logs -n postgres-job-queue job/queue-stats
```

Verify:

```text
Postgres runs in the cluster
migrations complete successfully
seed jobs insert queue rows
worker pods process jobs concurrently
failed jobs retry later
stuck jobs are recovered
jobs eventually become done or dead
```

Done when the local Docker Compose flow has a working Kubernetes equivalent.

## Milestone K15: Cleanup And Repeatability

Goal: make the Kubernetes setup easy to recreate.

Create:

```text
k8s/README.md
```

Document:

```text
how to create the cluster
how to build and load the image
how to render and install the Helm chart
how to run migrations
how to seed jobs
how to start workers
how to inspect stats
how to clean up
```

Delete project resources without deleting the whole kind cluster:

```bash
kubectl delete namespace postgres-job-queue
```

Delete the whole kind cluster:

```bash
kind delete cluster --name queue
```

Done when the whole Kubernetes environment can be deleted and recreated from documented commands.

## Future Milestones

After the local Kubernetes version works, explore production-shaped improvements:

```text
Kustomize overlays
external managed Postgres instead of in-cluster Postgres
HorizontalPodAutoscaler
PodDisruptionBudget
NetworkPolicy
ServiceAccount and RBAC
structured logs
Prometheus metrics
Grafana dashboard
OpenTelemetry tracing
GitHub Actions deployment
Argo CD or Flux GitOps
```

Do these later. They are more useful after the basic Kubernetes objects feel normal.
