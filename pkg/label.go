package pkg

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type NodeLabelList map[string]*Label

func (n NodeLabelList) GenerateCsvInfo() {
	file, _ := os.Create("node.csv")
	c := csv.NewWriter(file)
	err := c.Write([]string{"", "label", "number", "hosts"})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		c.Flush()
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	for k, v := range defaultNodeLabelList {
		l := append(make([]string, 1), k, strconv.Itoa(v.number), sliceToString(v.hosts))
		err := c.Write(l)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (n NodeLabelList) Read() {
	f, err := os.OpenFile("nodelabels.log", os.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	br := bufio.NewReader(f)
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			break
		} else {
			n.readLine(string(line))
		}
	}
}

func (n NodeLabelList) readLine(line string) {
	labels := strings.Split(line, ",")
	var host string
	for _, label := range labels {
		keys := strings.Split(label, "=")
		key, value := keys[0], keys[1]
		if key == "kubernetes.io/hostname" {
			host = value
			break
		}
	}
	for _, label := range labels {
		keys := strings.Split(label, "=")
		key, value := keys[0], keys[1]
		if key != "kubernetes.io/hostname" {
			if _, ok := n[label]; !ok {
				n[label] = NewLabel(key, value)
			}
			n[label].addHost(host)
		}
	}
}

var defaultNodeLabelList NodeLabelList

func init() {
	defaultNodeLabelList = make(map[string]*Label)
}

type Label struct {
	hosts      []string
	labelKey   string
	labelValue string
	number     int
}

func NewLabel(k, v string) *Label {
	label := new(Label)
	hosts := make([]string, 1, 1200)
	label.hosts = hosts
	label.labelKey = k
	label.labelValue = v
	label.number = 0
	return label
}
func (l *Label) addHost(host string) {
	l.hosts = append(l.hosts, host)
	l.number = l.number + 1
}

func sliceToString(slice []string) string {
	var result string
	for _, s := range slice {
		if result == "" {
			result = s
		} else {
			result = result + "," + s
		}
	}
	return result
}
