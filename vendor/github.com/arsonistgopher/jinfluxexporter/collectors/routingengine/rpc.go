package routingengine

// RoutingEngineRpc top level type
type RoutingEngineRpc struct {
	RouteEngine RouteEngine `xml:"route-engine"`
}

// RouteEngine type
type RouteEngine struct {
	Temperature        RouteEngineTemperature `xml:"temperature"`
	MemoryUtilization  float64                `xml:"memory-buffer-utilization"`
	CPUTemperature     RouteEngineTemperature `xml:"cpu-temperature"`
	CPUUser            float64                `xml:"cpu-user"`
	CPUBackground      float64                `xml:"cpu-background"`
	CPUSystem          float64                `xml:"cpu-system"`
	CPUInterrupt       float64                `xml:"cpu-interrupt"`
	CPUIdle            float64                `xml:"cpu-idle"`
	LoadAverageOne     float64                `xml:"load-average-one"`
	LoadAverageFive    float64                `xml:"load-average-five"`
	LoadAverageFifteen float64                `xml:"load-average-fifteen"`
}

// RouteEngineTemperature type
type RouteEngineTemperature struct {
	Value float64 `xml:"celsius,attr"`
}
