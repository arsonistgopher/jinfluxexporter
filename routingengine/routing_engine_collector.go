package routingengine

import (
	"fmt"

	"github.com/arsonistgopher/jkafkaexporter/collector"
	"github.com/arsonistgopher/jkafkaexporter/rpc"
)

type routingEngineCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	c := &routingEngineCollector{}
	return c
}

// Collect collects metrics from JunOS
func (c *routingEngineCollector) Collect(client rpc.Client, ch chan<- string, label string) error {
	x := &RoutingEngineRpc{}
	err := rpc.RunCommandAndParse(client, "<get-route-engine-information/>", &x)
	if err != nil {
		return err
	}

	stats := x.RouteEngine

	jsonReturn := "{Node: %s, REngine: {Temp: %f, MemoryUtilised: %f, CPUTemp: %f, CPUUser: %f, CPUBackground: %f, " +
		"CPUSystem: %f, CPUInterrupt: %f, CPUIdle: %f, LoadAverageOne: %f, LoadAverageFive: %f, LoadAverageFifteen: %f}}"

	ch <- fmt.Sprintf(jsonReturn, label, stats.Temperature.Value, stats.MemoryUtilization, stats.CPUTemperature.Value,
		stats.CPUUser, stats.CPUBackground, stats.CPUSystem, stats.CPUInterrupt, stats.CPUIdle, stats.LoadAverageOne,
		stats.LoadAverageFive, stats.LoadAverageFifteen)

	return nil
}
