package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type clusterNamespace struct {
	name           string
	numPods        int
	numDeployments int
	numServices    int
}

func printClusterInfo(m map[int]clusterNamespace) {
	var (
		totalPods        int
		totalDeployments int
		totalServices    int
	)

	for _, ns := range m {
		fmt.Printf("Namespace: %s\n", ns.name)
		fmt.Printf("\tPods: %d\n", ns.numPods)
		fmt.Printf("\tDeployments: %d\n", ns.numDeployments)
		fmt.Printf("\tServices: %d\n", ns.numServices)
		totalPods += ns.numPods
		totalDeployments += ns.numDeployments
		totalServices += ns.numServices
	}
	fmt.Printf("\nCluster totals:\n")
	fmt.Printf("\tPods: %d\n", totalPods)
	fmt.Printf("\tDeployments: %d\n", totalDeployments)
	fmt.Printf("\tServices: %d\n", totalServices)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// get the namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	clusterInfo := make(map[int]clusterNamespace)

	// loop through the namespaces
	for i, namespace := range namespaces.Items {
		namespaceName := namespace.Name

		// get the pods
		pods, err := clientset.CoreV1().Pods(namespaceName).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		// get the services
		services, err := clientset.CoreV1().Services(namespaceName).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		// get the deployments
		deployments, err := clientset.ExtensionsV1beta1().Deployments(namespaceName).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		clusterInfo[i] = clusterNamespace{
			name:           namespaceName,
			numPods:        len(pods.Items),
			numDeployments: len(deployments.Items),
			numServices:    len(services.Items),
		}
	}

	printClusterInfo(clusterInfo)
}
