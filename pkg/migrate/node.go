package migrate

import (
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NodeInfo struct {
	Hostname      string
	HostIP        string
	NodeNamespace string
	Node          *v1.Node
	clientSet     *kubernetes.Clientset
}

func NewNodeInfo(hostname string, clientSet *kubernetes.Clientset) (*NodeInfo, error) {
	var err error
	n := new(NodeInfo)
	n.clientSet = clientSet
	n.Node, err = n.clientSet.CoreV1().Nodes().Get(hostname, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}
	n.Hostname = hostname
	for _, addr := range n.Node.Status.Addresses {
		if addr.Type == v1.NodeInternalIP {
			n.HostIP = addr.Address
		}
	}
	return n, err
}

func (n *NodeInfo) MigrateNode(namespaces []string) error {
	err := checkNamespaces(namespaces)
	if err != nil {
		return err
	}
	if len(namespaces) == 1 {
		n.NodeNamespace = namespaces[0]
	} else {
		n.NodeNamespace = Namespace_MIME
	}
	n.Node.Labels["namespace"] = n.NodeNamespace
	labels := []string{
		"beta.kubernetes.io/arch",
		"beta.kubernetes.io/os",
		"kubernetes.io/arch",
		"kubernetes.io/hostname",
		"kubernetes.io/os",
	}
	for _, label := range labels {
		delete(n.Node.Labels, label)
	}
	return nil
}
