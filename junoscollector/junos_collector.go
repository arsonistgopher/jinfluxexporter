package junoscollector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/crypto/ssh"

	"github.com/arsonistgopher/jkafkaexporter/collector"
	"github.com/arsonistgopher/jkafkaexporter/rpc"
)

const prefix = "junos_"

var (
	scrapeDurationDesc *prometheus.Desc
	upDesc             *prometheus.Desc
)

func init() {
	upDesc = prometheus.NewDesc(prefix+"up", "Scrape of target was successful", []string{"target"}, nil)
	scrapeDurationDesc = prometheus.NewDesc(prefix+"collector_duration_seconds", "Duration of a collector scrape for one target", []string{"target"}, nil)
}

// JunosCollector for export
type JunosCollector struct {
	collectors map[string]collector.RPCCollector
	target     string
	sshconfig  *ssh.ClientConfig
	port       int
}

// NewJunosCollector for creating a new collector map
func NewJunosCollector(sshconfig *ssh.ClientConfig, port int, target string) *JunosCollector {
	return &JunosCollector{sshconfig: sshconfig, port: port, target: target, collectors: make(map[string]collector.RPCCollector)}
}

// Add a collector to the collector map
func (c *JunosCollector) Add(name string, newcollector collector.RPCCollector) error {
	c.collectors[name] = newcollector
	return nil
}

// Collect implements prometheus.Collector interface
func (c *JunosCollector) Collect(ch chan<- string, label string) {

	client, err := rpc.Create(c.sshconfig, c.target, c.port)

	if err != nil {
		fmt.Println(err)
	}

	c.collectForHost(client, ch, label)
	rpc.Close(client)
}

func (c *JunosCollector) collectForHost(client *rpc.Client, ch chan<- string, label string) {

	for k, col := range c.collectors {
		// fmt.Println("DEBUG: Collection > ", k)
		err := col.Collect(*client, ch, label)
		if err != nil && err.Error() != "EOF" {
			fmt.Print("ERROR: " + k + ": " + err.Error())
		}
	}
}
