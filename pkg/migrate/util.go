package migrate

import (
	"fmt"
	"strings"

	appsV1 "k8s.io/api/apps/v1"
)

var (
	modelDeploys = []string{
		"bigscale-click-fm-001-20190228-19006-aliyun",
		"bigscale-click-fm-001-20190228-19006-youfu",
		"bigscale-click-fm-002-20190228-19010-aliyun",
		"bigscale-click-fm-002-20190228-19010-youfu",
		"bigscale-interaction-fm-013-20181214-1-19005-aliyun",
		"bigscale-interaction-fm-013-20181214-1-19005-youfu",
		"bigscale-read-fm-001-20181213-19003-aliyun",
		"bigscale-read-fm-001-20181213-19003-youfu",
		"channel-lr-16561-dabailou",
		"channel-lr-16562-dabailou",
		"channel-lr-16565-dabailou",
		"channel-lr-16566-dabailou",
		"channel-lr-16569-dabailou",
		"online-model-test-youfu",
	}
	mimeDeploys = []string{
		"audio-classifier-17081",
		"face-attributes-on-17028",
		"labelling-v2-16223",
		"labelling-v2-16323",
		"logo1-recognition-17080",
		"ms-download-sina-x001-17100",
		"ms-download-sina-x001-17101",
		"ms-image-vgg-vec-download-17100",
		"ms-imagefinger-17059",
		"ms-label-v3-16224",
		"ms-label-v4-16225",
		"ms-logo-downloader-17000",
		"ms-logo-videoprocess-17090",
		"ocr-text-down-17000",
		"ocr-text-on-17027",
		"pic-down-17000",
	}
	//Namespace_Default = "default"
	Namespace_MIME    = "multimedia-service"
	Namespace_Model   = "model-service"
	Namespace_Weidis  = "weidis-service"
	Namespace_Weips   = "weips-service"
	Namespace_Feature = "feature-service"
	Namespace_Desired = []string{
		"multimedia-service",
		"model-service",
		"weidis-service",
		"weips-service",
		"feature-service",
	}
)

func ChangeNamespace(deploy *appsV1.Deployment) string {
	namespace := deploy.Namespace
	name := deploy.Name
	// multimedia-service
	if namespace == "chuanlong" || arrayContain(name, mimeDeploys) {
		namespace = Namespace_MIME
	}
	// model-service
	if namespace == "model-service" {
		namespace = Namespace_Model
	} else if namespace == "default" && strings.Contains(name, "mainp") {
		namespace = Namespace_Model
	} else if namespace == "default" && strings.Contains(name, "feed") && !strings.Contains(name, "weips") {
		namespace = Namespace_Model
	} else if namespace == "default" && strings.Contains(name, "channel") {
		namespace = Namespace_Model
	} else if namespace == "default" && arrayContain(name, modelDeploys) {
		namespace = Namespace_Model
	}
	// model service
	if namespace == "weiss" || namespace == "weidis" {
		namespace = Namespace_Weidis
	}
	// weips
	if namespace == "default" && strings.Contains(name, "weips") {
		namespace = Namespace_Weips
	}
	// feature-service
	if namespace == "feature-service" {
		namespace = Namespace_Feature
	}
	return namespace
}

func isDesiredNamespace(namespace string) bool {
	return arrayContain(namespace, Namespace_Desired)
}

func arrayContain(s string, slice []string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func checkNamespaces(namespaces []string) error {
	var err error
	if len(namespaces) == 1 {
		if !isDesiredNamespace(namespaces[0]) {
			err = fmt.Errorf(
				"namespace is not in the desired namespaces list, namespaces illegal, namespaces: %v",
				namespaces)
		}
	} else if arrayContain(Namespace_MIME, namespaces) {
		if !arrayContain(Namespace_Model, namespaces) || len(namespaces) > 2 {
			err = fmt.Errorf("namespaces are not exactly mime and model, namespaces illegal, namespaces: %v", namespaces)
		}
	} else {
		err = fmt.Errorf("namspaces are not contain mime, namespaces illegal, namespaces: %v", namespaces)
	}
	return err
}
