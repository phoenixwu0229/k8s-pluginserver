package main

import (
	"flag"
	"k8s-pluginserver/pkg/migrate"
	"log"
	"os"
)

var (
	oldConfigPath = new(string)
	newConfigPath = new(string)
	labelFilePath = new(string)
	hostname      = new(string)
)

func init() {
	flag.StringVar(oldConfigPath, "old_config_path", "", "")
	flag.StringVar(newConfigPath, "new_config_path", "", "")
	flag.StringVar(labelFilePath, "label_file_path", "", "")
	flag.StringVar(hostname, "hostname", "", "")
}
func run(hostname string) error {
	var labels string
	file, err := os.OpenFile(*labelFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer func() {
		err = file.Close()
	}()
	if err != nil {
		return err
	}

	host, err := migrate.NewHost(hostname, *oldConfigPath, *newConfigPath)
	if err != nil {
		return err
	}
	err = host.MigrateDeploy(host.NewClient)
	if err != nil {
		return err
	}
	err = host.MigrateNode(host.Namespaces)
	if err != nil {
		return err
	}
	for k, v := range host.Node.Labels {
		labels = labels + " " + k + "=" + v
	}
	_, err = file.Write([]byte(labels))
	return err
}

func main() {
	var err error
	flag.Parse()
	client := migrate.NewClient(*oldConfigPath)
	_, err = client.CoreV1().Nodes().Get(*hostname, migrate.GetEverything)
	if err != nil {
		log.Fatal(err)
	}
	err = run(*hostname)
	if err != nil {
		log.Fatal(err)
	}
}
