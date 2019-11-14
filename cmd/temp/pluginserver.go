package main

import (
	"flag"
	"k8s-pluginserver/cmd/temp/app"
	logs "k8s-pluginserver/pkg/log"

	"github.com/spf13/pflag"
	"k8s.io/klog"
)

func main() {
	command := app.Newk8sPluginServerCommand()
	logs.InitLogs()
	flag.Parse()
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	klog.Flush()
	defer logs.FlushLogs()
	if err := command.Execute(); err != nil {
		klog.Fatalf("root cmd execute failed")
	}
}
