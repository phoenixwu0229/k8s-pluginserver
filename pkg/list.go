package pkg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}
func getKubernetesConfig() *rest.Config {
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(homeDir(), ".kube", "flink_kubeconfig"))
	}
	if err != nil {
		log.Fatal("failed to get kubernetes configuration")
	}
	return config
}

func Run() {
	client, err := kubernetes.NewForConfig(getKubernetesConfig())
	if err != nil {
		log.Fatal("failed init client")
	}
	factory := informers.NewSharedInformerFactory(client, 0)
	fmt.Print(factory)
	//nodeLister := factory.Core().V1().Nodes().Lister()
	//nodes, err := nodeLister.List(labels.Everything())
	//if err != nil {
	//}
}
