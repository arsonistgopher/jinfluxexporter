package junoscollector

import (
	"fmt"

	"github.com/arsonistgopher/jkafkaexporter/alarm"
	"github.com/arsonistgopher/jkafkaexporter/collector"
	"github.com/arsonistgopher/jkafkaexporter/environment"
	"github.com/arsonistgopher/jkafkaexporter/interfaces"
	"github.com/arsonistgopher/jkafkaexporter/routingengine"
	"github.com/prometheus/client_golang/prometheus"

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
}

// NewJunosCollector for creating a new collector map
func NewJunosCollector() *JunosCollector {
	collectors := collectors()
	return &JunosCollector{collectors}
}

func collectors() map[string]collector.RPCCollector {
	m := map[string]collector.RPCCollector{
		"alarm": alarm.NewCollector(""),
	}

	// Include each of the stats here
	m["interfaces"] = interfaces.NewCollector()
	m["routing-engine"] = routingengine.NewCollector()
	m["environment"] = environment.NewCollector()
	m["interfaces"] = interfaces.NewCollector()

	return m
}

// Collect implements prometheus.Collector interface
func (c *JunosCollector) Collect(ch chan<- string, label string) {

	client, err := rpc.Create()

	if err != nil {
		fmt.Println(err)
	}

	c.collectForHost(client, ch, label)

	client.Close()
	fmt.Println("Exited from Collect()")
}

func (c *JunosCollector) collectForHost(client *rpc.Client, ch chan<- string, label string) {

	for k, col := range c.collectors {
		// fmt.Println("DEBUG: Collection > ", k)
		err := col.Collect(client, ch, label)
		if err != nil && err.Error() != "EOF" {
			fmt.Print("ERROR: " + k + ": " + err.Error())
		}
	}

	// fmt.Println("Exited from collectForHost")
}
