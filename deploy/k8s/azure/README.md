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