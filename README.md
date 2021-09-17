# K8s Operator

This is a simple k8s operator written in Go with operator framework. The idea is to watch a custom resource then make an nginx deployment with 

* Number of replicas.
* Message to be served via nginx.

## Deployment

### Prerequisties

* Minikube
* Go@1.16
* docker

### Create Your Image

To deploy the operator from the repo with your own docker image use

``` bash
./build-deploy-operator-on-custer.sh example/k8s-operator:1.0
```

This script will build and publish image, then deploy the operator in the current loaded kubeconfig. **P.S. This script login to docker using .docker-secret**

Next steps will be in [Testing Section](#testing)


### Use Existing Image

To deploy the operator with an existing image run

``` bash
make deploy IMG=mustafasaber/custom-k8s-operator:v9.1.0
```

### Deleting

To delete the deployment just run

``` bash
make undeploy
```

## Testing

To test the operator we apply a yaml file with the CR like this

``` bash
k apply -f ./config/samples/_v1_onekind.yaml
```

And now you can have nginx service up and running, expose the service with

``` bash
minikube service onekind-sample
```