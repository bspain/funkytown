# Deployment to Local Kubernetes

This POC is leveraging [Docker for Desktop with Kubernetes being hosted by the WSL2 engine](https://docs.docker.com/desktop/windows/wsl/) for local development

> NOTE: Browser based automation is very resource intensive, additional tuning may be required to allocate enough system resources to support local Kubernetes execution.  Recommend working with Docker from the base README.md to diagnose specific spec execution issues.


## Starting the application in local kubernetes (Windows WSL2 + Docker for Desktop)
If the `controller` image is not already built, [build it now](#ldd-build).  Otherwise, you should see it in your local docker registry

```
docker images

REPOSITORY  TAG     IMAGE ID       CREATED         SIZE
controller  latest  e4ea2e71c9fd   2 minutes ago   1.01GB
```

Confirm that `kubectl` is working on your system, and that your current context is set to `docker-desktop`
```
kubectl config get-contexts

CURRENT   NAME             CLUSTER          AUTHINFO         NAMESPACE
*         docker-desktop   docker-desktop   docker-desktop
```

Deploy the `controller` pod
```
kubectl create -f deploy/k8s/local/controller.yaml
```

You can confirm the controller is up and running by viewing the logs from it
```
kubectl logs funkytown-controller 

Starting redis...
Redis server v=5.0.14 sha=00000000:0 malloc=jemalloc-5.1.0 bits=64 build=20357126fbf912e0
9:C 28 May 2022 14:56:16.411 # oO0OoO0OoO0Oo Redis is starting oO0OoO0OoO0Oo
9:C 28 May 2022 14:56:16.411 # Redis version=5.0.14, bits=64, commit=00000000, modified=0, pid=9, just started
9:C 28 May 2022 14:56:16.411 # Configuration loaded
Starting controller...
2022/05/28 14:56:16 funkytown Controller has started...
```

## Deploy the controller hosted services (workqueue and reporter)

Deploy the `workqueue` and `reporter` services

```
kubectl create -f deploy/k8s/local/workqueue-service.yaml 
kubectl create -f deploy/k8s/local/reporter-service.yaml
```

Using the `redis-cli` you should be able to connect to the REDIS DB hosted by the `controller` at port `30379` (note, this is assuming Docker-for-Desktop on WSL2 which creates a `kubernetes.docker.internal` mapping to `127.0.0.1`)

```
redis-cli -h kubernetes.docker.internal -p 30379

kubernetes.docker.internal:30379> hgetall runmeta
1) "tasksfinished"
2) "0"
3) "finished"
4) "0"
5) "runid"
6) "a_new_run"
7) "tasksremaining"
8) "12"
```
> In "PROD" we would not expose the REDIS port

You should also be able to view the `reporter` in a browser at http://kubernetes.docker.internal:30000/results

## Diagnostic: Confirm the workqueue service is available within the cluster

Another optional diagnostic step here is to launch the `redis-cli` from within a new cluster pod, and confirm that the `funkytown-workqueue-service` is operational within the cluster.

```
kubectl run -i --tty temp --image redis --command "/bin/sh"

If you don't see a command prompt, try pressing enter.

# redis-cli -h funkytown-workqueue-service -p 40379
funkytown-workqueue-service:40379> hgetall runmeta
```

Delete your temp redis-cli pod using `kubectl delete po temp`

## Execute the worker as a Job

Deploy the `worker` job

```
docker build -f Dockerfile.worker -t worker:latest .
```

## Diagnostic: Confirm the worker can execute playwright specs within the cluster

An optional diagnostic step is to run the `worker-pod-debug` spec, and then execute a spec from within

```
kubectl create -f deploy/k8s/local/worker-pod-debug.yaml
```

After the pod is created, exec into the pod, move in the `/specs` path, and run a spec
```
k exec -it funkytown-worker-debug bash

root@funkytown-worker-debug:/app# 
cd specs

root@funkytown-worker-debug:/app/specs#
npx playwright test --project=desktop-chrome find_a_store

Running 1 test using 1 worker
```

Delete the debug pod with `kubectl delete po funkytown-worker-debug`
