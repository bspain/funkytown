# funkytown
A Distributed System Proof of Concept application

This project was designed to deal with a challenge in my company.  During a daily release of our core website, the dev team wanted 100's of Functional Specs (a mix of Webdriver, Cypress, and Playwright specs) to be executed "as fast as possible".  My idea was to implement this, using the ["Fine Parallel Processing Using a Work Queue" pattern](https://kubernetes.io/docs/tasks/job/fine-parallel-processing-work-queue/) in Kubernetes.

### The **Controller** application
The first part of the application is the `controller` pod, which hosts the **workqueue** in a local REDIS store.  It will load the list of specs and enqueue a **task** object for each **spec** / **browser** / **viewport** combination.  It will **filter out invalid combinations** (_such as **mobile** / **firefox**_). The `controller` will poll the **workqueue**, and update statistics as tasks finish.

The `controller` also hosts a `/results` endpoint (at port `:3000`), which is exposed via a **LoadBalancer** as the `reporter-service` (at port `:80`) which can be used to visualize all the **tasks** in the **workqueue**, their results, and the overall statistics.

The `controller` pod remains until terminated.

### The **Worker** application
The second part of the application is the `worker` pods, which are executed as a **batch Job** resource.  One to N `worker` pods pop **tasks** from the **workqueue** and execute the spec accordingly.  Each `worker` will then push a **result** back into the **workqueue**.  Once all `worker` pods exit, the batch job will be complete.


```mermaid
flowchart
    subgraph K["kubernetes"]
        CS[/controller start/] --> PL
        subgraph C ["controller - pod"]
            PL(parse list) --> PST(enqu tasks)
            PL --> ML(monitor loop)
            PST --->|task| RDB[(redis)]
            ML <-->|tasks remaining?| RDB
            PL --> HR(host reporter)
        end
        subgraph R ["workqueue service"]
            RC(ClusterIP) -->|redis commands| RDB
        end
        WS[/worker start/] --> ML2
        subgraph W ["worker(s) - job"]
            ML2(worker loop) <-->|task available?| RC
            ML2 -->|task| POT(dequ task)
            POT --> RS(run spec)
            RS --> PR(push result)
            PR -->|result| RC
            ML2 -->|tasks finished| WE[\worker exit\]
        end
        subgraph S ["reporter service"]
            HR --> RPS(LoadBalancer)
        end
    end
```

## Local Testing 
The `controller` and the `worker` can both be run locally using a docker hosted REDIS instance.

### Start local REDIS
Use docker to start a REDIS instance at localhost:6379
```
docker run -d -p 6379:6379 redislabs/redismod
```

> TIP: Using the `redis-cli`, you can easily reset the entire redis database using `FLUSHALL`

### Initialize the go workspace
This project uses `go work` features of golang `1.18` : https://go.dev/blog/get-familiar-with-workspaces

Create a workspace
```
go work init ./controller ./worker ./shared
```

### Start the controller
Start the controller providing the `REDIS_HOST`, `REDIS_PORT`, and `GROUP_TASKS_FILE`
```
REDIS_HOST=localhost REDIS_PORT=6379 GROUP_TASKS_FILE=specs/spec_context_map.json HTML_INDEX_FILE=controller/html/index.html go run github.com/bspain/funkytown/controller
```

Should be able to use the `redis-cli` to verify the "run metatdata" object was created sucessfully.
```
redis-cli
127.0.0.1:6379> hgetall runmeta
1) "runid"
2) "a_new_run"
3) "cmdcount"
4) "0"
5) "cmdfinishedcount"
6) "0"
7) "finished"
8) "0"
127.0.0.1:6379> exit
```

### Start the worker
Start the worker providing the `REDIS_HOST`, `REDIS_PORT`, and `SPEC_ROOT`
```
REDIS_HOST=localhost REDIS_PORT=6379 SPEC_ROOT=specs go run github.com/bspain/funkytown/worker
```

## Local Docker Development

### <a name="ldd-build"></a>Build the application images
Build the `controller` image
```
docker build -f Dockerfile.controller -t controller:latest .
```

Build the `worker` image
```
docker build -f Dockerfile.worker -t worker:latest .
```

### Create the local docker network
The `controller` and `worker` instances will be communicating with each other, therefore they need a local network service.

```
docker network create funkytown
```

### Start the controller
Start the `controller` image

```
docker run -it --rm --name controller --net funkytown -p 6379:6379 -p 80:3000 controller:latest
```

### Start the worker
Start the `worker` image

```
docker run -it --name worker --net funkytown --ipc=host worker:latest
```

# Kubernetes deployment

This POC contains information about local K8's deployment testing (using Windows Subsystem Linux 2 and Docker for Desktop)

See [`deploy/k8s/local/README.md`](deploy/k8s/local/README.md) for details


TODO: This POC contains information about Azure deployment, using Azure Container Registry (ACR) and Azure Kubernetes Service (AKS)

See TODO
