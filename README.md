# The Kubernetes Cluster Info tool

This tool is a quick and dirty test of the `client-go` package (for Go clients talking to a Kkubernetes cluster), at the moment this tool will fetch the number of pods, deployments and services in each namespace and then print this information to the user.

API authentication method are _out of cluster_ and you need your kubeconfig file that holds the context information to initialize a client. The same file that `kubectl` uses.

More info about `client-go` and how to use it in various ways can be found [here](https://github.com/kubernetes/client-go).

This tool is not _done_ yet and it can only get better!

## Install

Use the provided `Makefile` and run `make build`, if you want to build a Linux binary you can run `make build-linux`.

## Run

Example:

```
./k8sinfo
Namespace: default
        Pods: 2
        Deployments: 2
        Services: 3
Namespace: kube-public
        Pods: 0
        Deployments: 0
        Services: 0
Namespace: kube-system
        Pods: 10
        Deployments: 3
        Services: 2

Cluster totals:
        Pods: 12
        Deployments: 5
        Services: 5
```