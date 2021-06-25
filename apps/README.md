# apps
--
    import "k8s/apps"

Package apps exposes a few handy APIs which interact and retrieve information
from the kubernetes cluster

## Usage

```go
const (
	// InCluster refers to one of the configuration types by which the kubernetes cluster can be accessed.
	// This configuration helps in initializing the authentication to the Kubernetes API from an application running inside the Kubernetes cluster.
	// Remember to run the following command to create role binding which will grant the default service account view permissions.
	// Command: `kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default`
	InCluster configType = "In-Cluster"
	// OutOfCluster refers to one of the configuration types by which the kubernetes cluster can be accessed.
	// This type of configuration initializes the authentication to the Kubernetes API from an application running outside the Kubernetes cluster.
	OutOfCluster configType = "Out-Of-Cluster"
)
```

#### type Client

```go
type Client struct {
	// Clientset refers to the actual clientset of kubernetes go client that interacts with the Kubernetes API
	*kubernetes.Clientset
}
```

Client acts as a config holder which interacts with the Kubernetes API

#### func  NewClient

```go
func NewClient(confType configType) *Client
```
NewClient is a constructor function which initializes and returns the client
that can interact with the Kubernetes API based on the provided configuration
type

#### func (*Client) GetPods

```go
func (cli *Client) GetPods(namespace string) []Pod
```
GetPods is an API to initialize the connection with Kubernetes and fetch the
details of all the pods present in a given "namespace". namespace defaults to
the "default" if the argument passed is an empty string ("") (TODO) This can be
modified as a handler function incase if a REST server is exposed towards the
user

#### type Pod

```go
type Pod struct {
	// Name of the pod
	Name string
	// Status of the pod ex:"Running/CrashLoopBack/Error" etc.
	Status string
	// RestartCount refers to the sum of the restart counts of all the containers in a pod
	RestartCount int
	// UpTime represents the age of the pod
	UpTime float64
}
```

Pod represents the information of the pod present in the kubernetes cluster. The
info consists of Name of the pod, Status if the pod is Running, Total Restart
count of all the containers, The age of the pod since it is up
