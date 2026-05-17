# Kubernetes Setup

This directory contains the Kubernetes files for running the project in a local kind cluster.

Completed Kubernetes milestones:

- K1: create a local cluster and prove `kubectl` can talk to it
- K2: create a project namespace
- K3: build and load the local app image

## Requirements

- Docker
- kind
- kubectl

## Current Files

```text
k8s/
|-- kind-config.yaml  # Local kind cluster configuration
|-- namespace.yaml    # Project namespace
`-- README.md         # Kubernetes setup notes
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

## Useful Commands

List kind clusters:

```bash
kind get clusters
```

Switch `kubectl` back to this cluster:

```bash
kubectl config use-context kind-queue
```

Delete the whole local cluster:

```bash
kind delete cluster --name queue
```

## Project Namespace

The project namespace is defined in:

```text
k8s/namespace.yaml
```

Apply it with:

```bash
kubectl apply -f k8s/namespace.yaml
```

Verify it with:

```bash
kubectl get namespace postgres-job-queue
```

Expected status:

```text
postgres-job-queue   Active
```

Future project resources should be created in this namespace with:

```bash
kubectl apply -n postgres-job-queue -f <manifest.yaml>
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

Kubernetes manifests should use:

```yaml
image: postgres-job-queue:dev
imagePullPolicy: IfNotPresent
```

## Next Milestone

Next is K4: run Postgres inside the kind cluster.

## Memory Box

```text
kind creates the local Kubernetes cluster.
kubectl talks to the selected cluster context.
A Ready node means Kubernetes has a place to run containers.
kind load docker-image copies a host Docker image into the kind node image store.
```
