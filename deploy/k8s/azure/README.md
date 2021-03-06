# Deployment to Azure Kubernetes Service (AKS)

For this POC, will assume access to an Azure account, and working knowledge of the Azure CLI

## Creating the Azure Container Registry (ACR)

> These steps are detailed in [Kubernetes on Azure tutorial - Create a container registry - Azure Kubernetes Service | Microsoft Docs](https://docs.microsoft.com/en-us/azure/aks/tutorial-kubernetes-prepare-acr?tabs=azure-cli)

1. Creating a Resource Group
    ```
    az group create --name funkytown_rg --location centralus
    ```

2. Create a Container Service
    ```
    az acr create --resource-group funkytown_rg --name funkytownacr --sku Basic
    ```

## Setup Github Actions workflow to publish to ACR
For this POC, I wanted to configure Github Actions for the `funkytown` repo to build both `controller` and `worker` images, and publish them directly in my new ACR: `funkytownacr`

> These steps are detailed in [Deploy container instance by GitHub action - Azure Container Instances | Microsoft Docs](https://docs.microsoft.com/en-us/azure/container-instances/container-instances-github-action#configure-github-workflow)

1. Create a Service Principal
    TODO: `--sdk-auth` claims it is deprecated, what is the new integration story with GH Actions? (specifically, [the `azure/login@v1` action](https://docs.microsoft.com/en-us/azure/container-instances/container-instances-github-action#create-workflow-file))

    ```
    # $groupId from funkytown_rg
    az ad sp create-for-rbac --scope $groupId --role Contributor --sdk-auth
    ```

2. Get RegistryId 
    ```
    registryId=$(az acr show --name funkytownacr --query id --output tsv)
    ```

3. Assign the `AcrPush` role to the Service Principal from Step 1.  `clientId` comes from the output of Step 1.
    ```
    # $clientId from Service Principal creation step
    az role assignment create --assignee $clientId --scope $registryId --role AcrPush
    ```

4. Setup the Repo Secrets in Github

    Key | Value
    -- | --
    AZURE_CONTAINER_SDKAUTH	| _Full JSON output_ of SDKAuth in Step 1 
    REGISTRY_LOGIN_SERVER | `funkytownacr.azurecr.io`
    REGISTRY_USERNAME | `clientID` from SDKAuth in Step 1
    REGISTRY_PASSWORD | `clientSecret` from SDKAuth in Step 1
    RESOURCE_GROUP | `funkytown_rg`

5. Create action as shown at [`.github/workflows/publish.yml`](../../../.github/workflows/publish.yml)

## Create Azure Kubernetes Service
Following along from [Kubernetes on Azure tutorial - Deploy a cluster - Azure Kubernetes Service | Microsoft Docs](https://docs.microsoft.com/en-us/azure/aks/tutorial-kubernetes-deploy-cluster?tabs=azure-cli)


1. Create the AKS cluster
    ```
    az aks create --resource-group funkytown_rg --name funkytownaks --node-count 2 --generate-ssh-keys --attach-acr funkytownacr
    ```

2. Get AKS credentials.  This will directly setup certificate config in the Cloud shell `~/.kube/config`
    ```
    az aks get-credentials --resource-group funkytown_rg --name funkytownaks
    ```

## Deploy Application Resources
For this POC, I created **individual resource .yaml files** in `deploy/k8s/azure` - a typical deployment strategy would use one file, or orchestrate the application deployment with Helm

During the POC work, I found it helpful to develop the **resource.yaml** files locally, and then use git to 'push/pull' those files in the Azure CLI shell.

> TODO: I'm sure there are many more interesting ways to achieve this smoother.  Github codespaces perhaps?

1. Clone the `funkytown` repo into the Azure CLI shell
    ```
    git clone https://github.com/bspain/funkytown.git
    ```

2. Deploy the `controller` pod
    ```
    cd deploy/k8s/azure

    kubectl apply -f deploy/k8s/azure/controller.yaml
    ```

3. Deploy the `reporter` service
    ```
    kubectl apply -f reporter-service.yaml
    ```

    The LoadBalancer will eventually receive an external IP address.  Once that happens, the `/results` endpoint should be live at the IP.

    Run `get service`, and wait for the `EXTERNAL-IP` to get assigned.
    ```
    kubectl get service funkytown-
    reporter-service
    ```

4. Deploy the `workqueue` service

    ```
    kubectl apply -f workqueue-service.yaml
	```

    As an optional step, you can test that the workqueue service is available within the cluster using a temp/redis pod

    ```
	kubectl run -i --tty temp --im
	age redis --command "/bin/sh"
	
    If you don't see a command prompt, try pressing enter.
	
	# redis-cli -h funkytown-workqueue-service -p 6379
	funkytown-workqueue-service:6379> hgetall runmeta
	1) "runid"
	2) "a_new_run"
	3) "tasksremaining"
	4) "12"
	5) "tasksfinished"
	6) "0"
	7) "finished"
	8) "0"
    ```

5. Deploy the `worker` as a Job

    ```
    kubectl create -f deploy/k8s/azure/worker-job.yaml
    ```

    The worker should start working, execute the job, and the pod should shut down once the job is finished.

## Diagnostic: Confirm the worker can execute playwright specs within the cluster

An optional diagnostic step is to run the `worker-pod-debug` spec, and then execute a spec from within

```
kubectl create -f deploy/k8s/azure/worker-pod-debug.yaml
```

After the pod is created, exec into the pod, move in the `/specs` path, and run a spec
```
kubectl exec -it funkytown-worker-debug bash

root@funkytown-worker-debug:/app# 
cd specs

root@funkytown-worker-debug:/app/specs#
npx playwright test --project=desktop-chrome find_a_store

Running 1 test using 1 worker
```

Delete the debug pod with `kubectl delete po funkytown-worker-debug`