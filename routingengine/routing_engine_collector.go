package routingengine

import (
	"fmt"

	"github.com/arsonistgopher/jkafkaexporter/collector"
	"github.com/arsonistgopher/jkafkaexporter/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

const prefix string = "junos_route_engine_"

var (
	temperature        *prometheus.Desc
	memoryUtilization  *prometheus.Desc
	cpuTemperature     *prometheus.Desc
	cpuUser            *prometheus.Desc
	cpuBackground      *prometheus.Desc
	cpuSystem          *prometheus.Desc
	cpuInterrupt       *prometheus.Desc
	cpuIdle            *prometheus.Desc
	loadAverageOne     *prometheus.Desc
	loadAverageFive    *prometheus.Desc
	loadAverageFifteen *prometheus.Desc
)

func init() {
	l := []string{"target"}
	temperature = prometheus.NewDesc(prefix+"temp", "Temperature of the air flowing past the Routing Engine (in degrees C)", l, nil)
	memoryUtilization = prometheus.NewDesc(prefix+"memory_utilization", "Percentage of Routing Engine memory being used", l, nil)
	cpuTemperature = prometheus.NewDesc(prefix+"cpu_temp", "Temperature of the CPU (in degrees C)", l, nil)
	cpuUser = prometheus.NewDesc(prefix+"cpu_user_percent", "Percentage of CPU time being used by user processes", l, nil)
	cpuBackground = prometheus.NewDesc(prefix+"cpu_background_percent", "Percentage of CPU time being used by background processes", l, nil)
	cpuSystem = prometheus.NewDesc(prefix+"cpu_system_percent", "Percentage of CPU time being used by kernel processes", l, nil)
	cpuInterrupt = prometheus.NewDesc(prefix+"cpu_interrupt_percent", "Percentage of CPU time being used by interrupts", l, nil)
	cpuIdle = prometheus.NewDesc(prefix+"cpu_idle_percent", "Percentage of CPU time that is idle", l, nil)
	loadAverageOne = prometheus.NewDesc(prefix+"load_average_one", "Routing Engine load averages for the last 1 minute", l, nil)
	loadAverageFive = prometheus.NewDesc(prefix+"load_average_five", "Routing Engine load averages for the last 5 minutes", l, nil)
	loadAverageFifteen = prometheus.NewDesc(prefix+"load_average_fifteen", "Routing Engine load averages for the last 15 minutes", l, nil)
}

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
