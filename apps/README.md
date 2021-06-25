# apps
--
    import "k8s/apps"

Package apps exposes a few handy APIs which interact and retrieve information
from the kubernetes cluster

## Usage

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

#### func  GetPods

```go
func GetPods(namespace string) []Pod
```
GetPods is an API to initialize the connection with Kubernetes and fetch the
details of all the pods present in a given "namespace". namespace defaults to
the "default" if the argument passed is an empty string ("") (TODO) This can be
modified as a handler function incase if a REST server is exposed towards the
user
