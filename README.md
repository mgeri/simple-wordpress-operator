# Simple Wordpress Operator

An example of simple Kubernetes operator.

This operator was created using Operator-SDK 1.14

## Creating the Operator

### Create directory of your project

```shell
$ mkdir simple-wordpress-operator
$ cd simple-wordpress-operator
```

### Init the operator project

```shell
$ operator-sdk init --domain=simplewordpress.com --repo github.com/mgeri/simple-wordpress-operator
```

### Create the CRDs and Controller for your operator

```shell
$ operator-sdk create api --group=simple-wordpress --version=v1alpha1 --kind=SimpleWordpress --resource=true --controller=true
```

The command creates the following files:
```
api/v1/simplewordpress_types.go
controllers/simplewordpress_controller.go
```

Change the `api/v1/simplewordpress_types.go`to include a new `SqlRootPassword` attribute in the WordpressSpec struct.

```go
// SimpleWordpressSpec defines the desired state of SimpleWordpress
type SimpleWordpressSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// MySql root password
	SqlRootPassword string `json:"sqlRootPassword"`
}
```

After you make the change in `api/v1/simplewordpress_types.go`, run the following commands to ensure that the changes are reflected in the. custom resource definitions (CRDs) in the config directory:
```shell
$ make generate
$ make manifests
```

### Implement Controller

Change `controllers\simplewordpress_controller.go`.

The boiler-plate generated the`SetupWithManager` function creates watch for the primary resources:

```go
For(&simplewordpressv1alpha1.SimpleWordpress{})
```

Watch also Deployment and Services as additional resources:
```go
func (r *SimpleWordpressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&simplewordpressv1alpha1.SimpleWordpress{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
```

Change the `reconcile` loop function to check and create MySQL and WordPress deployment and service.

Note: In this simple example we don't check or create any PV and PVC.

### Run the Kubernetes Operator

You need a Kubernetes cluster like [K3s](https://k3d.io), [Kind](https://kind.sigs.k8s.io/), [Minikube](https://minikube.sigs.k8s.io/docs/start/).
Set it as your current kube context.

Enter the following command to build and run operator:
```shell
$ make install run
```

Your `config/samples/simple-wordpress_v1alpha1_simplewordpress.yaml` should look something like:

```yaml
apiVersion: simple-wordpress.simplewordpress.com/v1alpha1
kind: SimpleWordpress
metadata:
  name: simplewordpress-sample
spec:
  sqlRootPassword: "mypassword"

```

Create the `SimpleWorkdpress` deployment sample using Kustomize
```
$ kubectl create -k config/samples/
```

After few seconds the MySql and WordPress deployments and services are created.
If you try to delete a deployment it will be recreated by operator.

