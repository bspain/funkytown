# Helm Funkytown Workflow
This part of the POC uses Helm to create the full K8s deployment and trigger the spec execution job

## Install the funkytown application
Use `helm install <app name and version> ./helm` to create the application installation.
```
helm install funk-v1 .
```

The key element that **helm** adds, is that each resource will include the installation name (e.g. the `.Release.Name`) as part of it's `name`.  This allows for multiple installations of **funkytown** to exist side by side in the same cluster.


```
kubectl get po

NAME                 READY   STATUS    RESTARTS   AGE
funk-v1-controller   1/1     Running   0          3s
funk-v1-worker       1/1     Running   0          3s
```

> The use-case here is creating a separate instance of **funkytown** to test either a specific deployment of the Application Under Test, or to run an updated version of a spec.


## Uninstall the funkytown application 
Using the same `<app name and version>` that you created the application with, use `helm uninstall` to remove it.

```
helm uninstall funk-v1
```