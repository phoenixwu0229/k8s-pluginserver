package migrate

import (
	"fmt"
	"log"
	"strings"

	"k8s.io/client-go/kubernetes"

	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodsInfo struct {
	hostname   string
	Pods       []v1.Pod
	Deploys    map[string]*appsV1.Deployment
	Namespaces []string
	clientSet  *kubernetes.Clientset
}

func (p *PodsInfo) generatePodsAndDeploys() error {
	podList, err := p.clientSet.
		CoreV1().
		Pods("").
		List(metaV1.ListOptions{FieldSelector: "spec.nodeName=" + p.hostname})
	if err != nil {
		return fmt.Errorf("failed to list pod: %v", err)
	}
	p.Pods = podList.Items
	for _, item := range podList.Items {
		deploy, err := p.generateDeployWithRef(item.Namespace, item.OwnerReferences)
		if err != nil {
			return err
		}
		if deploy != nil {
			p.Deploys[item.Name] = deploy
		}
	}
	return nil
}

func (p *PodsInfo) generateNamespaces() error {
	var err error
	for _, deploy := range p.Deploys {
		namespace := ChangeNamespace(deploy)
		deploy.Namespace = namespace
		if !arrayContain(namespace, p.Namespaces) {
			p.Namespaces = append(p.Namespaces, namespace)
		}
	}
	err = checkNamespaces(p.Namespaces)
	return err
}

func (p *PodsInfo) generateDeployWithRef(namespace string, ownerRef []metaV1.OwnerReference) (*appsV1.Deployment, error) {
	var err error
	for _, ref := range ownerRef {
		if ref.Kind == KindReplicaSet {
			result, err := p.clientSet.AppsV1().ReplicaSets(namespace).Get(ref.Name, metaV1.GetOptions{})
			for err != nil {
				if strings.Contains(err.Error(), "allotted") || strings.Contains(err.Error(), "unexpected EOF") || strings.Contains(err.Error(), "INTERNAL_ERROR") {
					result, err = p.clientSet.AppsV1().ReplicaSets(namespace).Get(ref.Name, metaV1.GetOptions{})
				} else {
					return nil, err
				}
			}
			return p.generateDeployWithRef(namespace, result.OwnerReferences)
		} else if ref.Kind == KindDeployment {
			result, err := p.clientSet.AppsV1().Deployments(namespace).Get(ref.Name, metaV1.GetOptions{})
			for err != nil {
				if strings.Contains(err.Error(), "allotted") || strings.Contains(err.Error(), "unexpected EOF") || strings.Contains(err.Error(), "INTERNAL_ERROR") {
					result, err = p.clientSet.AppsV1().Deployments(namespace).Get(ref.Name, metaV1.GetOptions{})
				} else {
					return nil, err
				}
			}
			return result, err
		}
	}
	return nil, err
}

func (p *PodsInfo) formatDeploy() {
	//if h.Node.Labels == nil {
	//	h.Node.Labels = make(map[string]string)
	//}
	for _, deploy := range p.Deploys {
		deploy.Namespace = ChangeNamespace(deploy)
		if len(p.Namespaces) == 1 || deploy.Namespace == Namespace_MIME {
			deploy.Spec.Template.Spec.NodeSelector["namespace"] = deploy.Namespace
		}
		deploy.Status = appsV1.DeploymentStatus{}
		deploy.ObjectMeta = metaV1.ObjectMeta{
			Namespace:   deploy.ObjectMeta.Namespace,
			Annotations: deploy.ObjectMeta.Annotations,
			Labels:      deploy.ObjectMeta.Labels,
			Name:        deploy.ObjectMeta.Name,
		}
	}
}

func (p *PodsInfo) MigrateDeploy(newClientSet *kubernetes.Clientset) error {
	p.formatDeploy()
	for name, deploy := range p.Deploys {
		log.Printf("pod: %q belong to deployment: %q is migrating", name, deploy.Name)
		_, err := newClientSet.AppsV1().Deployments(deploy.Namespace).Get(deploy.Name, GetEverything)
		for err != nil && !strings.Contains(err.Error(), "not found") {
			_, err = newClientSet.AppsV1().Deployments(deploy.Namespace).Get(deploy.Name, GetEverything)
		}
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				_, err := newClientSet.AppsV1().Deployments(deploy.Namespace).Create(deploy)
				for err != nil {
					_, err = newClientSet.AppsV1().Deployments(deploy.Namespace).Create(deploy)
				}
			}
		}
		log.Printf("deployment: %q migrate success", deploy.Name)
	}
	return nil
}

func NewPods(hostname string, clientSet *kubernetes.Clientset) (*PodsInfo, error) {
	var err error
	p := new(PodsInfo)

	p.hostname = hostname
	p.clientSet = clientSet
	p.Deploys = make(map[string]*appsV1.Deployment)
	err = p.generatePodsAndDeploys()
	if err != nil {
		return nil, err
	}
	err = p.generateNamespaces()
	if err != nil {
		return nil, err
	}
	return p, nil
}
