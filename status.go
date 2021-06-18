// Package apps exposes a few handy APIs which interact and retrieve information from the kubernetes cluster
package apps

import (
	"context"
	"log"
	"time"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	//  defaultNamespace refers to the kubernetes' "default" namespace
	defaultNamespace = "default"
)

// Pod represents the information of the pod present in the kubernetes cluster.
// The info consists of Name of the pod, Status if the pod is Running, Total Restart count of all the containers,
// The age of the pod since it is up
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

// getPodPhaseStatus returns the pod status depending upon its containers' statuses
func getPodPhaseStatus(pod apiv1.Pod) string {
	log.Println("Getting the pod status, Pod: %s", pod.ObjectMeta.Name)
	containerStatuses := pod.Status.ContainerStatuses
	for index := 0; index < len(containerStatuses); index++ {
		// returning the reason if a container is in waiting state.
		// The status of a given pod is considered 'Running' only if all the containers inside that pod are 'Running'
		if containerStatuses[index].State.Waiting != nil {
			return containerStatuses[index].State.Waiting.Reason
		}
	}
	// returning the pod status if all the containers are in non-Waiting state
	return string(pod.Status.Phase)
}

// getPodRestartCount returns the restart count of a pod.
// Restart Count is the sum of the restart counts of all the containers present in the given pod.
func getPodRestartCount(pod apiv1.Pod) int32 {
	log.Println("Getting the pod restart count, Pod: %s", pod.ObjectMeta.Name)
	containerStatuses := pod.Status.ContainerStatuses
	var restartCount int32
	for index := 0; index < len(containerStatuses); index++ {
		restartCount += containerStatuses[index].RestartCount
	}
	return restartCount
}

// GetPods is an API to initialize the connection with Kubernetes and fetch the details of all the pods present in a given "namespace". namespace defaults to the "default" if the argument passed is an empty string ("")
// (TODO) This can be modified as a handler function incase if a REST server is exposed towards the user
func GetPods(namespace string) []Pod {
	if namespace == "" {
		namespace = defaultNamespace
	}
	log.Println("Getting the pods information, Namespace: %s", namespace)
	var pods []Pod

	// (TODO) Incluster or an Out-of-cluster Configuration can be initialized.
	// (TODO) There can be a seperate API to initialize the connection
	// (TODO) As of now, proceeding with the creation of an In-Cluster configuration which limits user scope to call this API only from within the kubernetes cluster
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Println("Creating InCluster Configuration failed, Error: %v", err)
		return nil
	}

	// Creating a clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Println("Clientset creation failed, Error: %v", err)
		return nil
	}

	// Getting Pod information
	response, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	for _, info := range response.Items {
		pod := new(Pod)
		pod.Name = info.ObjectMeta.Name
		pod.Status = getPodPhaseStatus(info)
		pod.RestartCount = int(getPodRestartCount(info))
		pod.UpTime = float64(time.Now().Unix() - info.Status.StartTime.Unix())
		pods = append(pods, *pod)
	}
	log.Println("Fetched information successfully, Info: %v", pods)
	return pods
}
