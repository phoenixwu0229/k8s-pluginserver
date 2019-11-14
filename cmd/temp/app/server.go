package app

import (
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func Newk8sPluginServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "k8s-pluginserver",
		Long: "The Kubernetes scheduler provides a relationship with node and namespace",
		Run: func(cmd *cobra.Command, args []string) {
			klog.Info("hello world")
		},
	}
	//cmd.Usage()
	return cmd
}
