# Kubernetes Setup

This directory now contains only the local kind cluster setup for the project.

The Kubernetes app resources are managed by Helm from:

```text
charts/postgres-job-queue
```

Completed Kubernetes milestones:

- K1: create a local cluster and prove `kubectl` can talk to it
- K2: create a project namespace
- K3: build and load the local app image
- K4: run Postgres inside the kind cluster
- K5: provide the app database URL through Kubernetes configuration
- K6: run database migrations as a Kubernetes Job
- K7: run Postgres as a StatefulSet with per-pod storage
- K8: package the Kubernetes app resources with Helm

## Requirements

- Docker
- kind
- kubectl
- Helm

## Current Files

```text
k8s/
|-- kind-config.yaml  # Local kind cluster configuration
`-- README.md         # Kubernetes setup notes
```

Helm chart:

```text
charts/postgres-job-queue/
|-- Chart.yaml
|-- values.yaml
`-- templates/
    |-- _helpers.tpl
    |-- app-secret.yaml
    |-- migrate-job.yaml
    |-- postgres-headless-service.yaml
    |-- postgres-secret.yaml
    |-- postgres-service.yaml
    `-- postgres-statefulset.yaml
```

Memory Box:

```text
k8s/kind-config.yaml creates the local learning cluster.
charts/postgres-job-queue is the source of truth for app Kubernetes resources.
```

## Create The Local Cluster

The cluster is created from:

```text
k8s/kind-config.yaml
```

Current config:

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
```

This creates a local Kubernetes cluster named `queue`.

In kind, Kubernetes nodes are Docker containers. The single `control-plane` node is enough for this learning setup because it can also run normal workloads.

## Verify The Cluster

Check which Kubernetes context `kubectl` is using:

```bash
kubectl config current-context
```

Expected context:

```text
kind-queue
```

Show cluster API information:

```bash
kubectl cluster-info
```

Show cluster nodes:

```bash
kubectl get nodes
```

K1 is complete when the node is `Ready`:

```text
NAME                  STATUS   ROLES           VERSION
queue-control-plane   Ready    control-plane   v1.35.0
```

Useful cluster commands:

```bash
kind get clusters
kubectl config use-context kind-queue
kind delete cluster --name queue
```

## Container Image

The Docker image is built from the project root `Dockerfile`.

Build the local image:

```bash
docker build -t postgres-job-queue:dev .
```

Load it into the kind cluster:

```bash
kind load docker-image postgres-job-queue:dev --name queue
```

Verify the image exists inside the kind node:

```bash
docker exec "queue-control-plane" crictl images postgres-job-queue
```

The Helm chart defaults to:

```yaml
app:
  image:
    repository: postgres-job-queue
    tag: dev
    pullPolicy: IfNotPresent
```

If the Postgres image is already available locally and the kind node cannot pull from Docker Hub, load it too:

```bash
kind load docker-image postgres:16-alpine --name queue
```

## Helm Chart

K8 packages the current Kubernetes app resources as a Helm chart.

The chart renders:

```text
Secret/postgres
Secret/queue-app
Service/postgres
Service/postgres-headless
StatefulSet/postgres
Job/queue-migrate
```

The chart does not render a Namespace. Create the namespace outside the chart or let Helm create it during install with `--create-namespace`.

Render the chart without installing it:

```bash
helm lint ./charts/postgres-job-queue
helm template queue ./charts/postgres-job-queue > /tmp/queue-rendered.yaml
```

Validate the rendered Kubernetes manifests:

```bash
kubectl apply --dry-run=client -f /tmp/queue-rendered.yaml
```

Install into the project namespace:

```bash
helm install queue ./charts/postgres-job-queue --namespace postgres-job-queue --create-namespace
```

If you want a clean test namespace, use:

```bash
helm install queue ./charts/postgres-job-queue --namespace postgres-job-queue-helm --create-namespace
```

Inspect the Helm release:

```bash
helm list -n postgres-job-queue
helm status queue -n postgres-job-queue
helm get values queue -n postgres-job-queue
helm get manifest queue -n postgres-job-queue
```

If you installed into `postgres-job-queue-helm`, use that namespace in the commands above.

Verify the Kubernetes resources:

```bash
kubectl get all,pvc -n postgres-job-queue
kubectl rollout status statefulset/postgres -n postgres-job-queue
kubectl wait -n postgres-job-queue --for=condition=complete job/queue-migrate --timeout=120s
kubectl logs -n postgres-job-queue job/queue-migrate
kubectl exec -n postgres-job-queue pod/postgres-0 -- psql -U queue -d queue -c "\d jobs"
```

The migration Job can retry if it starts before Postgres accepts connections. That is normal for this learning chart as long as the Job eventually reaches `Complete`.

Practice an upgrade that changes the StatefulSet Pod template:

```bash
helm upgrade queue ./charts/postgres-job-queue -n postgres-job-queue --set postgres.image.pullPolicy=Never
kubectl rollout status statefulset/postgres -n postgres-job-queue
helm history queue -n postgres-job-queue
```

Practice a rollback:

```bash
helm rollback queue 1 -n postgres-job-queue
kubectl rollout status statefulset/postgres -n postgres-job-queue
helm history queue -n postgres-job-queue
```

For local cleanup:

```bash
helm uninstall queue -n postgres-job-queue
kubectl delete namespace postgres-job-queue
```

## Debugging Split

Use Helm to inspect the release:

```bash
helm status queue -n postgres-job-queue
helm get values queue -n postgres-job-queue
helm get manifest queue -n postgres-job-queue
helm history queue -n postgres-job-queue
```

Use Kubernetes to inspect what the cluster actually did:

```bash
kubectl get pods -n postgres-job-queue
kubectl describe pod -n postgres-job-queue postgres-0
kubectl get events -n postgres-job-queue --sort-by=.lastTimestamp
kubectl logs -n postgres-job-queue job/queue-migrate
```

Memory Box:

```text
Helm renders Kubernetes YAML from chart templates and values.
helm get manifest shows the YAML Helm installed.
kubectl shows what the cluster actually did with that YAML.
Production debugging needs both views.
```

## Next Milestone

App Milestone 4 is complete: the CLI can show queue state counts with `queue stats`.

Next is app Milestone 5: add job claiming so one worker can safely take one available job.

The next Kubernetes milestone is K9: One-Off CLI Jobs.
