// Package apps exposes a few handy APIs which interact and retrieve information from the kubernetes cluster
package apps

import (
	"context"
	"flag"
	"log"
	"path/filepath"
	"time"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	//  defaultNamespace refers to the kubernetes' "default" namespace
	defaultNamespace = "default"
)

// configType refers to the types of modes through which the Kubernetes API can be accessed.
type configType string

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

// Client acts as a config holder which interacts with the Kubernetes API
type Client struct {
	// Clientset refers to the actual clientset of kubernetes go client that interacts with the Kubernetes API
	*kubernetes.Clientset
}

// NewClient is a constructor function which initializes and returns the client that can interact with the Kubernetes API based on the provided configuration type
func NewClient(confType configType) *Client {
	log.Printf("Initializing the client configuration, Config Type: %v\n", confType)
	if confType == InCluster {
		config, err := rest.InClusterConfig()
		if err != nil {
			log.Printf("Creating InCluster Configuration failed, Error: %v\n", err)
			return nil
		}

		// Creating a clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Printf("Clientset creation failed, Error: %v\n", err)
			return nil
		}
		return &Client{clientset}

	} else if confType == OutOfCluster {
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Printf("Creating Out of Cluster Configuration failed, Error: %v\n", err)
			return nil
		}
		// Creating a clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Printf("Clientset creation failed, Error: %v\n", err)
			return nil
		}
		return &Client{clientset}
	}
	log.Printf("Initializing the configuration failed, Invalid Config type: %v\n", confType)
	return nil
}

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
	containerStatuses := pod.Status.ContainerStatuses
	var restartCount int32
	for index := 0; index < len(containerStatuses); index++ {
		restartCount += containerStatuses[index].RestartCount
	}
	return restartCount
}

// GetPods is an API to fetch the details of all the pods present in a given "namespace". namespace defaults to the "default" if the argument passed is an empty string ("")
func (cli *Client) GetPods(namespace string) []Pod {
	if namespace == "" {
		namespace = defaultNamespace
	}
	log.Printf("Getting the pods information, Namespace: %s\n", namespace)
	var pods []Pod

	// Getting Pod information
	response, err := cli.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed getting response from k8s API, Err: %v", err)
		return nil
	}
	for _, info := range response.Items {
		pod := new(Pod)
		pod.Name = info.ObjectMeta.Name
		pod.Status = getPodPhaseStatus(info)
		pod.RestartCount = int(getPodRestartCount(info))
		pod.UpTime = float64(time.Now().Unix() - info.Status.StartTime.Unix())
		pods = append(pods, *pod)
	}
	log.Printf("Fetched information successfully, Info: %v\n", pods)
	return pods
}

// GetEvents is an API to fetch the events that were recorded in the kubernetes cluster
// "namespace" defaults to the "default" if provided as an empty string("")
func (cli *Client) GetEvents(namespace string) interface{} {
	if namespace == "" {
		namespace = defaultNamespace
	}
	log.Printf("Getting the events information, Namespace: %s\n", namespace)
	events, err := cli.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed getting response from k8s API, Err: %v", err)
		return nil
	}
	return events
}
