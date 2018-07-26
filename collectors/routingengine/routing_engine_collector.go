package routingengine

import (
	"time"

	"github.com/arsonistgopher/jinfluxdbexporter/collector"
	"github.com/arsonistgopher/jinfluxdbexporter/internal/channels"
	"github.com/arsonistgopher/jinfluxdbexporter/rpc"
)

type routingEngineCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	c := &routingEngineCollector{}
	return c
}

// Collect collects metrics from JunOS
func (c *routingEngineCollector) Collect(client rpc.Client, ch chan<- channels.InfluxDBMeasurement, label string, measurement string) error {
	x := &RoutingEngineRpc{}
	err := rpc.RunCommandAndParse(client, "<get-route-engine-information/>", &x)
	if err != nil {
		return err
	}

	stats := x.RouteEngine

	tagset := make(map[string]string)
	tagset["host"] = label

	fieldset := map[string]interface{}{
		"temp":           uint64(stats.Temperature.Value),
		"memoryUtilised": uint64(stats.MemoryUtilization),
		"cpuTemp":        uint64(stats.CPUTemperature.Value),
		"cpuUser":        uint64(stats.CPUUser),
		"cpuBackground":  uint64(stats.CPUBackground),
		"cpuSystem":      uint64(stats.CPUSystem),
		"cpuInterrupt":   uint64(stats.CPUInterrupt),
		"cpuIdle":        uint64(stats.CPUIdle),
		"loadAverage1":   uint64(stats.LoadAverageOne),
		"loadAverage5":   uint64(stats.LoadAverageFive),
		"loadAverage15":  uint64(stats.LoadAverageFifteen),
	}

	ch <- channels.InfluxDBMeasurement{Measurement: measurement, TagSet: tagset, FieldSet: fieldset, TimeStamp: time.Now()}

	return nil
}
