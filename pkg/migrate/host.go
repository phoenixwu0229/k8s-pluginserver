package migrate

import (
	"log"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	KindDeployment = "Deployment"
	KindReplicaSet = "ReplicaSet"
)

var (
	GetEverything = metaV1.GetOptions{}
	//ListEverything = metaV1.ListOptions{}
)

type HostInfo struct {
	*NodeInfo
	*PodsInfo
	oldClient *kubernetes.Clientset
	NewClient *kubernetes.Clientset
}

func NewHost(hostname string, oldConfig string, newConfig string) (*HostInfo, error) {
	var err error
	host := new(HostInfo)
	host.NewClient = NewClient(newConfig)
	host.oldClient = NewClient(oldConfig)
	host.PodsInfo, err = NewPods(hostname, host.oldClient)
	if err != nil {
		return host, err
	}
	host.NodeInfo, err = NewNodeInfo(hostname, host.oldClient)
	return host, err
}

func NewClient(path string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		log.Fatal(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return clientSet
}
