# Kubernetes Setup

This directory contains the Kubernetes files for running the project in a local kind cluster.

Completed Kubernetes milestones:

- K1: create a local cluster and prove `kubectl` can talk to it
- K2: create a project namespace

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

## Next Milestone

Next is K3: build a Docker image containing the `queue` binary and load it into the kind cluster.

## Memory Box

```text
kind creates the local Kubernetes cluster.
kubectl talks to the selected cluster context.
A Ready node means Kubernetes has a place to run containers.
```
