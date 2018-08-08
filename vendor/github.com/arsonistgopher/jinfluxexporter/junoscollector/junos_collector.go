package junoscollector

import (
	"fmt"

	channels "github.com/arsonistgopher/jinfluxexporter/rootchannels"

	"golang.org/x/crypto/ssh"

	"github.com/arsonistgopher/jinfluxexporter/collector"
	"github.com/arsonistgopher/jinfluxexporter/rpc"
)

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

// Len (number) of collectors registered
func (c *JunosCollector) Len() int {
	return len(c.collectors)
}

// Collect implements interface
func (c *JunosCollector) Collect(ch chan<- channels.InfluxDBMeasurement, label string) {

	client, err := rpc.Create(c.sshconfig, c.target, c.port)

	if err != nil {
		fmt.Println(err)
	}

	c.collectForHost(client, ch, label)

	err = rpc.Close(client)

	if err != nil {
		fmt.Println(err)
	}
}

func (c *JunosCollector) collectForHost(client *rpc.Client, ch chan<- channels.InfluxDBMeasurement, label string) {

	for k, col := range c.collectors {
		err := col.Collect(*client, ch, label, k)
		if err != nil && err.Error() != "EOF" {
			fmt.Print("ERROR: " + k + ": " + err.Error())
		}
	}
}
